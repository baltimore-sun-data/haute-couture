package checker

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	CSS, HTMLDir string
	Include      *regexp.Regexp
	Exclude      *regexp.Regexp
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
	paths, err := conf.ListFiles()
	for _, path := range paths {
		fmt.Println(path)
	}
	return err
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
