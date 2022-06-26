package commands

import (
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"gopkg.in/ini.v1"
)

const xdgMimetypeFile = "/etc/xdg/mimeapps.list"

func init() {
	register(xdgMimeDefault)
}

var xdgMimeDefault = &Command{
	"xdg.mimedefault",
	"Set the default application for a mimetype",
	[]string{
		"Mimetype",
		"Name of the application",
	},
	nil,
	"Use Chromium\n  xdg.mimedefault text/html chromium_chromium",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		mimetype := args[0]
		app := args[1]
		if !strings.HasSuffix(".desktop", app) {
			app += ".desktop"
		}

		current, err := datastores.XdgMimetypes.Get(mimetype)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if current == app {
			return nil, fmt.Sprintf("Default app for %s is already %s", mimetype, app), nil, StatusSuccess, current, app
		}

		apply = func(out Output) bool {
			out.Runningf("Setting default app for %s to %s", mimetype, app)
			conf, err := ini.LooseLoad(xdgMimetypeFile)
			if err != nil {
				out.Errorf("Could not load mimetypes file: %v", err)
				return false
			}
			conf.Section("Default Applications").Key(mimetype).SetValue(app)
			if err := conf.SaveTo(xdgMimetypeFile); err != nil {
				out.Errorf("Could not save mimetypes file: %v", err)
				return false
			}
			out.Successf("Default app for %s set to %s", mimetype, app)
			return true
		}
		return nil, fmt.Sprintf("Need to set default app for %s to %s", mimetype, app), apply, StatusInfo, current, app
	},
	datastores.XdgMimetypes.Reset,
}
