package styles

import (
	"bytes"
	"testing"
)

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("err != nil: %v", err)
	}
}

func TestParsing(t *testing.T) {
	var tcs = []struct {
		name    string
		classes []string
		ids     []string
		css     string
	}{
		{
			"basic-class", []string{"foo"}, nil,
			`h1 .foo { background-color: green; }`,
		},
		{
			"basic-id", nil, []string{"bar"},
			`h1 #bar { background-color: green; }`,
		},
		{
			"media-rule", []string{"baz"}, []string{"spam"},
			`.baz { color: inherit; }
			@media max-width: 100px {
				h1 #spam { background-color: green; }
			}`,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewBufferString(tc.css)
			classes, ids, err := ExtractClassesAndIds(r)
			check(t, err)
			if len(classes) != len(tc.classes) {
				t.Fatalf("wrong number of classes: %d != %d",
					len(classes), len(tc.classes))
			}
			for _, class := range tc.classes {
				if _, ok := classes[class]; !ok {
					t.Fatalf("missing class: %q", class)
				}
			}
			if len(ids) != len(tc.ids) {
				t.Fatalf("wrong number of ids: %d != %d",
					len(ids), len(tc.ids))
			}
			for _, id := range tc.ids {
				if _, ok := ids[id]; !ok {
					t.Fatalf("missing id: %q", id)
				}
			}
		})
	}
}
