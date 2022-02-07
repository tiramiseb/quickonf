package commands

import (
	"fmt"
	"path"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
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

func (g gitHubRelease) Extract(pattern string) (tagname, url string, err error) {
	tagname = g.TagName
	if pattern == "" {
		return
	}
	matching := []int{}
	for i, asset := range g.Assets {
		var ok bool
		ok, err = path.Match(pattern, asset.Name)
		if err != nil {
			return
		}
		if ok {
			matching = append(matching, i)
		}
	}
	if len(matching) == 0 {
		err = fmt.Errorf(`no asset matching "%s"`, pattern)
		return
	}
	url = g.Assets[matching[0]].URL
	if len(matching) > 1 {
		names := make([]string, len(matching))
		for i, m := range matching {
			names[i] = g.Assets[m].Name
		}
		err = fmt.Errorf("too many assets matching pattern: %s", strings.Join(names, ", "))
	}
	return
}

func init() {
	register(githubLatest)
}

var githubLatest = Command{
	"github.latest",
	"Check the latest release from a GitHub repository",
	[]string{
		"GitHub repository name",
		"Pattern to match an asset",
	},
	[]string{
		"Release name",
		"Asset URL",
	},
	"Kmonad\n  release url = github.latest kmonad/kmonad kmonad-*-linux",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		repository := args[0]
		pattern := args[1]
		data := make([]gitHubRelease, 1)
		if err := helper.DownloadJSON("https://api.github.com/repos/"+repository+"/releases?per_page=1", &data); err != nil {
			return nil, fmt.Sprintf("Could not get information for %s: %s", repository, err), nil, StatusError
		}
		name, url, err := data[0].Extract(pattern)
		if err != nil {
			return nil, fmt.Sprintf("Could not extract information for %s: %s", repository, err), nil, StatusError
		}
		return []string{name, url}, fmt.Sprintf("Got information for %s", repository), nil, StatusSuccess
	},
	nil,
}
