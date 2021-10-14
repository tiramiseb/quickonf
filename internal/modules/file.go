package modules

import (
	"errors"
	"os"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("file", File)
	Register("executable-file", ExecutableFile)
	Register("restricted-file", RestrictedFile)
	Register("read-file", ReadFile)
}

const (
	filePermissionStandard   os.FileMode = 0644
	filePermissionExecutable os.FileMode = 0755
	filePermissionRestricted os.FileMode = 0600
)

// File creates or replaces files
func File(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace file")
	return file(in, out, filePermissionStandard)
}

// ExecutableFile creates or replaces files with executable flag
func ExecutableFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace executable file")
	return file(in, out, filePermissionExecutable)
}

// RestrictedFile creates or replaces files only readable by the owner
func RestrictedFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Create or replace restricted file")
	return file(in, out, filePermissionRestricted)
}

func file(in interface{}, out output.Output, permission os.FileMode) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for path, content := range data {
		path = helper.Path(path)
		if err != nil {
			return err
		}
		out.ShowLoader()
		result, err := helper.File(path, []byte(content), permission, false)
		out.HideLoader()
		switch result {
		case helper.ResultAlready:
			out.Infof("%s already has the needed content", path)
		case helper.ResultDryrun:
			out.Infof("Would create or modify %s", path)
		case helper.ResultError:
			return err
		case helper.ResultSuccess:
			out.Successf("%s created or modified", path)
		}
	}
	return nil
}

// ReadFile reads the content of a file and places it in the store
func ReadFile(in interface{}, out output.Output) error {
	out.InstructionTitle("Read file")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for path, storeKey := range data {
		path = helper.Path(path)
		content, err := os.ReadFile(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		helper.Store(storeKey, string(content))
		out.Successf("Read %s", path)
	}
	return nil
}
