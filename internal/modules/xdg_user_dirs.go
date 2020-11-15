package modules

import (
	"errors"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

var xdgAllUserDirs = map[string]bool{
	"DESKTOP":     true,
	"DOCUMENTS":   true,
	"DOWNLOAD":    true,
	"MUSIC":       true,
	"PICTURES":    true,
	"PUBLICSHARE": true,
	"TEMPLATES":   true,
	"VIDEOS":      true,
}

func init() {
	Register("xdg-user-dir", XdgUserDir)
}

// XdgUserDir sets XDG user dirs
func XdgUserDir(in interface{}, out output.Output) error {
	out.InstructionTitle("XDG user dir")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for name, path := range data {
		name = strings.ToUpper(name)
		if !xdgAllUserDirs[name] {
			return errors.New("User dir \"" + name + "\" does not exist")
		}
		path = helper.Path(path)
		if _, err := helper.Exec("xdg-user-dirs-update", "--set", name, path); err != nil {
			return err
		}
		out.Success("Changed user dir " + name + " to " + path)
	}
	return nil
}
