package html

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func ExtractClassesAndIDs(r io.Reader) (classes, ids map[string]bool, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, nil, err
	}

	classes = map[string]bool{}
	ids = map[string]bool{}
	findClassesAndIDs(doc, classes, ids)
	return
}

func findClassesAndIDs(n *html.Node, classes, ids map[string]bool) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" {
				vals := strings.Fields(a.Val)
				for _, val := range vals {
					classes[val] = true
				}
			} else if a.Key == "id" {
				ids[a.Val] = true
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findClassesAndIDs(c, classes, ids)
	}
}
