package instructions

import (
	"fmt"
	"os"
)

type FilePresent struct {
	Path string
}

func (f *FilePresent) Compare(vars *Variables) bool {
	path := vars.TranslateVariables(f.Path)
	_, err := os.Stat(path)
	return err == nil
}

func (f *FilePresent) String() string {
	return fmt.Sprintf("file.present %s", f.Path)
}
