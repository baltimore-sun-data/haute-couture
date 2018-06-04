package checker

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/baltimore-sun-data/haute-couture/html"
	"github.com/baltimore-sun-data/haute-couture/styles"
)

type Config struct {
	CSS, HTMLDir, Output string
	Include              *regexp.Regexp
	Exclude              *regexp.Regexp
}

func NewConfig() Config {
	return Config{
		Include: regexp.MustCompile(`\.html?$`),
		Exclude: regexp.MustCompile(`^\.`),
	}
}

func FromArgs(args []string) Config {
	conf := NewConfig()
	fl := flag.NewFlagSet("haute-couture", flag.ExitOnError)
	fl.StringVar(&conf.CSS, "css", "", "CSS file to match against")
	fl.StringVar(&conf.HTMLDir, "html-dir", "public", "directory to search for HTML files")
	fl.StringVar(&conf.Output, "output", "extra-css.txt",
		"file to save any found extra CSS identifiers in")
	fl.Var(&regexpFlag{conf.Include}, "include", "regexp for HTML files to process")
	fl.Var(&regexpFlag{conf.Exclude}, "exclude", "regexp for sub-directories to exclude")
	fl.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`haute-couture

Usage of haute-couture:

`,
		)
		fl.PrintDefaults()
	}
	_ = fl.Parse(args)

	return conf
}

type regexpFlag struct{ *regexp.Regexp }

func (r *regexpFlag) String() string {
	return fmt.Sprintf(`"%s"`, r.Regexp)
}

func (r *regexpFlag) Set(s string) error {
	var err error
	r.Regexp, err = regexp.Compile(s)
	return err
}

func (conf Config) Execute() error {
	classes, ids, err := conf.ReadCSS()
	if err != nil {
		return err
	}
	paths, err := conf.ListFiles()
	if err != nil {
		return fmt.Errorf("problem listing HTML files: %v", err)
	}
	if err = conf.CheckPaths(classes, ids, paths); err != nil {
		return err
	}
	if len(classes) != 0 || len(ids) != 0 {
		fmt.Fprintf(os.Stderr, "found %d extra classes; %d extra ids\n",
			len(classes), len(ids))
		return conf.WriteOutput(classes, ids)
	}
	return nil
}

func (conf Config) ListFiles() ([]string, error) {
	var paths []string
	err := filepath.Walk(conf.HTMLDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(path)
		if info.IsDir() && conf.Exclude.MatchString(base) {
			return filepath.SkipDir
		}
		if !info.IsDir() && conf.Include.MatchString(base) {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

func (conf Config) ReadCSS() (classes, ids map[string]bool, err error) {
	f, err := os.Open(conf.CSS)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open CSS file: %v", err)
	}
	defer f.Close()
	classes, ids, err = styles.ExtractClassesAndIDs(f)
	return
}

func (conf Config) CheckPaths(classes, ids map[string]bool, paths []string) (err error) {
	for _, path := range paths {
		if err := conf.CheckPath(classes, ids, path); err != nil {
			return err
		}
	}
	return nil
}

func (conf Config) CheckPath(classes, ids map[string]bool, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open HTML file: %v", err)
	}
	defer f.Close()

	HTMLclasses, HTMLids, err := html.ExtractClassesAndIDs(f)
	if err != nil {
		return err
	}
	for class := range HTMLclasses {
		delete(classes, class)
	}
	for id := range HTMLids {
		delete(ids, id)
	}
	return nil
}

func (conf Config) WriteOutput(classes, ids map[string]bool) (err error) {
	f, err := os.Create(conf.Output)
	if err != nil {
		return fmt.Errorf("could not save results: %v", err)
	}

	var data struct {
		Classes []string
		IDs     []string
	}

	for class := range classes {
		data.Classes = append(data.Classes, class)
	}
	for id := range ids {
		data.IDs = append(data.IDs, id)
	}
	sort.Strings(data.Classes)
	sort.Strings(data.IDs)

	b, _ := json.MarshalIndent(&data, "", "  ")
	_, err = f.Write(b)
	if err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
