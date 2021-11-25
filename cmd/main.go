package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tiramiseb/quickonf/state"
)

func main() {

	var conf string
	flag.StringVar(&conf, "config", "quickonf.qconf", "List all steps")
	flag.StringVar(&conf, "c", "quickonf.qconf", "List all steps (shorthand)")

	var dryrun bool
	flag.BoolVar(&dryrun, "dry-run", false, "Try all steps without modifying the system")
	flag.BoolVar(&dryrun, "d", false, "Try all steps without modifying the system (shorthand)")

	var slow bool
	flag.BoolVar(&slow, "slow", false, "Run the steps slowly (wait 500ms between instructions)")
	flag.BoolVar(&slow, "s", false, "Run the steps slowly (shorthand)")

	var help bool
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")

	flag.Parse()

	if help {
		showHelp()
		return
	}

	args := flag.Args()
	options := state.Options{
		DryRun: dryrun,
		Slow:   slow,
	}
	if len(args) == 0 {
		apply(conf, nil, options)
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
	apply(conf, args, options)
}

func showHelp() {
	fmt.Print(`
Quickonf usage:

Options:

  -config, -c: path to the configuration file (default quickonf.qconf)
  -dry-run, -d: simulate steps without modifying the system
  -slow, -s: run steps slowly (1s between two instructions - useful for debug)
  -help, -h: show this help

  Commands:

  list: list all available instructions
  help <instruction>: get help on the given instruction

Without any argument: apply all steps from the configuration file
With any number of arguments: apply only the given steps (as patterns) from the configuration file
`)
}
