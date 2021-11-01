package helper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// File writes the provided content to the provided path, with the given permission. The file owner is root if the root boolean is true.
//
// If the file already exists and already has the same content, ResultAlready is returned.
//
// In dry-run mode, comparison is done (thus, ResultAlready may be returned) but ResultDryrun is returned if file need to be created or modified.
func File(path string, content []byte, permission os.FileMode, root bool) (ResultStatus, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return ResultError, fmt.Errorf("%s is a directory", path)
		}
		// TODO Read root restricted content...
		current, err := ioutil.ReadFile(path)
		if err != nil {
			return ResultError, err
		}
		if bytes.Compare(content, current) == 0 {
			if !root {
				// TODO also chmod if owned by root
				if !Dryrun {
					err := os.Chmod(path, permission)
					if err != nil {
						return ResultError, err
					}

				}
			}
			return ResultAlready, nil
		}
	} else if !os.IsNotExist(err) {
		return ResultError, err
	}
	if Dryrun {
		return ResultDryrun, nil
	}

	realPath := path
	if root {
		f, err := ioutil.TempFile("", "quickonf-root-file-")
		if err != nil {
			return ResultError, err
		}
		path = f.Name()
		defer os.Remove(path)
		if err := f.Close(); err != nil {
			return ResultError, err
		}
	}

	if err := ioutil.WriteFile(path, content, permission); err != nil {
		return ResultError, err
	}
	// Added chmod because writefile doesn't seem to honor the permissions, even with a correct mask
	if err := os.Chmod(path, permission); err != nil {
		return ResultError, err
	}

	if root {
		if _, err := ExecSudo(nil, "", "cp", "--preserve=mode", path, realPath); err != nil {
			return ResultError, err
		}
	}

	return ResultSuccess, nil
}

// Directory creates a new directory, and its parents if needed.
//
// If the directory already exists, it returns ResultAlready.
//
// In dry-run mode, ResultDryrun is returned.
func Directory(path string) (ResultStatus, error) {
	info, err := os.Lstat(path)
	if err == nil {
		if info.IsDir() {
			return ResultAlready, nil
		}
		return ResultError, fmt.Errorf("%s is not a directory", path)
	}
	if !os.IsNotExist(err) {
		return ResultError, err
	}
	if Dryrun {
		return ResultDryrun, nil
	}
	if err = os.MkdirAll(path, 0755); err != nil {
		return ResultError, err
	}
	return ResultSuccess, nil
}

// Remove removes a file or empty directory from the system.
//
// If the path is a non-empty directory, it returns an error.
//
// If the path does not exist, it returns ResultAlready.
//
// In dry-run mode, file or directory is not removed and ResultDryrun is returned.
func Remove(path string) (ResultStatus, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ResultAlready, nil
		}
		return ResultError, err
	}
	if info.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			return ResultError, err
		}
		defer f.Close()
		_, err = f.Readdirnames(1)
		if err == nil {
			return ResultError, fmt.Errorf("directory %s contains files, cannot be deleted", path)
		}
		if err != io.EOF {
			return ResultError, err
		}
	}

	if Dryrun {
		return ResultDryrun, nil
	}
	if err := os.Remove(path); err != nil {
		return ResultError, err
	}
	return ResultSuccess, nil

}
