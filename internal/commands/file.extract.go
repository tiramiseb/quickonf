package commands

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
)

var fileExtractors = map[string]func(archive, dest string) error{
	"tar.xz": extractTarxz,
	"zip":    extractZip,
}

func init() {
	register(fileExtract)
}

var fileExtract = Command{
	"file.extract",
	"Extract a (compressed) archive",
	[]string{
		"Archive format (one of zip, tar.xz)",
		"Archive absolute path",
		"Target absolute path",
	},
	nil,
	"Extract something\n  file.extract zip /foo/bar.zip /bar", func(args []string) (result []string, msg string, apply Apply, status Status) {
		format := args[0]
		archive := args[1]
		dest := args[2]
		extractor, ok := fileExtractors[format]
		if !ok {
			return nil, fmt.Sprintf("Archive format %s unknown", format), nil, StatusError
		}
		if !filepath.IsAbs(archive) {
			return nil, fmt.Sprintf("%s is not an absolute path", archive), nil, StatusError
		}
		if !filepath.IsAbs(dest) {
			return nil, fmt.Sprintf("%s is not an absolute path", dest), nil, StatusError
		}

		apply = func(out Output) bool {
			out.Runningf("Extracting %s to %s", archive, dest)
			if err := extractor(archive, dest); err != nil {
				out.Errorf("Could not extract %s to %s: %s", archive, dest, err)
				return false
			}
			out.Successf("Extracted %s to %s", archive, dest)
			return true
		}
		return nil, fmt.Sprintf("Need to extract %s to %s", archive, dest), apply, StatusInfo
	},
	func() {
		datastores.Dconf.Reset()
		datastores.Users.Reset()
	},
}

func extractTarxz(archive, dest string) error {
	fread, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer fread.Close()

	xread, err := xz.NewReader(fread)
	if err != nil {
		return err
	}

	tread := tar.NewReader(xread)
	for {
		header, err := tread.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(target, header.FileInfo().Mode()); err != nil {
						return err
					}
					continue
				}
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tread)
			f.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func extractZip(archive, dest string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0o755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
			return err
		}

		outfile, err := os.Create(fpath)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outfile.Close()
			return err
		}

		_, err = io.Copy(outfile, rc)
		outfile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
