package modules

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("file", File)
	Register("executable-file", ExecutableFile)
	Register("restricted-file", RestrictedFile)
	Register("root-file", RootFile)
	Register("executable-root-file", ExecutableRootFile)
	Register("restricted-root-file", ExecutableRootFile)
}

const (
	filePermissionStandard   os.FileMode = 0644
	filePermissionExecutable os.FileMode = 0755
	filePermissionRestricted os.FileMode = 0600
)

func File(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Create or replace file")
	return file(in, out, false, filePermissionStandard, store)
}

func ExecutableFile(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Create or replace executable file")
	return file(in, out, false, filePermissionExecutable, store)
}

func RestrictedFile(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Create or replace restricted file")
	return file(in, out, false, filePermissionRestricted, store)
}

func RootFile(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Create or replace file as root")
	return file(in, out, true, filePermissionStandard, store)
}

func ExecutableRootFile(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Create or replace executable file as root")
	return file(in, out, true, filePermissionExecutable, store)
}

func RestrictedRootFile(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Create or replace restricted file as root")
	return file(in, out, true, filePermissionRestricted, store)
}

func file(in interface{}, out output.Output, root bool, permission os.FileMode, store map[string]interface{}) error {
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for path, content := range data {
		path = helper.Path(path)
		result, err := helper.Template(content)
		if err != nil {
			return err
		}
		info, err := os.Stat(path)
		if err == nil {
			if info.IsDir() {
				return errors.New(path + " is a directory")
			}
			current, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if bytes.Compare(result, current) == 0 {
				out.Info(path + " already has the needed content")
                if !root {
					err := os.Chmod(path, permission)
					if err != nil {
						return err
					}
                }
				return nil
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		if root {
			f, err := ioutil.TempFile("", "quickonf-root-file")
			if err != nil {
				return err
			}
			fName := f.Name()
			defer os.Remove(fName)
			defer f.Close()
			_, err = f.Write(result)
			if err != nil {
				return err
			}
			err = os.Chmod(fName, permission)
			if err != nil {
				return err
			}
			err = helper.ExecSudo("cp", fName, path)
			if err != nil {
				return err
			}
		} else {
			if err := ioutil.WriteFile(path, result, permission); err != nil {
				return err
			}
		}
		out.Success(path + " created or modified")
	}
	return nil
}
