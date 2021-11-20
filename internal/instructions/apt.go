package instructions

import (
	"bytes"
	"time"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	register("apt", Apt, 1, 0)
}

// Apt installs a package from apt repositories
func Apt(args []string, out *output.Instruction) ([]string, bool) {
	pkg := args[0]
	var buf bytes.Buffer
	wait, err := helper.Exec(nil, &buf, "dpkg", "--get-selections", pkg)
	if err != nil {
		out.Errorf("cannot check if %s is installed: %s", pkg, err)
		return nil, false
	}
	if err := wait(); err != nil {
		out.Errorf("cannot check if %s is installed: %s", pkg, err)
		return nil, false
	}
	time.Sleep(4 * time.Second)
	if bytes.Contains(buf.Bytes(), []byte("install")) {
		out.Successf("%s is already installed", pkg)
		return nil, true
	}
	if Dryrun {
		out.Successf("would install %s", pkg)
		return nil, true
	}
	wait, err = helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "--quiet", "install", pkg)
	if err != nil {
		out.Errorf("cannot install %s: %s", pkg, err)
		return nil, false
	}
	out.Loadf("installing %s", pkg)
	if err := wait(); err != nil {
		out.Errorf("cannot install %s: %s", pkg, err)
		return nil, false
	}
	out.Successf("installed %s", pkg)
	return nil, true
}
