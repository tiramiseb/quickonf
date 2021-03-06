package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("vscode-extension", VscodeExtension)
	Register("vscodium-extension", VscodiumExtension)
}

// VscodeExtension installs VS Code extensions
func VscodeExtension(in interface{}, out output.Output) error {
	out.InstructionTitle("Install Vscode extension")
	return vscodeExtension(in, out, "code")
}

// VscodiumExtension installs VS Codium extensions
func VscodiumExtension(in interface{}, out output.Output) error {
	out.InstructionTitle("Install Vscodium extension")
	return vscodeExtension(in, out, "codium")
}

func vscodeExtension(in interface{}, out output.Output, cmd string) error {
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, extension := range data {
		if Dryrun {
			out.Infof("Would install %s", extension)
			continue
		}
		out.Infof("Installing %s", extension)
		out.ShowLoader()
		_, err := helper.Exec(nil, "", cmd, "--install-extension", extension)
		out.HideLoader()
		if err != nil {
			return nil
		}
	}
	return nil
}
