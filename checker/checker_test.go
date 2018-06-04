package checker

import (
	"testing"
)

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("err != nil: %v", err)
	}
}

func equals(t *testing.T, msg string, a, b interface{}) {
	t.Helper()
	if a != b {
		t.Fatalf(msg, a, b)
	}
}

func TestListing(t *testing.T) {
	var tcs = []struct {
		name  string
		path  string
		files []string
	}{
		{
			"basic-test", "testingdata/a", []string{
				"testingdata/a/1.htm", "testingdata/a/b/4.html",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			conf := NewConfig()
			conf.HTMLDir = tc.path
			files, err := conf.ListFiles()
			check(t, err)
			equals(t, "wrong number of files: %d != %d",
				len(files), len(tc.files))
			for i := range files {
				equals(t, "wrong file listed: %q != %q",
					files[i], tc.files[i])
			}
		})
	}
}
