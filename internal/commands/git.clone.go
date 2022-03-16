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
)

func init() {
	register(gitClone)
}

var gitClone = Command{
	"git.clone",
	"Clone a git repository (if the repository already exists locally, pull the last commit) and switch to the requested reference (branch, tag...)",
	[]string{
		"Remote repository URI",
		"Clone target (absolute path)",
		"Reference (branch, tag...)",
	},
	nil,
	"Temporarily clone git repository\n  tmp = temppath\n  git.clone https://www.example.com/foobar.git <tmp>",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uri := args[0]
		dest := args[1]
		ref := args[2]

		if !filepath.IsAbs(dest) {
			return nil, fmt.Sprintf("%s is not an absolute path", dest), nil, StatusError, "", ""
		}

		lst, err := datastores.GitRemotes.List(uri)
		if err != nil {
			return nil, fmt.Sprintf("Could not list remote references for %s: %s", uri, err), nil, StatusError, "", ""
		}
		var reference *plumbing.Reference
		for _, l := range lst {
			if l.Name().Short() == ref {
				reference = l
			}
		}
		if reference == nil {
			return nil, fmt.Sprintf(`Reference "%s" does not exist in %s`, ref, uri), nil, StatusError, "", ""
		}

		// Check if destination already exists
		finfo, err := os.Stat(dest)
		if errors.Is(err, fs.ErrNotExist) {
			apply = func(out Output) bool {
				out.Runningf("Cloning %s into %s", uri, dest)
				repo, err := git.PlainClone(dest, false, &git.CloneOptions{
					URL:               uri,
					RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
				})
				if err != nil {
					out.Errorf("Could not clone %s: %s", uri, err)
					return false
				}
				worktree, err := repo.Worktree()
				if err != nil {
					out.Errorf("Could not work with %s: %s", dest, err)
					return false
				}
				out.Infof("Checking out %s in %s", ref, dest)
				if err := worktree.Checkout(&git.CheckoutOptions{
					Hash: plumbing.NewHash(ref),
				}); err != nil {
					out.Errorf("Could not checkout %s in %s: %s", ref, dest, err)
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
		repo, err := git.PlainOpen(dest)
		if err != nil {
			return nil, fmt.Sprintf("%s is not a repository: %s", dest, err), nil, StatusError, "", ""
		}

		// Is the destination repository the one we want?
		remotes, err := repo.Remotes()
		if err != nil {
			return nil, fmt.Sprintf("Could not list remotes of %s: %s", dest, err), nil, StatusError, "", ""
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
			return nil, fmt.Sprintf("%s is not a clone of %s", dest, uri), nil, StatusError, "", ""
		}

		// Get the working directory for the repository
		worktree, err := repo.Worktree()
		if err != nil {
			return nil, fmt.Sprintf("Could not work with %s: %s", dest, err), nil, StatusError, "", ""
		}

		apply = func(out Output) bool {
			out.Infof("Pulling in %s", dest)
			if err := worktree.Pull(&git.PullOptions{}); err != nil {
				out.Errorf("Could not pull in %s: %s", dest, err)
				return false
			}
			out.Infof("Checking out %s in %s", ref, dest)
			if err := worktree.Checkout(&git.CheckoutOptions{
				Hash: plumbing.NewHash(ref),
			}); err != nil {
				out.Errorf("Could not checkout %s in %s: %s", ref, dest, err)
				return false
			}
			out.Successf("Pulled in %s", dest)
			return true
		}

		return nil, fmt.Sprintf("Need to pull updates in %s", dest), apply, StatusInfo, "", ""
	},
	datastores.GitRemotes.Reset,
}
