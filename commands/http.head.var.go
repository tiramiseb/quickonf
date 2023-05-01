package commands

import (
	"fmt"
	"net/http"
)

func init() {
	register(httpHeadVar)
}

var httpHeadVar = &Command{
	"http.head.var",
	"Download headers with a HTTP HEAD request (headers are downloaded when checking)",
	[]string{
		"URI of the content for which to download headers",
		"Name of the header",
	},
	[]string{
		"Content of the named header",
	},
	"VS Code\n  headers = http.head.var https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64 Content-Disposition\n  version = regexp.submatch \"filename=\\\"code_([0-9]+\\.[0-9]+\\.[0-9]+)-.*_amd64.deb\\\"\" <headers>\n  ...",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uri := args[0]
		header := args[1]
		resp, err := http.Head(uri)
		if err != nil {
			return []string{""}, fmt.Sprintf("Could not download headers of %s: %s", uri, err), nil, StatusError, "", ""
		}
		resp.Body.Close()
		headerValue := resp.Header.Get(header)
		result = []string{headerValue}
		return result, fmt.Sprintf("Downloaded %s", uri), nil, StatusSuccess, "", headerValue
	},
	nil,
}
