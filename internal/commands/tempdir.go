package commands

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

func init() {
	register(tempdir)
}

var tempdir = Command{
	"tempdir",
	"Create a temporary directory (directory is automatically deleted when closing the application)",
	nil,
	[]string{
		"Path of the created directory",
	},
	"Temporarily clone git repository\n  tmp = tempdir\n  git.clone https://www.example.com/foobar.git <tmp>",
	func(args []string) (result []string, msg string, apply Apply, status Status) {
		path := filepath.Join(os.TempDir(), fmt.Sprintf("quickonf-%d", rand.Int()))

		apply = func(out Output) bool {
			out.Info("Creating temporary directory")
			if err := os.MkdirAll(path, 0o644); err != nil {
				out.Errorf("Could not create temporary directory: %s", err)
				return false
			}
			registerClean(func() error {
				return os.RemoveAll(path)
			})
			out.Success("Created temporary directory")
			return true
		}
		return []string{path}, "Need to create a temporary directory", apply, StatusInfo
	},
	nil,
}
