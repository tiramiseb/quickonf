package helper

import (
	"fmt"

	"os"
	"os/user"
)

var Home string

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Home = u.HomeDir
}
