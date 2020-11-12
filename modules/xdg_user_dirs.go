package modules

import (
	"errors"
	"strings"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

var xdgAllUserDirs = []string{"DESKTOP", "DOWNLOAD", "TEMPLATES", "PUBLICSHARE", "DOCUMENTS", "MUSIC", "PICTURES", "VIDEOS"}

func init() {
	Register("xdg-user-dirs", XdgUserDirs)
}

func XdgUserDirs(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("XDG user dir")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for dir, path := range data {
		dir = strings.ToUpper(dir)
		if !helper.IsInStrings(dir, xdgAllUserDirs) {
			return errors.New("User dir \"" + dir + "\" does not exist")
		}
		path = helper.Path(path)
		if _, err := helper.Exec("xdg-user-dirs-update", "--set", dir, path); err != nil {
			return err
		}
		out.Success("Changed user dir " + dir + " to " + path)
	}
	return nil
}
