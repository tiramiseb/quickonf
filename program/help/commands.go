package help

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

func (m *Model) commandsDoc(dark bool) string {
	var pattern string
	if dark {
		pattern = fmt.Sprintf("*%s*.dark.msg", m.commandFilter)
	} else {
		pattern = fmt.Sprintf("*%s*.light.msg", m.commandFilter)
	}

	var result strings.Builder

	if err := fs.WalkDir(commandsFS, "commands", func(filePath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := dirEntry.Name()
		if ok, err := path.Match(pattern, name); err != nil {
			return err
		} else if !ok {
			return nil
		}
		src, err := commandsFS.Open("commands/" + name)
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
