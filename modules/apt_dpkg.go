package modules

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("dpkg", Dpkg)
	Register("dpkg-version", DpkgVersion)
	Register("apt", Apt)
	Register("apt-remove", AptRemove)
	Register("apt-upgrade", AptUpgrade)
	Register("apt-autoremove-purge", AptAutoremovePurge)
}

func Dpkg(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Install .deb package")
	data, err := input.SliceString(in, store)
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

func DpkgVersion(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Checking package version")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	pkg, ok := data["package"]
	if !ok {
		return errors.New("Missing package name")
	}
	cmdout, _ := exec.Command("dpkg-query", "--showformat=${Version}", "--show", pkg).Output()
	if storeAs, ok := data["store"]; ok {
		store[storeAs] = string(cmdout)
	}
	if len(cmdout) == 0 {
		out.Info("Package " + pkg + " is not installed")
	} else {
		out.Info("Package " + pkg + " version is " + string(cmdout))
	}
	return nil
}

func Apt(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("APT install package")
	data, err := input.SliceString(in, store)
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

func AptRemove(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("APT remove package")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		cmdout, err := exec.Command("dpkg", "--get-selections", pkg).CombinedOutput()
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

func AptUpgrade(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("APT upgrade")
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

func AptAutoremovePurge(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("APT autoremove")
	out.Info("Removing unneeded dependencies")
	out.ShowLoader()
	err := helper.ExecSudo("apt-get", "--yes", "autoremove", "--purge")
	out.HideLoader()
	return err
}
