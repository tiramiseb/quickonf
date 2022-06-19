package conf

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) parseCookbook(uri *token) {
	variables := instructions.NewVariablesSet()
	target := variables.TranslateVariables(uri.content)
	var reader io.Reader
	switch {
	case strings.HasPrefix(target, "http://"), strings.HasPrefix(target, "https://"):
		response, err := http.Get(target)
		if err != nil {
			p.errs = append(p.errs, uri.errorf(`cannot download %s: %s`, target, err.Error()))
		}
		defer response.Body.Close()
		reader = response.Body
	case filepath.IsAbs(target):
		f, err := os.Open(target)
		if err != nil {
			p.errs = append(p.errs, uri.errorf(`cannot open file %s: %s`, target, err.Error()))
		}
		defer f.Close()
		reader = f
	default:
		p.errs = append(p.errs, uri.errorf(`cannot understand URI "%s" (supports local files (absolute path), HTTP and HTTPS)`, target))
		return
	}
	recipes, errs := Read(reader)
	for _, err := range errs {
		p.errs = append(p.errs, uri.errorf(`error in cookbook "%s": %s`, target, err.Error()))
	}
	for _, grp := range recipes.All() {
		if _, ok := p.recipes[grp.Name]; ok {
			p.errs = append(p.errs, uri.errorf(`recipe "%s" found in cookbook "%s" is already defined`, grp.Name, target))
			continue
		}
		p.recipes[grp.Name] = grp.Instructions
	}
}
