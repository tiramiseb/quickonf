package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/conf"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program"
)

func main() {
	config := "quickonf.qconf"
	if len(os.Args) > 1 {
		config = os.Args[1]
	}

	r, err := os.Open(config)
	if err != nil {
		fmt.Println("Could not open configuration file", config)
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()
	instructions.NewGlobalVar("confdir", filepath.Dir(config))
	state, errs := conf.Read(r)
	if errs != nil {
		fmt.Println("Configuration file", config, "is invalid:")
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Could not get current user")
		fmt.Println(err)
		os.Exit(1)
	}
	if usr.Name != "root" {
		fmt.Println("Quickonf must run as root")
		os.Exit(1)
	}
	program.Run(state)
}
