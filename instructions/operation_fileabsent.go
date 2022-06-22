package instructions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type FileAbsent struct {
	Path string
}

func (f *FileAbsent) Compare(vars Variables) bool {
	path := vars.TranslateVariables(f.Path)
	_, err := os.Stat(path)
	return err != nil && errors.Is(err, fs.ErrNotExist)
}

func (f *FileAbsent) String() string {
	return fmt.Sprintf("file.absent %s", f.Path)
}
