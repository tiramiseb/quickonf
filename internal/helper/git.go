package helper

import (
	"github.com/go-git/go-git/v5"

	"github.com/tiramiseb/quickonf/internal/output"
)

// GitClone clones a Git repository
func GitClone(repo, dest string, depth int, out output.Output) error {
	if Dryrun {
		if out != nil {
			out.Infof("Would clone %s into %s", repo, dest)
		}
		return nil
	}
	if out != nil {
		out.Infof("Cloning %s into %s", repo, dest)
	}
	options := git.CloneOptions{
		URL: repo,
	}
	if depth > 0 {
		options.Depth = depth
	}

	if out != nil {
		out.ShowLoader()
	}
	_, err := git.PlainClone(dest, false, &options)
	if out != nil {
		out.HideLoader()
	}
	return err
}
