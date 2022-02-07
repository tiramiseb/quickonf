package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(dpkgInstall)
}

var dpkgInstall = Command{
	"dpkg.install",
	"Install a package using dpkg",
	[]string{"Absolute path to the package file"},
	nil,
	"Install that awesome package\n  dpkg.install <confdir>/my-awesome-stuff.deb",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		path := args[0]
		if !filepath.IsAbs(path) {
			return nil, fmt.Sprintf("%s is not an absolute path", path), nil, StatusError
		}
		_, err := os.Lstat(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("%s does not exist", path), nil, StatusError
			}
			return nil, err.Error(), nil, StatusError
		}

		var out bytes.Buffer
		if err := helper.Exec(nil, &out, "dpkg-deb", "-f", path, "Package"); err != nil {
			return nil, fmt.Sprintf("could not check package name of %s: %s", path, helper.ExecErr(err)), nil, StatusError
		}
		candidatename := out.String()
		out.Reset()
		if err := helper.Exec(nil, &out, "dpkg-deb", "-f", path, "Version"); err != nil {
			return nil, fmt.Sprintf("could not check version of %s: %s", path, helper.ExecErr(err)), nil, StatusError
		}
		candidateversion := out.String()

		pkg, ok, err := datastores.DpkgPackages.Get(candidatename)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}
		if ok && pkg.Version == candidateversion {
			return nil, fmt.Sprintf("%s is already installed in version %s", candidatename, candidateversion), nil, StatusSuccess
		}

		apply = &Apply{
			"dpkg.install",
			fmt.Sprintf("Will install %s", path),
			func(out Output) bool {
				out.Infof("Waiting for dpkg to be available to install %s", path)
				datastores.DpkgMutex.Lock()
				defer datastores.DpkgMutex.Unlock()
				out.Runningf("Installing %s", path)
				if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "dpkg", "--install", path); err != nil {
					out.Errorf("Could not install %s: %s", path, helper.ExecErr(err))
					return false
				}
				out.Successf("Installed %s", path)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to install %s", path), apply, StatusInfo
	},
	datastores.DpkgPackages.Reset,
}
