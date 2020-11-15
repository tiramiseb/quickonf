package modules

import (
	// 	"bytes"
	"bytes"
	"errors"
	"os/exec"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("dpkg", Dpkg)
	Register("dpkg-version", DpkgVersion)
	Register("apt", Apt)
	Register("apt-remove", AptRemove)
	Register("apt-upgrade", AptUpgrade)
	Register("apt-autoremove-purge", AptAutoremovePurge)
}

// Dpkg installs a .deb package
func Dpkg(in interface{}, out output.Output) error {
	out.InstructionTitle("Install .deb package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, path := range data {
		out.Info("Installing " + path)
		out.ShowLoader()
		err := helper.ExecSudo("dpkg", "--install", path)
		if err != nil {
			return err
		}
		out.HideLoader()
	}
	return nil
}

// DpkgVersion checks an installed deb package version
func DpkgVersion(in interface{}, out output.Output) error {
	out.InstructionTitle("Checking package version")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	pkg, ok := data["package"]
	if !ok {
		return errors.New("Missing package name")
	}
	cmdout, err := helper.Exec("dpkg-query", "--showformat=${Version}", "--show", pkg)
	if err != nil {
		out.Info("Package " + pkg + " is not installed")
		if storeAs, ok := data["store"]; ok {
			helper.Store(storeAs, "0.0.0")
		}
	} else {
		out.Info("Package " + pkg + " version is " + string(cmdout))
		if storeAs, ok := data["store"]; ok {
			helper.Store(storeAs, string(cmdout))
		}
	}
	return nil
}

// Apt installs packages from apt repositories
func Apt(in interface{}, out output.Output) error {
	out.InstructionTitle("APT install package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}

	for _, pkg := range data {
		cmdout, err := exec.Command("dpkg", "--get-selections", pkg).CombinedOutput()
		if err != nil {
			return errors.New(string(cmdout))
		}
		if bytes.Index(cmdout, []byte("install")) >= 0 {
			out.Info(pkg + " is already installed")
			continue
		}
		out.Info("Installing " + pkg)
		out.ShowLoader()
		err = helper.ExecSudo("apt-get", "--yes", "--quiet", "install", "--no-install-recommends", pkg)
		out.HideLoader()
		if err != nil {
			return err
		}
	}
	return nil
}

// AptRemove removes deb packages from system
func AptRemove(in interface{}, out output.Output) error {
	out.InstructionTitle("APT remove package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		cmdout, err := helper.Exec("dpkg", "--get-selections", pkg)
		if err != nil {
			return errors.New(string(cmdout))
		}
		if bytes.Index(cmdout, []byte("install")) == -1 {
			out.Info(pkg + " is not installed")
			continue
		}
		out.Info("Removing " + pkg)
		out.ShowLoader()
		err = helper.ExecSudo("apt-get", "--yes", "--quiet", "remove", pkg)
		out.HideLoader()
		if err != nil {
			return err
		}
	}
	return nil
}

// AptUpgrade upgrades all packages in the system
func AptUpgrade(in interface{}, out output.Output) error {
	out.InstructionTitle("APT upgrade")
	out.Info("Updating packages list")
	out.ShowLoader()
	err := helper.ExecSudo("apt-get", "--yes", "update")
	out.HideLoader()
	if err != nil {
		return err
	}
	out.Info("Upgrading packages")
	out.ShowLoader()
	err = helper.ExecSudo("apt-get", "--yes", "upgrade")
	out.HideLoader()
	if err != nil {
		return err
	}
	return nil
}

// AptAutoremovePurge cleans unneeded packages from the system
func AptAutoremovePurge(in interface{}, out output.Output) error {
	out.InstructionTitle("APT autoremove")
	out.Info("Removing unneeded dependencies")
	out.ShowLoader()
	err := helper.ExecSudo("apt-get", "--yes", "autoremove", "--purge")
	out.HideLoader()
	return err
}
