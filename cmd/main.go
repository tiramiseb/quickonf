package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/tiramiseb/quickonf/conf"
	"github.com/tiramiseb/quickonf/instructions"
	"github.com/tiramiseb/quickonf/program"
)

func main() {
	rand.Seed(time.Now().Unix())

	config := "quickonf.qconf"
	var configFromFlag string
	flag.StringVar(&configFromFlag, "config", "", "path to the configuration file")
	flag.StringVar(&configFromFlag, "c", "", "path to the configuration file (shorthand)")
	checkStdin := flag.Bool("check-stdin", false, "check configuration on stdin")
	flag.Parse()

	if *checkStdin {
		instructions.NewGlobalVar("confdir", "-")
		err := conf.Check(os.Stdin)
		if err != nil {
			fmt.Printf(`{"error": "Could not check input data: %s"}\n`, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if configFromFlag == "" {
		args := flag.Args()
		if len(args) > 0 {
			config = args[0]
		}
	} else {
		config = configFromFlag
	}

	r, err := os.Open(config)
	if err != nil {
		fmt.Println("Could not open configuration file", config)
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()
	instructions.NewGlobalVar("confdir", filepath.Dir(config))
	groups, errs := conf.Read(r)
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
	program.Run(groups)
}
