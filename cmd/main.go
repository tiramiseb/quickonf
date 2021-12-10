package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tiramiseb/quickonf/state"
)

func main() {

	var conf string
	flag.StringVar(&conf, "config", "quickonf.qconf", "")
	flag.StringVar(&conf, "c", "quickonf.qconf", "")

	var dryrun bool
	flag.BoolVar(&dryrun, "dry-run", false, "")
	flag.BoolVar(&dryrun, "d", false, "")

	var slow bool
	flag.BoolVar(&slow, "slow", false, "")
	flag.BoolVar(&slow, "s", false, "")

	var nbconcurrent int
	flag.IntVar(&nbconcurrent, "nbconcurrent", 8, "")
	flag.IntVar(&nbconcurrent, "n", 8, "")

	var help bool
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&help, "h", false, "")

	flag.Parse()

	if help {
		showHelp()
		return
	}

	args := flag.Args()
	if len(args) > 0 {
		if args[0] == "list" {
			list()
			return
		} else if args[0] == "help" {
			if len(args) == 1 {
				fmt.Println("Please specify the command you need help for")
				os.Exit(1)
			}
			commandHelp(args[1])
			return
		}
	}
	apply(conf, args, state.Options{
		DryRun:             dryrun,
		Slow:               slow,
		NbConcurrentGroups: nbconcurrent,
	})
}

func showHelp() {
	fmt.Print(`
Quickonf usage:

Options:

  -config, -c: path to the configuration file (default quickonf.qconf)
  -dry-run, -d: simulate without modifying the system
  -slow, -s: run slowly (0.5-1s (random) between two instructions - useful for debugging)
  -nbconcurrent, -n: number of concurrent groups running (default 8)
  -help, -h: show this help

  Commands:

  list: list all available commands
  help <command>: get help on the given command

Without any argument: apply all groups from the configuration file
With any number of arguments: apply only the given groups (filtered with patterns) from the configuration file
`)
}
