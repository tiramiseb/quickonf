package helper

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/tiramiseb/quickonf/internal/output"
)

// Passthru is inspired by https://stackoverflow.com/a/25645804
type passThru struct {
	io.ReadCloser
	current int64
	length  int64
	out     output.Output
}

func (pt *passThru) Read(p []byte) (int, error) {
	n, err := pt.ReadCloser.Read(p)
	if n > 0 {
		pt.current += int64(n)
		pt.out.ShowPercentage(int(
			float64(pt.current) / float64(pt.length) * float64(100),
		))
	}
	return n, err
}

// DownloadFileWithPercent downloads the given URL to the given path, writing the percentage to the given output.
func DownloadFileWithPercent(url, path string, out output.Output) error {
	return downloadFile(url, path, out)
}

// DownloadFile downloads the given URL to the given path, without any output.
func DownloadFile(url, path string) error {
	return downloadFile(url, path, nil)
}

func downloadFile(url, path string, out output.Output) error {
	if out != nil {
		out.ShowLoader()
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body := resp.Body
	if out != nil {
		body = &passThru{
			ReadCloser: resp.Body,
			length:     resp.ContentLength,
			out:        out,
		}
		defer out.HidePercentage()
		out.HideLoader()
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, body)
	return err
}

// DownloadJSON downloads the given URL and decodes it to the given destination
func DownloadJSON(url string, destination interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(destination)
}
