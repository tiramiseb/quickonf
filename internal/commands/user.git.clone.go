package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(userGitClone)
}

var userGitClone = Command{
	"user.git.clone",
	"Clone a git repository as provided user (if the repository already exists locally, pull the last commit) and switch to the requested reference (branch, tag...)",
	[]string{
		"Username",
		"Remote repository URI",
		"Clone target (if path is relative, it is relative to the user's nome directory)",
		"Reference (branch, tag...)",
	},
	nil,
	"Oh My Bash\n  user.git.clone alice https://github.com/ohmybash/oh-my-bash.git .oh-my-bash master",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		uri := args[1]
		dest := args[2]
		ref := args[3]

		usr, err := datastores.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}

		if !filepath.IsAbs(dest) {
			dest = filepath.Join(usr.User.HomeDir, dest)
		}

		lst, err := datastores.GitRemotes.List(uri)
		if err != nil {
			return nil, fmt.Sprintf("Could not list remote references for %s: %s", uri, err), nil, StatusError
		}
		var reference *plumbing.Reference
		for _, l := range lst {
			if l.Name().Short() == ref {
				reference = l
			}
		}
		if reference == nil {
			return nil, fmt.Sprintf(`Reference "%s" does not exist in %s`, ref, uri), nil, StatusError
		}

		// Check if destination already exists
		finfo, err := os.Stat(dest)
		if errors.Is(err, fs.ErrNotExist) {
			apply = &Apply{
				"user.git.clone",
				fmt.Sprintf("Will clone %s into %s", uri, dest),
				func(out Output) bool {
					out.Infof("Cloning %s into %s", uri, dest)
					if err := helper.ExecAs(usr.User, nil, nil, "git", "clone", uri, dest); err != nil {
						out.Errorf("Could not clone %s: %s", uri, err)
						return false
					}
					out.Infof("Checking out %s in %s", ref, dest)
					if err := helper.ExecAs(usr.User, nil, nil, "git", "-C", dest, "checkout", ref); err != nil {
						out.Errorf("Could not checkout %s in %s: %s", ref, dest, err)
						return false
					}
					out.Successf("Cloned %s into %s", uri, dest)
					return true
				},
			}
			return nil, fmt.Sprintf("Need to clone %s into %s", uri, dest), apply, StatusInfo
		}
		if err != nil {
			return nil, fmt.Sprintf("Could not check stats of %s: %s", dest, err), nil, StatusError
		}

		// Destination is not a directory, thus not a repository!
		if !finfo.IsDir() {
			return nil, fmt.Sprintf("%s is not a directory", dest), nil, StatusError
		}

		// Is destination a repository?
		repo, err := git.PlainOpen(dest)
		if err != nil {
			return nil, fmt.Sprintf("%s is not a repository: %s", dest, err), nil, StatusError
		}

		// Is the destination repository the one we want?
		remotes, err := repo.Remotes()
		if err != nil {
			return nil, fmt.Sprintf("Could not list remotes of %s: %s", dest, err), nil, StatusError
		}
		var isTheCorrectRepository bool
		for _, r := range remotes {
			for _, u := range r.Config().URLs {
				if u == uri {
					isTheCorrectRepository = true
				}
			}
		}

		if !isTheCorrectRepository {
			return nil, fmt.Sprintf("%s is not a clone of %s", dest, uri), nil, StatusError
		}

		apply = &Apply{
			"user.git.clone",
			fmt.Sprintf("Will pull updates in %s", dest),
			func(out Output) bool {
				out.Infof("Pulling in %s", dest)
				if err := helper.ExecAs(usr.User, nil, nil, "git", "-C", dest, "pull"); err != nil {
					out.Errorf("Could not pull in %s: %s", dest, err)
					return false
				}
				out.Infof("Checking out %s in %s", ref, dest)
				if err := helper.ExecAs(usr.User, nil, nil, "git", "-C", dest, "checkout", ref); err != nil {
					out.Errorf("Could not checkout %s in %s: %s", ref, dest, err)
					return false
				}
				out.Successf("Pulled in %s", dest)
				return true
			},
		}

		return nil, fmt.Sprintf("Need to pull updates in %s", dest), apply, StatusInfo
	},
	nil,
}
