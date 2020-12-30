package helper

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

// GitClone clones a Git repository to the provided destination.
//
// depth is the number of commits to fetch for clone (if 0, only fetch HEAD).
//
// If the provided destination is already this repository, return ResultAlready
//
// In dry-run mode, do not clone, return ResultDryrun
func GitClone(repo, dest string, depth int) (ResultStatus, error) {
	finfo, err := os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			// destination does not exist, clone the repository
			if Dryrun {
				return ResultDryrun, nil
			}
			options := git.CloneOptions{
				URL: repo,
			}
			if depth > 0 {
				options.Depth = depth
			}
			if _, err := git.PlainClone(dest, false, &options); err != nil {
				return ResultError, err
			}
			return ResultSuccess, nil
		}
		return ResultError, err
	}
	// File exists, check if it is the requested repository
	if !finfo.IsDir() {
		return ResultError, fmt.Errorf("%s exists but is not a directory", dest)
	}
	rep, err := git.PlainOpen(dest)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return ResultError, fmt.Errorf("%s is not a Git repository", dest)
		}
		return ResultError, err
	}
	remotes, err := rep.Remotes()
	if err != nil {
		return ResultError, err
	}
	for _, remote := range remotes {
		urls := remote.Config().URLs
		if len(urls) == 0 {
			return ResultError, fmt.Errorf("No remote for repository in %s", dest)
		}
		if urls[0] == repo {
			return ResultAlready, nil
		}
	}
	return ResultError, fmt.Errorf("Local repository does not have %s ad an origin", repo)
}

// GitPull pulls the latest HEAD of the given Git repository from origin.
//
// If the provided destination is already at the latest HEAD from origin, return ResultAlready.
//
// In dry-run mode, do not pull, return ResultDryrun.
func GitPull(path string) (ResultStatus, error) {

	finfo, err := os.Stat(path)
	if err != nil {
		return ResultError, err
	}

	if !finfo.IsDir() {
		return ResultError, fmt.Errorf("%s exists but is not a directory", path)
	}
	// Directory exists, verify it is a repository and pull it
	rep, err := git.PlainOpen(path)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return ResultError, fmt.Errorf("%s is not a Git repository", path)
		}
		return ResultError, err
	}
	if Dryrun {
		return ResultDryrun, nil
	}
	w, err := rep.Worktree()
	if err != nil {
		return ResultError, err
	}
	if err := w.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return ResultAlready, nil
		}
		return ResultError, err
	}
	return ResultSuccess, nil
}
