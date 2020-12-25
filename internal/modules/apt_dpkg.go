package modules

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("dpkg", Dpkg)
	Register("dpkg-dependencies", DpkgDependencies)
	Register("dpkg-reconfigure", DpkgReconfigure)
	Register("dpkg-version", DpkgVersion)
	Register("debconf-set", DebconfSet)
	Register("apt", Apt)
	Register("apt-remove", AptRemove)
	Register("apt-upgrade", AptUpgrade)
	Register("apt-autoremove-purge", AptAutoremovePurge)
	Register("apt-clean-cache", AptCleanCache)
}

// Dpkg installs a .deb package
func Dpkg(in interface{}, out output.Output) error {
	out.InstructionTitle("Install .deb package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, path := range data {
		if Dryrun {
			out.Infof("Would install %s", path)
			continue
		}
		out.Infof("Installing %s", path)
		out.ShowLoader()
		_, err := helper.ExecSudo(nil, "dpkg", "--install", path)
		out.HideLoader()
		if err != nil {
			return err
		}
	}
	return nil
}

// DpkgDependencies installs dependencies for a .deb package
func DpkgDependencies(in interface{}, out output.Output) error {
	out.InstructionTitle("Install dependencies for .deb package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, path := range data {
		out.Infof("Dependencies for %s", path)
		depsB, err := helper.Exec(nil, "dpkg-deb", "--show", "--showformat=${Depends}", path)
		if err != nil {
			return err
		}
		deps := strings.Split(string(depsB), ", ")
		for _, pkg := range deps {
			if Dryrun {
				out.Infof("Would install %s", pkg)
				continue
			}
			out.Infof("Installing %s", pkg)
			out.ShowLoader()
			_, err = helper.ExecSudo([]string{"DEBIAN_FRONTEND=noninteractive"}, "apt-get", "--yes", "--quiet", "install", pkg)
			out.HideLoader()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DpkgReconfigure reconfigures a package interactively
func DpkgReconfigure(in interface{}, out output.Output) error {
	out.InstructionTitle("Reconfiguring deb package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, name := range data {
		if Dryrun {
			out.Infof("Would reconfigure %s", name)
			continue
		}
		out.Infof("Reconfiguring %s", name)
		out.ShowLoader()
		_, err := helper.ExecSudo(nil, "dpkg-reconfigure", "--frontend", "noninteractive", name)
		out.HideLoader()
		if err != nil {
			return err
		}
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
	cmdout, err := helper.Exec(nil, "dpkg-query", "--showformat=${Version}", "--show", pkg)
	if err != nil {
		out.Infof("Package %s is not installed", pkg)
		if storeAs, ok := data["store"]; ok {
			helper.Store(storeAs, "0.0.0")
		}
	} else {
		out.Infof("Package %s version is %s", pkg, cmdout)
		if storeAs, ok := data["store"]; ok {
			helper.Store(storeAs, string(cmdout))
		}
	}
	return nil
}

// DebconfSet sets a debconf variable
func DebconfSet(in interface{}, out output.Output) error {
	out.InstructionTitle("Set debconf variable")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	pkg, ok := data["package"]
	if !ok {
		return errors.New("Missing package name")
	}
	variable, ok := data["variable"]
	if !ok {
		return errors.New("Missing variable name")
	}
	value, ok := data["value"]
	if !ok {
		return errors.New("Missing value")
	}
	if Dryrun {
		out.Infof("Would set %s to %s for %s", variable, value, pkg)
		return nil
	}
	tmpfile, err := ioutil.TempFile("", "quickonf-debconf")
	if err != nil {
		return err
	}
	defer tmpfile.Close()
	if _, err := tmpfile.WriteString(strings.Join([]string{pkg, variable, "select", value}, " ")); err != nil {
		return err
	}
	if _, err := helper.ExecSudo(nil, "debconf-set-selections", tmpfile.Name()); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}
	return os.Remove(tmpfile.Name())
}

// Apt installs packages from apt repositories
func Apt(in interface{}, out output.Output) error {
	out.InstructionTitle("APT install package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}

	for _, pkg := range data {
		cmdout, err := helper.Exec(nil, "dpkg", "--get-selections", pkg)
		if err != nil {
			return errors.New(string(cmdout))
		}
		if bytes.Contains(cmdout, []byte("install")) {
			out.Infof("%s is already installed")
			continue
		}
		if Dryrun {
			out.Infof("Would install %s", pkg)
			continue
		}
		out.Infof("Installing %s", pkg)
		out.ShowLoader()
		_, err = helper.ExecSudo([]string{"DEBIAN_FRONTEND=noninteractive"}, "apt-get", "--yes", "--quiet", "install", pkg)
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
		cmdout, err := helper.Exec(nil, "dpkg", "--get-selections", pkg)
		if err != nil {
			return err
		}
		if !bytes.Contains(cmdout, []byte("install")) {
			out.Infof("%s is not installed", pkg)
			continue
		}
		if Dryrun {
			out.Infof("Would remove %s", pkg)
			continue
		}
		out.Infof("Removing %s", pkg)
		out.ShowLoader()
		_, err = helper.ExecSudo([]string{"DEBIAN_FRONTEND=noninteractive"}, "apt-get", "--yes", "--quiet", "remove", pkg)
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
	if Dryrun {
		out.Info("Would update packages list and upgrade packages")
		return nil
	}
	out.Info("Updating packages list")
	out.ShowLoader()
	_, err := helper.ExecSudo([]string{"DEBIAN_FRONTEND=noninteractive"}, "apt-get", "--yes", "update")
	out.HideLoader()
	if err != nil {
		return err
	}
	out.Info("Upgrading packages")
	out.ShowLoader()
	_, err = helper.ExecSudo([]string{"DEBIAN_FRONTEND=noninteractive"}, "apt-get", "--yes", "upgrade")
	out.HideLoader()
	if err != nil {
		return err
	}
	return nil
}

// AptAutoremovePurge cleans unneeded packages from the system
func AptAutoremovePurge(in interface{}, out output.Output) error {
	out.InstructionTitle("APT autoremove")
	if Dryrun {
		out.Info("Would clean the system from unneeded packages")
		return nil
	}
	out.Info("Removing unneeded dependencies")
	out.ShowLoader()
	_, err := helper.ExecSudo([]string{"DEBIAN_FRONTEND=noninteractive"}, "apt-get", "--yes", "autoremove", "--purge")
	out.HideLoader()
	return err
}

// AptCleanCache remove archives .deb files from cache
func AptCleanCache(in interface{}, out output.Output) error {
	out.InstructionTitle("Cleaning APT cache")
	if Dryrun {
		out.Info("Would clean the APT cache")
		return nil
	}
	out.Info("Cleaning the APT cache")
	out.ShowLoader()
	_, err := helper.ExecSudo([]string{"DEBIAN_FRONTEND=noninteractive"}, "apt-get", "clean")
	out.HideLoader()
	return err
}
