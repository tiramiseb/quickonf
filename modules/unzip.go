package modules

import (
	"archive/zip"
	"io"
	"os"
	"path"

	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("unzip", Unzip)
}

func Unzip(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Unzip")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for zipFile, destPath := range data {
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
		out.Info("Uncompressed " + zipFile)
	}
	return nil
}
