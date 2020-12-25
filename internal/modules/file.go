package modules

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("file", File)
	Register("executable-file", ExecutableFile)
	Register("restricted-file", RestrictedFile)
	Register("root-file", RootFile)
	Register("executable-root-file", ExecutableRootFile)
	Register("restricted-root-file", RestrictedRootFile)
}

const (
	filePermissionStandard   os.FileMode = 0644
	filePermissionExecutable os.FileMode = 0755
	filePermissionRestricted os.FileMode = 0600
)

// File creates or replaces files
func File(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace file")
	return file(in, out, false, filePermissionStandard)
}

// ExecutableFile creates or replaces files with executable flag
func ExecutableFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace executable file")
	return file(in, out, false, filePermissionExecutable)
}

// RestrictedFile creates or replaces files only readable by the owner
func RestrictedFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace restricted file")
	return file(in, out, false, filePermissionRestricted)
}

// RootFile creates or replaces files as root
func RootFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace file as root")
	return file(in, out, true, filePermissionStandard)
}

// ExecutableRootFile creates or replaces executable files as root
func ExecutableRootFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace executable file as root")
	return file(in, out, true, filePermissionExecutable)
}

// RestrictedRootFile creates or replaces files only readable by root
func RestrictedRootFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace restricted file as root")
	return file(in, out, true, filePermissionRestricted)
}

func file(in interface{}, out output.Output, root bool, permission os.FileMode) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for path, content := range data {
		path = helper.Path(path)
		if err != nil {
			return err
		}
		bcontent := []byte(content)
		info, err := os.Stat(path)
		if err == nil {
			if info.IsDir() {
				return fmt.Errorf("%s is a directory", path)
			}
			// TODO Read root restricted content...
			current, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if bytes.Compare(bcontent, current) == 0 {
				out.Infof("%s already has the needed content", path)
				if !root {
					if !Dryrun {
						err := os.Chmod(path, permission)
						if err != nil {
							return err
						}
					}
				}
				continue
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		if Dryrun {
			out.Infof("Would create or modify %s", path)
			continue
		}
		if root {
			f, err := ioutil.TempFile("", "quickonf-root-file")
			if err != nil {
				return err
			}
			fName := f.Name()
			defer os.Remove(fName)
			defer f.Close()
			if _, err = f.Write(bcontent); err != nil {
				return err
			}
			if err := os.Chmod(fName, permission); err != nil {
				return err
			}
			if _, err := helper.ExecSudo(nil, "cp", fName, path); err != nil {
				return err
			}
		} else {
			if err := ioutil.WriteFile(path, bcontent, permission); err != nil {
				return err
			}
		}
		out.Successf("%s created or modified", path)
	}
	return nil
}
