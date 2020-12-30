package modules

import (
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
		out.Infof("Extracting %s to %s", tarxzFile, dest)
		out.ShowLoader()
		result, err := helper.ExtractTarxz(tarxzFile, dest)
		out.HideLoader()
		switch result {
		case helper.ResultDryrun:
			out.Info("... not really extracted")
		case helper.ResultSuccess:
			out.Success("... extracted successfully")
		case helper.ResultError:
			return err
		}
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
		result, err := helper.ExtractZip(zipFile, dest, out)
		switch result {
		case helper.ResultDryrun:
			out.Infof("Would extract %s into %s", zipFile, dest)
		case helper.ResultError:
			return err
		case helper.ResultSuccess:
			out.Successf("Extracted %s into %s", zipFile, dest)
		}
	}
	return nil
}
