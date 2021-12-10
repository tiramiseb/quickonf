package commands

import (
	"os/user"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/helper"
)

const xdgUserMimetypeFile = "/etc/xdg/mimeapps.list"

func init() {
	register(xdgUserMimeDefault)
}

var xdgUserMimeDefault = Command{
	"xdg.user.mime-default",
	"Set the default application for a mimetype, for a specific user",
	"Do not change default application configuration",
	[]string{
		"Username",
		"Mimetype",
		"Name of the application",
	},
	nil,
	"Use Chromium\n  xdg.user.mime-default john text/html chromium_chromium",
	func(args []string, out output, dry bool) ([]string, bool) {
		username := args[0]
		mimetype := args[1]
		app := args[2]

		usr, err := user.Lookup(username)
		if err != nil {
			out.Errorf("could not identify user %s: %v", username, err)
			return nil, false
		}

		theFile := filepath.Join(usr.HomeDir, xdgUserMimetypeFile)

		ok := xdgCommonMimeDefault(theFile, mimetype, app, out, dry)
		if !ok {
			return nil, false
		}
		if !dry {
			if err := helper.Chown(theFile, usr); err != nil {
				out.Error(err.Error())
				return nil, false
			}
		}
		return nil, true
	},
}
