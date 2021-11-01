package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"

	"github.com/tiramiseb/quickonf/internal/output"
)

// passThru is inspired by https://stackoverflow.com/a/25645804
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

// DownloadFile downloads the given URL to the given path, writing the percentage to the given output if provided.
//
// If in dry-run mode, it returns ResultDryrun right away, without downloading.
func DownloadFile(url, path string, out output.Output) (ResultStatus, error) {
	if Dryrun {
		return ResultDryrun, nil
	}
	if out != nil {
		out.ShowLoader()
	}
	resp, err := http.Get(url)
	if err != nil {
		if out != nil {
			out.HideLoader()
		}
		return ResultError, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		if out != nil {
			out.HideLoader()
		}
		return ResultError, errors.New("404 not found")
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
		return ResultError, err
	}
	defer f.Close()
	if _, err := io.Copy(f, body); err != nil {
		return ResultError, err
	}
	return ResultSuccess, nil
}

// DownloadJSON downloads the given URL and decodes it to the given destination.
//
// The URL is downloaded even in dry-run mode.
func DownloadJSON(url string, destination interface{}) error {
	body, err := httpReq("GET", url, nil)
	if err != nil {
		return err
	}
	defer body.Close()
	if err := json.NewDecoder(body).Decode(destination); err != nil {
		return err
	}
	return nil
}

// HTTPGet downloads the given URL and returns it as a []byte.
//
// The URL is downloaded even in dry-run mode.
func HTTPGet(url string) ([]byte, error) {
	body, err := httpReq("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(body)
	return buf.Bytes(), err
}

// HTTPPost uses method POST on the given URL and returns the result as a []byte.
//
// The URL is downloaded even in dry-run mode.
func HTTPPost(url string, payload []byte) ([]byte, error) {
	body, err := httpReq("POST", url, payload)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(body)
	return buf.Bytes(), err
}

func httpReq(method, url string, payload []byte) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(payload))
		contentType := mimetype.Detect(payload).String()
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, errors.New("not found")
	}
	return resp.Body, nil
}
