package modules

import (
	"archive/zip"
	"io"
	"os"
	"path"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("unzip", Unzip)
}

// Unzip extracts a zip archive
func Unzip(in interface{}, out output.Output) error {
	out.InstructionTitle("Unzip")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for zipFile, destPath := range data {
		destPath = helper.Path(destPath)
		r, err := zip.OpenReader(zipFile)
		if err != nil {
			return err
		}
		for _, f := range r.File {
			content, err := f.Open()
			if err != nil {
				return err
			}
			defer content.Close()
			fileDest := path.Join(destPath, f.FileHeader.Name)
			if f.FileInfo().IsDir() {
				if err := os.MkdirAll(fileDest, 0755); err != nil {
					return err
				}
				continue
			} else {
				dir := path.Dir(fileDest)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
				dest, err := os.Create(fileDest)
				if err != nil {
					return err
				}
				defer dest.Close()
				if _, err = io.Copy(dest, content); err != nil {
					return err
				}
			}
		}
		out.Info("Uncompressed " + zipFile)
	}
	return nil
}
