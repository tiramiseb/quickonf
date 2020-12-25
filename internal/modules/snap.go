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
	Register("snap-edge", SnapEdge)
	Register("snap-refresh", SnapRefresh)
	Register("snap-version", SnapVersion)
}

// Snap installs Snap packages
func Snap(in interface{}, out output.Output) error {
	return snap(in, out)
}

// SnapClassic installs Snap packages in classic mode
func SnapClassic(in interface{}, out output.Output) error {
	return snap(in, out, "--classic")
}

// SnapDangerous installs Snap packages without verifying their signatures
func SnapDangerous(in interface{}, out output.Output) error {
	return snap(in, out, "--dangerous")
}

// SnapEdge installs Snap packages from the edge channel
func SnapEdge(in interface{}, out output.Output) error {
	return snap(in, out, "--edge")
}

func snap(in interface{}, out output.Output, options ...string) error {
	out.InstructionTitle("Installing snap package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		if Dryrun {
			out.Infof("Would install %s", pkg)
			continue
		}
		out.Infof("Installing %s", pkg)
		out.ShowLoader()
		var err error
		cmd := []string{"snap", "install"}
		cmd = append(cmd, options...)
		cmd = append(cmd, pkg)
		if _, err = helper.ExecSudo(nil, cmd...); err != nil {
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
	_, err := helper.ExecSudo(nil, "snap", "refresh")
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
	cmdout, err := helper.Exec(nil, "snap", "info", pkg)
	if err != nil {
		out.Infof("Package %s is not installed", pkg)
		if storeAs, ok := data["store"]; ok {
			helper.Store(storeAs, "")
		}
	}
	for _, l := range bytes.Split(cmdout, []byte{'\n'}) {
		fields := bytes.Fields(l)
		if bytes.Equal(fields[0], []byte("installed:")) {
			out.Infof("Package %s version is %s", pkg, fields[1])
			if storeAs, ok := data["store"]; ok {
				helper.Store(storeAs, string(fields[1]))
			}
			return nil
		}
	}
	out.Infof("Could not determine package %s version", pkg)
	if storeAs, ok := data["store"]; ok {
		helper.Store(storeAs, "")
	}
	return nil
}
