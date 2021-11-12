package modules

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

// Not using the GraphQL API because it would need to be authenticated... too bad, for public data
type gitHubRelease struct {
	Name       string `json:"name"`
	TagName    string `json:"tag_name"`
	Draft      bool   `json:"draft"`
	PreRelease bool   `json:"prerelease"`
	Assets     []struct {
		URL  string `json:"browser_download_url"`
		Name string `json:"name"`
	} `json:"assets"`
}

func init() {
	Register("github-latest", GithubLatest)
	Register("github-latest-stable", GithubLatestStable)
}

// GithubLatest checks latest release (including drafts and pre-releases) of a GitHub repository
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

	result := make([]gitHubRelease, 1)
	out.ShowLoader()
	err = helper.DownloadJSON("https://api.github.com/repos/"+repository+"/releases?per_page=1", &result)
	out.HideLoader()
	if err != nil {
		return err
	}
	var suffix string
	switch {
	case result[0].Draft && result[0].PreRelease:
		suffix = " (draft for pre-release)"
	case result[0].Draft:
		suffix = " (draft)"
	case result[0].PreRelease:
		suffix = " (pre-release)"
	default:
		suffix = " (stable)"
	}
	// out.Infof("Latest release for %s is %s (%s)", repository, result[0].TagName, suffix)
	return extractGithubRelease(result[0], repository, data["store"], data["pattern"], data["store-url"], suffix, out)

}

// GithubLatestStable checks latest stable release of a GitHub repository
func GithubLatestStable(in interface{}, out output.Output) error {
	out.InstructionTitle("Check latest stable release from GitHub")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	repository, ok := data["repository"]
	if !ok {
		return errors.New("missing repository owner/name")
	}

	result := gitHubRelease{}
	out.ShowLoader()
	err = helper.DownloadJSON("https://api.github.com/repos/"+repository+"/releases/latest", &result)
	out.HideLoader()
	if err != nil {
		return err
	}
	return extractGithubRelease(result, repository, data["store"], data["pattern"], data["store-url"], "", out)
}

func extractGithubRelease(release gitHubRelease, repository, store, pattern, urlStore string, suffix string, out output.Output) error {
	out.Infof("Latest release for %s is %s", repository, release.TagName)
	if store != "" {
		helper.Store(store, release.TagName)
	}
	if pattern == "" {
		return nil
	}
	matching := []int{}
	for i, asset := range release.Assets {
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
			names[i] = release.Assets[m].Name
		}
		return fmt.Errorf("too many assets matching pattern in %s: %s", repository, strings.Join(names, ", "))
	}
	url := release.Assets[matching[0]].URL
	out.Infof("Download URL for latest release is %s", url)
	if urlStore != "" {
		helper.Store(urlStore, url)
	}
	return nil

}
