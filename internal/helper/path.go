package helper

import (
	"fmt"

	"os"
	"os/user"
)

var homepath string

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	homepath = u.HomeDir
	Store("home", homepath)
}

// Path returns the absolute path for the given file path.
// Relative paths are prefixed with the home directory.
// if the argument is empty, it returns the home directory path.
func Path(str string) string {
	if len(str) == 0 {
		return homepath
	}
	if str[0] == '/' {
		return str
	}
	return homepath + "/" + str
}
