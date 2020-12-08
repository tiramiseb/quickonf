package modules

import (
	"bytes"
	"errors"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("snap", Snap)
	Register("snap-classic", SnapClassic)
	Register("snap-dangerous", SnapDangerous)
	Register("snap-refresh", SnapRefresh)
	Register("snap-version", SnapVersion)
}

// Snap installs Snap packages
func Snap(in interface{}, out output.Output) error {
	return snap(in, out, false, false)
}

// SnapClassic installs Snap packages in classic mode
func SnapClassic(in interface{}, out output.Output) error {
	return snap(in, out, true, false)
}

// SnapDangerous installs Snap packages without verifying their signatures
func SnapDangerous(in interface{}, out output.Output) error {
	return snap(in, out, false, true)
}

func snap(in interface{}, out output.Output, classic bool, dangerous bool) error {
	out.InstructionTitle("Installing snap package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		if Dryrun {
			out.Info("Would install " + pkg)
			continue
		}
		out.Info("Installing " + pkg)
		out.ShowLoader()
		var err error
		if classic {
			_, err = helper.ExecSudo("snap", "install", "--classic", pkg)
		} else if dangerous {
			_, err = helper.ExecSudo("snap", "install", "--dangerous", pkg)
		} else {
			_, err = helper.ExecSudo("snap", "install", pkg)
		}
		if err != nil {
			return err
		}
		out.HideLoader()
	}
	return nil
}

// SnapRefresh refreshes Snap packages
func SnapRefresh(in interface{}, out output.Output) error {
	out.InstructionTitle("Refreshing snap packages")
	if Dryrun {
		out.Info("Would refresh packages")
		return nil
	}
	out.ShowLoader()
	_, err := helper.ExecSudo("snap", "refresh")
	out.HideLoader()
	if err != nil {
		return err
	}
	out.Success("Refreshed")
	return nil
}

// SnapVersion gets a Snap package version
func SnapVersion(in interface{}, out output.Output) error {
	out.InstructionTitle("Getting snap package version")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	pkg, ok := data["package"]
	if !ok {
		return errors.New("Missing package name")
	}
	cmdout, err := helper.Exec("snap", "info", pkg)
	if err != nil {
		out.Info("Package " + pkg + " is not installed")
		if storeAs, ok := data["store"]; ok {
			helper.Store(storeAs, "")
		}
	}
	for _, l := range bytes.Split(cmdout, []byte{'\n'}) {
		fields := bytes.Fields(l)
		if bytes.Equal(fields[0], []byte("installed:")) {
			out.Info("Package " + pkg + " version is " + string(fields[1]))
			if storeAs, ok := data["store"]; ok {
				helper.Store(storeAs, string(fields[1]))
			}
			return nil
		}
	}
	out.Info("Could not determine package " + pkg + " version")
	if storeAs, ok := data["store"]; ok {
		helper.Store(storeAs, "")
	}
	return nil
}
