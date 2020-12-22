package helper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/output"
)

// UnzipFile unzips the given file into the given destination
func UnzipFile(zipfilepath, dest string, out output.Output) error {
	if Dryrun {
		out.Info("Would extract " + zipfilepath + " to " + dest)
		return nil
	}
	out.Info("Extracting " + zipfilepath + " to " + dest)
	if out != nil {
		defer out.HidePercentage()
	}

	r, err := zip.OpenReader(zipfilepath)
	if err != nil {
		return err
	}
	defer r.Close()
	totalfiles := len(r.File)

	for i, f := range r.File {
		if out != nil {
			out.ShowXonY(i, totalfiles)
		}

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
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
