package commands

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func init() {
	register(httpGetVar)
}

var httpGetVar = &Command{
	"http.get.var",
	"Download content with a HTTP GET request (URI is downloaded when checking)",
	[]string{
		"URI of the content to download",
	},
	[]string{
		"Downloaded content",
	},
	"Download example\n  foobar = http.get.var http://www.example.com",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uri := args[0]
		resp, err := http.Get(uri)
		if err != nil {
			return []string{""}, fmt.Sprintf("Could not download %s: %s", uri, err), nil, StatusError, "", ""
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return []string{""}, fmt.Sprintf("%s not found", uri), nil, StatusError, "", ""
		}
		var buf strings.Builder
		if _, err := io.Copy(&buf, resp.Body); err != nil {
			return []string{""}, fmt.Sprintf("Could not read data: %s", err), nil, StatusError, "", ""
		}
		varbl := buf.String()
		result = []string{varbl}
		return result, fmt.Sprintf("Downloaded %s", uri), nil, StatusSuccess, "", varbl
	},
	nil,
}
