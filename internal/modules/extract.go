package modules

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("extract-tarxz", ExtractTarxz)
	Register("extract-zip", ExtractZip)
}

// ExtractTarxz extracts a tar.xz archive
func ExtractTarxz(in interface{}, out output.Output) error {
	out.InstructionTitle("Extract tar.xz")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for tarxzFile, dest := range data {
		tarxzFile = helper.Path(tarxzFile)
		dest = helper.Path(dest)
		if Dryrun {
			out.Info("Would extract " + tarxzFile + " to " + dest)
			continue
		}
		out.Info("Extracting " + tarxzFile + " to " + dest)
		fread, err := os.Open(tarxzFile)
		if err != nil {
			return err
		}
		defer fread.Close()

		xread, err := xz.NewReader(fread)
		if err != nil {
			return err
		}

		tread := tar.NewReader(xread)
		out.ShowLoader()
		for {
			header, err := tread.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				out.HideLoader()
				return err
			}

			target := filepath.Join(dest, header.Name)

			switch header.Typeflag {
			case tar.TypeDir:
				if _, err := os.Stat(target); err != nil {
					if os.IsNotExist(err) {
						if err := os.MkdirAll(target, header.FileInfo().Mode()); err != nil {
							out.HideLoader()
							return err
						}
						continue
					}
					out.HideLoader()
					return err
				}
			case tar.TypeReg:
				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, header.FileInfo().Mode())
				if err != nil {
					out.HideLoader()
					return err
				}
				_, err = io.Copy(f, tread)
				f.Close()
				if err != nil {
					out.HideLoader()
					return err
				}
			}
		}
		out.HideLoader()
	}
	return nil
}

// ExtractZip extracts a zip archive
func ExtractZip(in interface{}, out output.Output) error {
	out.InstructionTitle("Extract Zip")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for zipFile, dest := range data {
		zipFile = helper.Path(zipFile)
		dest = helper.Path(dest)
		if err := helper.UnzipFile(zipFile, dest, out); err != nil {
			return err
		}
	}
	return nil
}
