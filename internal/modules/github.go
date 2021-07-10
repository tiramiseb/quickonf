package modules

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("github-latest", GithubLatest)
}

// GithubLatest checks latest release of a GitHub repository
func GithubLatest(in interface{}, out output.Output) error {
	out.InstructionTitle("Check latest release from GitHub")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	repository, ok := data["repository"]
	if !ok {
		return errors.New("missing repository owner/name")
	}

	result := struct {
		Name    string `json:"name"`
		TagName string `json:"tag_name"`
		Assets  []struct {
			URL  string `json:"browser_download_url"`
			Name string `json:"name"`
		} `json:"assets"`
	}{}
	// Not using the GraphQL API because it would need to be authenticated... too bad, for public data
	out.ShowLoader()
	err = helper.DownloadJSON("https://api.github.com/repos/"+repository+"/releases/latest", &result)
	out.HideLoader()
	if err != nil {
		return err
	}
	out.Infof("Latest release for %s is %s", repository, result.TagName)
	if storeRelease, ok := data["store"]; ok {
		helper.Store(storeRelease, result.TagName)
	}
	pattern, ok := data["pattern"]
	if !ok {
		// No pattern requested, stop here
		return nil
	}
	matching := []int{}
	for i, asset := range result.Assets {
		ok, err := path.Match(pattern, asset.Name)
		if err != nil {
			return err
		}
		if ok {
			matching = append(matching, i)
		}
	}
	if len(matching) == 0 {
		return fmt.Errorf(`no asset matching "%s" in %s`, pattern, repository)
	}
	if len(matching) > 1 {
		names := make([]string, len(matching))
		for i, m := range matching {
			names[i] = result.Assets[m].Name
		}
		return fmt.Errorf("too many assets matching pattern in %s: %s", repository, strings.Join(names, ", "))
	}
	url := result.Assets[matching[0]].URL
	out.Infof("Download URL for latest release is %s", url)
	if storeURL, ok := data["store-url"]; ok {
		helper.Store(storeURL, url)
	}
	return nil
}
