package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func init() {
	register(httpPostVar)
}

var httpPostVar = Command{
	"http.post.var",
	"Send content with a HTTP POST request (data is sent when checking)",
	[]string{
		"URI of the content to download",
		"Body (payload)",
	},
	[]string{
		"Downloaded content",
	},
	"Download example\n  foobar = http.post.var http://www.example.com \"example data\"",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uri := args[0]
		payload := []byte(args[1])
		contentType := mimetype.Detect(payload).String()

		resp, err := http.Post(uri, contentType, bytes.NewReader(payload))
		if err != nil {
			return []string{""}, fmt.Sprintf("Could not send data to %s: %s", uri, err), nil, StatusError, "", ""
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
		return result, fmt.Sprintf("Sent to %s and received response", uri), nil, StatusSuccess, "", varbl
	},
	nil,
}
