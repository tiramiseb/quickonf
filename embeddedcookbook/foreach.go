package embeddedcookbook

import (
	"embed"
	"io"

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
	groups, err := conf.Read(r)
	if err != nil {
		return err
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
