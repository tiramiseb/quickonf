package modules

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("git-config", GitConfig)
	Register("git-clone-pull", GitCloneOrPull)
}

// GitConfig sets got configuration parameters
func GitConfig(in interface{}, out output.Output) error {
	out.InstructionTitle("Git configuration")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for param, value := range data {
		if Dryrun {
			out.Infof("Would set %s to %s", param, value)
			continue
		}
		if _, err := helper.Exec(nil, "git", "config", "--global", param, value); err != nil {
			return err
		}
		out.Successf("Set %s to %s", param, value)
	}
	return nil
}

// GitCloneOrPull clones a Git repository, with depth 1, or pulls latest commit if it already exists
func GitCloneOrPull(in interface{}, out output.Output) error {
	out.InstructionTitle("Git clone or pull")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for repo, dir := range data {
		dir = helper.Path(dir)
		finfo, err := os.Stat(dir)
		if err == nil {
			if !finfo.IsDir() {
				out.Error(fmt.Errorf("%s exists but is not a directory", dir))
				continue
			}
			// Directory exists, verify it is a repository and pull it
			rep, err := git.PlainOpen(dir)
			if err != nil {
				if errors.Is(err, git.ErrRepositoryNotExists) {
					return fmt.Errorf("%s is not a Git repository", dir)
				}
				return err
			}
			if Dryrun {
				out.Infof("Would pull latest commit in %s", dir)
				continue
			}
			out.Infof("%s already exists, pulling latest commit", dir)
			w, err := rep.Worktree()
			if err != nil {
				return err
			}
			if err := w.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
				if errors.Is(err, git.NoErrAlreadyUpToDate) {
					out.Info("... already up-to-date")
					continue
				}
				return err
			}
			continue
		}
		if os.IsNotExist(err) {
			// It does not already exist, cloning the repo
			if err := helper.GitClone(repo, dir, 1, out); err != nil {
				return err
			}
			continue
		}
		return err
	}
	return nil
}
