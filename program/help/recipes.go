package help

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

func (m *Model) recipesDoc(dark bool) string {
	var pattern string
	if dark {
		pattern = fmt.Sprintf("*%s*.dark.msg", m.recipeFilter)
	} else {
		pattern = fmt.Sprintf("*%s*.light.msg", m.recipeFilter)
	}

	var result strings.Builder

	if err := fs.WalkDir(recipesFS, "content/cookbook", func(filePath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := strings.ToLower(dirEntry.Name())
		if ok, err := path.Match(pattern, name); err != nil {
			return err
		} else if !ok {
			return nil
		}
		src, err := recipesFS.Open("content/cookbook/" + name)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(&result, src)
		return err
	}); err != nil {
		return "Could not render documentation: " + err.Error() + "\n-----\n" + result.String()
	}
	return result.String()
}
