package styles

import (
	"fmt"
	"io"
	"strings"

	"github.com/tdewolff/parse/css"
)

func ExtractClassesAndIDs(r io.Reader) (classes, ids map[string]bool, err error) {
	p := css.NewParser(r, false)
	classes = map[string]bool{}
	ids = map[string]bool{}
	for {
		gt, _, _ := p.Next()
		if gt == css.ErrorGrammar {
			err = p.Err()
			if err == io.EOF {
				err = nil
			}
			if err != nil {
				err = fmt.Errorf("encountered error parsing CSS: %v", err)
			}
			return
		} else if gt == css.BeginRulesetGrammar || gt == css.QualifiedRuleGrammar {
			captureClass := false
			for _, val := range p.Values() {
				data := string(val.Data)
				if captureClass {
					captureClass = false
					classes[data] = true
				} else if data == "." {
					captureClass = true
				} else if strings.HasPrefix(data, "#") {
					ids[data[1:]] = true
				}
			}
		}
	}

}
