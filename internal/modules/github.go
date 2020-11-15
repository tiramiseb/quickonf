package modules

import (
	"errors"
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
		return errors.New("Missing repository owner/name")
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
	out.Info("Latest release for " + repository + " is " + result.TagName)
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
		return errors.New("No asset matching \"" + pattern + "\" in " + repository)
	}
	if len(matching) > 1 {
		names := make([]string, len(matching))
		for i, m := range matching {
			names[i] = result.Assets[m].Name
		}
		return errors.New("Too many assets matching pattern in " + repository + ": " + strings.Join(names, ", "))
	}
	url := result.Assets[matching[0]].URL
	out.Info("Download URL for latest release is " + url)
	if storeURL, ok := data["store-url"]; ok {
		helper.Store(storeURL, url)
	}
	return nil
}
