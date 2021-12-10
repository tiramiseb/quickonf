package commands

import (
	"bytes"
	"sync"

	"github.com/tiramiseb/quickonf/internal/helper"
)

func init() {
	register(apt)
}

var aptMutex sync.Mutex

var apt = Command{
	"apt",
	"Install a package using apt",
	"Do not install the package",
	[]string{"Name of the package to install"},
	nil,
	"Install the \"ipcalc\" tool\n  apt ipcalc",
	func(args []string, out output, dry bool) ([]string, bool) {
		pkg := args[0]
		var buf bytes.Buffer
		wait, err := helper.Exec(nil, &buf, "dpkg", "--get-selections", pkg)
		if err != nil {
			out.Errorf("could not check if %s is installed: %s", pkg, err)
			return nil, false
		}
		if err := wait(); err != nil {
			out.Errorf("could not check if %s is installed: %s", pkg, err)
			return nil, false
		}
		if bytes.Contains(buf.Bytes(), []byte("install")) {
			out.Successf("%s is already installed", pkg)
			return nil, true
		}
		if dry {
			out.Successf("would install %s", pkg)
			return nil, true
		}
		out.Infof("waiting for apt to be available to install %s", pkg)
		aptMutex.Lock()
		defer aptMutex.Unlock()
		wait, err = helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "--quiet", "install", pkg)
		if err != nil {
			out.Errorf("could not install %s: %s", pkg, err)
			return nil, false
		}
		out.Infof("installing %s", pkg)
		if err := wait(); err != nil {
			out.Errorf("could not install %s: %s", pkg, err)
			return nil, false
		}
		out.Successf("installed %s", pkg)
		return nil, true
	},
}
