package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(gitClone)
}

var gitClone = &Command{
	"git.clone",
	"Clone a git repository (if the repository already exists locally, pull the last commit) and switch to the requested reference (branch, tag...)",
	[]string{
		"Remote repository URI",
		"Clone target (absolute path)",
		"Reference (branch, tag...)",
	},
	nil,
	"Temporarily clone git repository\n  tmp = temppath\n  git.clone https://www.example.com/foobar.git <tmp> main",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uri := args[0]
		dest := args[1]
		ref := args[2]

		if !filepath.IsAbs(dest) {
			return nil, fmt.Sprintf("%s is not an absolute path", dest), nil, StatusError, "", ""
		}

		// Check if destination already exists
		finfo, err := os.Stat(dest)
		if errors.Is(err, fs.ErrNotExist) {
			apply = func(out Output) bool {
				out.Runningf("Cloning %s into %s", uri, dest)
				if err := helper.Exec(nil, nil, "git", "clone", "--branch", ref, "--single-branch", uri, dest); err != nil {
					out.Errorf("Could not clone %s: %s", uri, helper.ExecErr(err))
					return false
				}
				out.Successf("Cloned %s into %s", uri, dest)
				return true
			}
			return nil, fmt.Sprintf("Need to clone %s into %s", uri, dest), apply, StatusInfo, "", ""
		}
		if err != nil {
			return nil, fmt.Sprintf("Could not check stats of %s: %s", dest, err), nil, StatusError, "", ""
		}

		// Destination is not a directory, thus not a repository!
		if !finfo.IsDir() {
			return nil, fmt.Sprintf("%s is not a directory", dest), nil, StatusError, "", ""
		}

		// Is destination a repository?
		info, err := os.Stat(filepath.Join(dest, ".git"))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Sprintf("%s is not a repository: %s", dest, err), nil, StatusError, "", ""
			}
			return nil, fmt.Sprintf("Could not check if %s is a repository: %s", dest, err), nil, StatusError, "", ""
		}
		if !info.IsDir() {
			return nil, fmt.Sprintf("%s is not a repository", dest), nil, StatusError, "", ""
		}

		// Is the destination repository the one we want?
		var out strings.Builder
		if err := helper.Exec(nil, &out, "git", "-C", dest, "remote", "--verbose"); err != nil {
			return nil, fmt.Sprintf("Could not list remotes of %s: %s", dest, err), nil, StatusError, "", ""
		}
		if !strings.Contains(out.String(), "\t"+uri+" ") {
			return nil, fmt.Sprintf("%s is not a clone of %s", dest, uri), nil, StatusError, "", ""
		}

		apply = func(out Output) bool {
			out.Infof("Pulling in %s", dest)
			if err := helper.Exec(nil, nil, "git", "-C", dest, "pull"); err != nil {
				out.Errorf("Could not pull in %s: %s", dest, err)
				return false
			}
			out.Infof("Checking out %s in %s", ref, dest)
			if err := helper.Exec(nil, nil, "git", "-C", dest, "checkout", ref); err != nil {
				out.Errorf("Could not checkout %s in %s: %s", ref, dest, err)
				return false
			}
			out.Successf("Pulled in %s", dest)
			return true
		}

		return nil, fmt.Sprintf("Need to pull updates in %s", dest), apply, StatusInfo, "", ""
	},
	nil,
}
