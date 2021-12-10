package commands

import (
	"strings"

	"gopkg.in/ini.v1"
)

const xdgMimetypeFile = "/etc/xdg/mimeapps.list"

func init() {
	register(xdgMimeDefault)
}

var xdgMimeDefault = Instruction{
	"xdg.mime-default",
	"Set the default application for a mimetype",
	"Do not change default application configuration",
	[]string{
		"Mimetype",
		"Name of the application",
	},
	nil,
	"Use Chromium\n  xdg.mime-default text/html chromium_chromium",
	func(args []string, out output, dry bool) ([]string, bool) {
		mimetype := args[0]
		app := args[1]

		ok := xdgCommonMimeDefault(xdgMimetypeFile, mimetype, app, out, dry)
		return nil, ok
	},
}

func xdgCommonMimeDefault(file, mimetype, app string, out output, dry bool) bool {

	if !strings.HasSuffix(".desktop", app) {
		app += ".desktop"
	}

	conf, err := ini.LooseLoad(file)
	if err != nil {
		out.Errorf("could not read %s: %v", file, err)
		return false
	}
	section := conf.Section("Default Applications")

	key := section.Key(mimetype)

	if app == key.Value() {
		out.Infof("default app for %s is already %s", mimetype, app)
		return true
	}

	if dry {
		out.Infof("would set default app for %s to %s", mimetype, app)
		return true
	}
	key.SetValue(app)
	if err := conf.SaveTo(file); err != nil {
		out.Errorf("could not save mimetypes file: %v", err)
		return false
	}
	return true
}
