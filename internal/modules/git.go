package modules

import (
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
		out.ShowLoader()
		result, err := helper.GitClone(repo, dir, 0)
		out.HideLoader()
		switch result {
		case helper.ResultDryrun:
			out.Infof("Would clone %s in %s", repo, dir)
			continue
		case helper.ResultSuccess:
			out.Successf("Cloned %s in %s", repo, dir)
			continue
		case helper.ResultError:
			return err
		}
		// The only possible situation here is the repository already existing: we can pull
		out.ShowLoader()
		result, err = helper.GitPull(dir)
		out.HideLoader()
		switch result {
		case helper.ResultDryrun:
			out.Infof("Would pull in %s", dir)
		case helper.ResultAlready:
			out.Infof("%s already in sync with origin", dir)
		case helper.ResultError:
			return err
		case helper.ResultSuccess:
			out.Successf("Pulled latest modifications in %s", dir)
		}
	}
	return nil
}
