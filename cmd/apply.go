package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/tiramiseb/quickonf/conf"
	"github.com/tiramiseb/quickonf/program"
	"github.com/tiramiseb/quickonf/state"
)

func apply(config string, filters []string, options state.Options) {
	r, err := os.Open(config)
	if err != nil {
		fmt.Println("Could not open configuration file", config)
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()
	state, errs := conf.Read(r, filters)
	if errs != nil {
		fmt.Println("Configuration file", config, "is invalid:")
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	state.Options = options
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
