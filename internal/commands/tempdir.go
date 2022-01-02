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
	"Create a temporary directory (directory is NOT deleted after usage, you must use file.absent to remove it at the end of the group)",
	nil,
	[]string{
		"Available temporary path",
	},
	"Temporarily clone git repository\n  tmp = temppath\n  git.clone https://www.example.com/foobar.git <tmp>",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		path := filepath.Join(os.TempDir(), fmt.Sprintf("quickonf-%d", rand.Int()))

		apply = &Apply{
			"tempdir",
			"Will create a temporary directory",
			func(out Output) bool {
				out.Info("Creating temporary directory")
				if err := os.MkdirAll(path, 0644); err != nil {
					out.Errorf("Could not create temporary directory: %s", err)
					return false
				}
				out.Success("Created temporary directory")
				return true
			},
		}
		return []string{path}, "Need to create a temporary directory", apply, StatusInfo
	},
	nil,
}
