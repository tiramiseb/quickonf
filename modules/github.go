package modules

import (
	"errors"
	"path"
	"strings"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("github-latest", GithubLatest)
}

func GithubLatest(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Check latest release from GitHub")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	repository, ok := data["repository"]
	if !ok {
		return errors.New("Missing repository owner/name")
	}

	result := struct {
		Name   string `json:"name"`
		Assets []struct {
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
	out.Info("Latest release for " + repository + " is " + result.Name)
	if storeRelease, ok := data["store-release"]; ok {
		store[storeRelease] = result.Name
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
		store[storeURL] = url
	}
	return nil
}
