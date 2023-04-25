package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/user"
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

	if configFromFlag == "" {
		args := flag.Args()
		if len(args) > 0 {
			config = args[0]
		}
	} else {
		config = configFromFlag
	}
	program.Run(config)
}
