package instructions

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/commands"
)

type CookbookRecipe struct {
	Doc          string
	VarsDoc      map[string]string
	Instructions []Instruction
}

var recipes = map[string]CookbookRecipe{}

type Cookbook struct {
	URI    string
	ReadFn func(r io.Reader) (*Groups, error)
}

func (c *Cookbook) Name() string {
	return "cookbook"
}

func (c *Cookbook) RunCheck(vars *Variables, signalTarget chan bool, level int) ([]*CheckReport, bool) {
	uri := vars.TranslateVariables(c.URI)
	var reader io.Reader
	switch {
	case strings.HasPrefix(uri, "http://"), strings.HasPrefix(uri, "https://"):
		response, err := http.Get(uri)
		if err != nil {
			return []*CheckReport{{
				Name:         "cookbook",
				level:        level,
				status:       commands.StatusError,
				message:      fmt.Sprintf("Cannot download %s", uri),
				signalTarget: signalTarget,
			}}, false
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			return []*CheckReport{{
				Name:         "cookbook",
				level:        level,
				status:       commands.StatusError,
				message:      fmt.Sprintf(`Cannot download %s: %s`, uri, response.Status),
				signalTarget: signalTarget,
			}}, false
		}
		reader = response.Body
	case filepath.IsAbs(uri):
		f, err := os.Open(uri)
		if err != nil {
			return []*CheckReport{{
				Name:         "cookbook",
				level:        level,
				status:       commands.StatusError,
				message:      fmt.Sprintf("Cannot download open file %s: %s", uri, err.Error()),
				signalTarget: signalTarget,
			}}, false
		}
		defer f.Close()
		reader = f
	default:
		return []*CheckReport{{
			Name:         "cookbook",
			level:        level,
			status:       commands.StatusError,
			message:      fmt.Sprintf(`Cannot understand URI "%s" (supports local files (absolute path), HTTP and HTTPS)`, uri),
			signalTarget: signalTarget,
		}}, false
	}

	newRecipes, err := c.ReadFn(reader)

	reports := []*CheckReport{}
	if err != nil {
		reports = append(reports, &CheckReport{
			Name:         "cookbook",
			level:        level,
			status:       commands.StatusError,
			message:      fmt.Sprintf("Error in cookbook %s: %s", uri, err.Error()),
			signalTarget: signalTarget,
		})
	}

	for _, recipe := range newRecipes.groups {
		if _, ok := recipes[recipe.Name]; ok {
			reports = append(reports, &CheckReport{
				Name:         "cookbook",
				level:        level,
				status:       commands.StatusError,
				message:      fmt.Sprintf("Recipe %q is already defined", recipe.Name),
				signalTarget: signalTarget,
			})
		}
		recipes[recipe.Name] = CookbookRecipe{Instructions: recipe.Instructions}
	}

	reports = append(reports, &CheckReport{
		Name:         "cookbook",
		level:        level,
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf("Successfully read cookbook %s", uri),
		signalTarget: signalTarget,
	})
	return reports, true
}

func (c *Cookbook) NotRunReports(level int) []*CheckReport {
	return []*CheckReport{{
		Name:    "cookbook",
		level:   level,
		status:  commands.StatusNotRun,
		message: "cookbook " + c.URI,
	}}
}

func (c *Cookbook) Reset() {}

func (c *Cookbook) String() string {
	return c.indentedString(0)
}

func (c *Cookbook) indentedString(level int) string {
	var content stringBuilder
	content.add("cookbook")
	content.add(c.URI)
	return content.string(level)
}

func (c *Cookbook) hasConfigError() bool {
	return false
}
