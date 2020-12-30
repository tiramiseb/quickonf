package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("download", Download)
}

// Download downloads files from given URLs
func Download(in interface{}, out output.Output) error {
	out.InstructionTitle("Download")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for url, path := range data {
		path = helper.Path(path)
		result, err := helper.DownloadFile(url, path, out)
		switch result {
		case helper.ResultDryrun:
			out.Infof("Would download %s to %s", url, path)
		case helper.ResultError:
			return err
		case helper.ResultSuccess:
			out.Successf("Downloaded %s to %s", url, path)
		}
	}
	return nil
}
