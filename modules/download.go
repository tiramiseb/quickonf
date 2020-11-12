package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("download", Download)
}

func Download(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Download")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for url, path := range data {
		out.Info("Downloading " + url + " to " + path)
		if err := helper.DownloadFileWithPercent(url, path, out); err != nil {
			return err
		}
	}
	return nil
}
