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

func setFrom(ss []string) map[string]bool {
	set := make(map[string]bool, len(ss))
	for _, s := range ss {
		set[s] = true
	}
	return set
}

func diff(a, b map[string]bool) []string {
	var ss []string
	for key := range a {
		if _, ok := b[key]; !ok {
			ss = append(ss, key)
		}
	}
	return ss
}

func equal(t *testing.T, msg string, ss []string, set map[string]bool) {
	t.Helper()

	newSet := setFrom(ss)
	diff1 := diff(newSet, set)
	diff2 := diff(set, newSet)
	if len(diff1) != 0 || len(diff2) != 0 {
		t.Fatalf("%s: has %v; missing %v", msg, diff1, diff2)
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
		{
			"qualified-case", []string{
				"cand-name",
				"candidate-box",
				"candidates-list",
				"center",
				"container",
				"desktop",
				"nav-item",
				"question-candidates",
				"selected",
				"span-70",
			}, []string{
				"all-candidates",
				"cand-name",
				"donut-container",
				"inline-nav",
				"mobile-nav",
				"star-and-logo",
				"top-nav",
				"top-sun-logo",
			},
			`
@media only screen and (max-width: 1000px) {
    /* This makes the running against list on the candidates page one col */
    div.candidates-list.span-70.center {
        clear: both;
    }
    /* Footer height needs to be increased at smaller sizes, as more links are added */
    footer {
        height: 800px;
    }
    .container,
    #inline-nav {
        width: 90%;
        margin-left: 5%;
    }
    nav {
        display: none;
    }
    #cand-name,
    .cand-name {
        font-size: 2em;
        width: auto;
    }
    #all-candidates ul li {
        padding: 0px 10px;
    }
    .question-candidates ul li.selected:before {
        right: 156px;
    }
    .nav-item {
        margin: 0px 5%;
    }
    #all-candidates {
        margin: 20px 0px 30px 0px;
    }
    .candidates-list {
        width: auto;
    }
    #top-nav,
    #top-sun-logo {
        display: none;
    }
    #star-and-logo {
        margin-left: 5px;
    }
    #mobile-nav {
        display: inline-block;
    }
    #donut-container,
    .candidate-box.selected:after {
        display: none;
    }
    #inline-nav.desktop {
        display: none;
    }
}
			`,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewBufferString(tc.css)
			classes, ids, err := ExtractClassesAndIDs(r)
			check(t, err)
			equal(t, "classes not equal", tc.classes, classes)
			equal(t, "ids not equal", tc.ids, ids)
		})
	}
}
