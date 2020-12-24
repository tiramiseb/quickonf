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
		if err := helper.DownloadFileWithPercent(url, path, out); err != nil {
			return err
		}
	}
	return nil
}
