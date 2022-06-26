package embeddedcookbook

import (
	"embed"
	"errors"
	"io"
	"strings"

	"github.com/tiramiseb/quickonf/conf"
	"github.com/tiramiseb/quickonf/instructions"
)

//go:embed *.qconf
var cookbook embed.FS

func ForEach(fn func(*instructions.Group) error) error {
	files, err := cookbook.ReadDir(".")
	if err != nil {
		return err
	}
	for _, fsEntry := range files {
		f, err := cookbook.Open(fsEntry.Name())
		if err != nil {
			return err
		}
		if err := forEach(f, fn); err != nil {
			return err
		}
	}
	return nil
}

func forEach(r io.Reader, fn func(*instructions.Group) error) error {
	groups, errs := conf.Read(r)
	if len(errs) > 0 {
		errmsg := make([]string, len(errs))
		for i, err := range errs {
			errmsg[i] = err.Error()
		}
		return errors.New(strings.Join(errmsg, "\n"))
	}

	group := groups.FirstGroup()
	for {
		if err := fn(group); err != nil {
			return err
		}
		newGrp := group.Next(1, true)
		if newGrp == group {
			break
		}
		group = newGrp
	}
	return nil
}
