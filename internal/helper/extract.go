package helper

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/output"
	"github.com/ulikunitz/xz"
)

// ExtractTarxz extracts a tar.xz archive.
//
// In dry-run mode, only return ResultDryrun.
func ExtractTarxz(tarxzfilepath, dest string) (ResultStatus, error) {
	if Dryrun {
		return ResultDryrun, nil
	}
	fread, err := os.Open(tarxzfilepath)
	if err != nil {
		return ResultError, err
	}
	defer fread.Close()

	xread, err := xz.NewReader(fread)
	if err != nil {
		return ResultError, err
	}

	tread := tar.NewReader(xread)
	for {
		header, err := tread.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return ResultError, err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(target, header.FileInfo().Mode()); err != nil {
						return ResultError, err
					}
					continue
				}
				return ResultError, err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, header.FileInfo().Mode())
			if err != nil {
				return ResultError, err
			}
			_, err = io.Copy(f, tread)
			f.Close()
			if err != nil {
				return ResultError, err
			}
		}
	}
	return ResultSuccess, nil
}

// ExtractZip unzips the given file into the given destination.
//
// In dry-run mode, only return ResultDryrun.
//
// If out is provided, a "X on Y" loader is provided to the user.
func ExtractZip(zipfilepath, dest string, out output.Output) (ResultStatus, error) {
	if Dryrun {
		return ResultDryrun, nil
	}
	if out != nil {
		defer out.HideXonY()
	}

	r, err := zip.OpenReader(zipfilepath)
	if err != nil {
		return ResultError, err
	}
	defer r.Close()
	totalfiles := len(r.File)

	for i, f := range r.File {
		if out != nil {
			out.ShowXonY(i, totalfiles)
		}

		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return ResultError, err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return ResultError, err
		}

		outfile, err := os.Create(fpath)
		if err != nil {
			return ResultError, err
		}

		rc, err := f.Open()
		if err != nil {
			outfile.Close()
			return ResultError, err
		}

		_, err = io.Copy(outfile, rc)
		outfile.Close()
		rc.Close()
		if err != nil {
			return ResultError, err
		}
	}
	return ResultSuccess, nil
}
