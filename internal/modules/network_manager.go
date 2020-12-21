package modules

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("nm-import-openvpn", NetworkManagerImport)
}

// NetworkManagerImport imports a configuration into network manager
func NetworkManagerImport(in interface{}, out output.Output) error {
	out.InstructionTitle("Importing configuration into network manager")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
RANGE:
	for _, path := range data {
		path = helper.Path(path)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				out.Error(err)
				continue
			}
			return err
		}
		nameParts := strings.Split(filepath.Base(path), ".")
		var nameS string
		switch len(nameParts) {
		case 0:
			return errors.New("No file name in " + path)
		case 1, 2:
			nameS = nameParts[0]
		default:
			nameS = strings.Join(nameParts[0:len(nameParts)-2], ".")
		}
		path = helper.Path(path)
		nameB := []byte(nameS)

		shown, err := helper.Exec("nmcli", "connection", "show")
		if err != nil {
			return err
		}

		for _, line := range bytes.Split(shown, []byte{'\n'}) {
			fields := bytes.Fields(line)
			var connName []byte
			if len(fields) > 0 {
				connName = fields[0]
			}
			if bytes.Equal(connName, nameB) {
				out.Info(nameS + " is already configured")
				continue RANGE
			}
		}

		if Dryrun {
			out.Info("Would import " + path)
			continue
		}

		if _, err := helper.Exec("nmcli", "connection", "import", "type", "openvpn", "file", path); err != nil {
			return err
		}
	}
	return nil
}
