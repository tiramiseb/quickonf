package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"

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
	if Dryrun {
		out.Infof("Would download %s to %s", url, path)
		return nil
	}
	out.Infof("Downloading %s to %s", url, path)
	if out != nil {
		out.ShowLoader()
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return errors.New("404 not found")
	}
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
	if resp.StatusCode == http.StatusNotFound {
		return errors.New("404 not found")
	}
	return json.NewDecoder(resp.Body).Decode(destination)
}

// Download downloads the given URL and returns it as a []byte
func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("404 not found")
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}

// Post uses method POST in the given URL
func Post(url string, payload []byte) ([]byte, error) {
	contentType := mimetype.Detect(payload).String()
	body := bytes.NewReader(payload)
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("404 not found")
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}
