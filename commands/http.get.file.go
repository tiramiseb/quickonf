package commands

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(httpGetFile)
}

var httpGetFile = &Command{
	"http.get.file",
	"Download a file with a HTTP GET request (URI is downloaded when applying)",
	[]string{
		"URI of the file to download",
		"Path of the destination file",
	},
	nil,
	"Download example\n  http.get.file http://www.example.com /opt/example",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uri := args[0]
		path := args[1]
		_, err := os.Stat(path)
		if err == nil {
			return nil, fmt.Sprintf("%s already exists", path), nil, StatusSuccess, "", ""
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Sprintf("Could not check if %s exists: %s", path, err), nil, StatusError, "", ""
		}
		apply = func(out Output) bool {
			percentChan := out.RunningPercentf("Downloading %s to %s", uri, path)
			f, err := os.Create(path)
			if err != nil {
				out.Errorf("Could not create %s: %s", path, err)
				return false
			}
			defer f.Close()
			resp, err := http.Get(uri)
			if err != nil {
				out.Errorf("Could not download %s: %s", uri, err)
				return false
			}
			defer resp.Body.Close()
			body := helper.NewProgressReader(resp.Body, resp.ContentLength, percentChan)
			if resp.StatusCode == http.StatusNotFound {
				out.Errorf("%s is not found", uri)
				return false
			}
			if _, err := io.Copy(f, body); err != nil {
				out.Errorf("Could not write content to %s: %s", path, err)
				return false
			}
			out.Successf("Downloaded %s to %s", uri, path)
			return true
		}
		return nil, fmt.Sprintf("Need to download %s to %s", uri, path), apply, StatusInfo, "", ""
	},
	nil,
}
