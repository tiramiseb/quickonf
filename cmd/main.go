package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	var conf string
	flag.StringVar(&conf, "config", "quickonf.yaml", "List all steps")
	flag.StringVar(&conf, "c", "quickonf.yaml", "List all steps (shorthand)")

	var dryrun bool
	flag.BoolVar(&dryrun, "dry-run", false, "Try all steps without modifying the system")
	flag.BoolVar(&dryrun, "d", false, "Try all steps without modifying the system (shorthand)")

	var help bool
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")

	flag.Parse()

	if help {
		showHelp()
		return
	}

	args := flag.Args()

	if len(args) == 0 {
		apply(conf, nil)
		return
	} else if args[0] == "list" {
		list()
		return
	} else if args[0] == "help" {
		if len(args) == 1 {
			fmt.Println("Please specify the instruction you need help for")
			os.Exit(1)
		}
		instructionHelp(args[1])
		return
	}
	apply(conf, args)
}

func showHelp() {
	fmt.Print(`
Quickonf usage:

Options:

  -config, -c: path to the configuration file (default quickonf.yaml)
  -dry-run, -d: simulate steps without modifying the system
  -help, -h: show this help

  Commands:

  list: list all available instructions
  help <instruction>: get help on the given instruction

Without any argument: apply all steps from the configuration file
With any number of arguments: apply only the given steps (as patterns) from the configuration file
`)
}
