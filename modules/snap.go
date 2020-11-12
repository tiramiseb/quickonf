package modules

import (
	"bytes"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("snap", Snap)
	Register("snap-classic", SnapClassic)
	Register("snap-dangerous", SnapDangerous)
	Register("snap-refresh", SnapRefresh)
	Register("snap-version", SnapVersion)
}

func Snap(in interface{}, out output.Output, store map[string]interface{}) error {
	return snap(in, out, false, false, store)
}

func SnapClassic(in interface{}, out output.Output, store map[string]interface{}) error {
	return snap(in, out, true, false, store)
}

func SnapDangerous(in interface{}, out output.Output, store map[string]interface{}) error {
	return snap(in, out, false, true, store)
}

func snap(in interface{}, out output.Output, classic bool, dangerous bool, store map[string]interface{}) error {
	out.ModuleName("Installing snap package")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		out.Info("Installing " + pkg)
		out.ShowLoader()
		var err error
		if classic {
			err = helper.ExecSudo("snap", "install", "--classic", pkg)
		} else if dangerous {
			err = helper.ExecSudo("snap", "install", "--dangerous", pkg)
		} else {
			err = helper.ExecSudo("snap", "install", pkg)
		}
		if err != nil {
			return err
		}
		out.HideLoader()
	}
	return nil
}

func SnapRefresh(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Refreshing snap packages")
	out.ShowLoader()
	err := helper.ExecSudo("snap", "refresh")
	if err != nil {
		return err
	}
	out.HideLoader()
	out.Success("Refreshed")
	return nil
}

func SnapVersion(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Getting snap package version")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for pkg, storekey := range data {
		output, err := helper.Exec("snap", "list", pkg)
		if err != nil {
			// Error probably means that package is not installed...
			// Package not installed is represented by 0.0
			store[storekey] = "0.0"
			continue
			// return err
		}
		lines := bytes.Split(output, []byte{'\n'})
		for _, l := range lines {
			line := bytes.Fields(l)
			if string(line[0]) == pkg {
				store[storekey] = string(line[1])
				break
			}
		}
	}
	return nil
}
