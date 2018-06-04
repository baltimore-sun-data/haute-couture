package html

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
		html    string
	}{
		{
			"basic-class", []string{"foo"}, nil,
			`<html><h1 class="foo">Some text</h1>`,
		},
		{
			"basic-id", nil, []string{"bar"},
			`<html><h1 id="bar">Some text</h1>`,
		},
		{
			"multiple-classes-and-ids", []string{"foo", "bar"}, []string{"bar", "baz"},
			`
            <html>
                <body class="foo bar" id="bar">
                    <p id="baz" class="foo">
                </body>
            </html>
            `,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewBufferString(tc.html)
			classes, ids, err := ExtractClassesAndIDs(r)
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
