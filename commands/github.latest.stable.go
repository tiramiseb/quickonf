package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(githubLatestStable)
}

var githubLatestStable = Command{
	"github.latest.stable",
	"Check the latest stable release from a GitHub repository",
	[]string{
		"GitHub repository name",
		"Pattern to match an asset",
	},
	[]string{
		"Release name",
		"Asset URL",
	},
	"Kmonad\n  release url = github.latest.stable kmonad/kmonad kmonad-*-linux",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		repository := args[0]
		pattern := args[1]
		data := gitHubRelease{}
		if err := helper.DownloadJSON("https://api.github.com/repos/"+repository+"/releases/latest", &data); err != nil {
			return nil, fmt.Sprintf("Could not get information for %s: %s", repository, err), nil, StatusError, "", ""
		}
		name, url, err := data.Extract(pattern)
		if err != nil {
			return nil, fmt.Sprintf("Could not extract information for %s: %s", repository, err), nil, StatusError, "", ""
		}
		return []string{name, url}, fmt.Sprintf("Got information for %s", repository), nil, StatusSuccess, "", fmt.Sprintf("Name: %s, URL: %s", name, url)
	},
	nil,
}
