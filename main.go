package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/tiramiseb/quickonf/internal/quickonf"
)

func main() {
	var conf string
	var list bool
	var dryrun bool
	flag.StringVar(&conf, "config", "quickonf.yaml", "List all steps")
	flag.StringVar(&conf, "c", "quickonf.yaml", "List all steps (shorthand)")
	flag.BoolVar(&list, "list", false, "List all steps")
	flag.BoolVar(&list, "l", false, "List all steps (shorthand)")
	flag.BoolVar(&dryrun, "dry-run", false, "Try all steps without modifying the system")
	flag.BoolVar(&dryrun, "r", false, "Try all steps without modifying the system (shorthand)")
	flag.Parse()

	yamlFile, err := ioutil.ReadFile(conf)
	if err != nil {
		fmt.Println("Could not read " + conf)
		fmt.Println(err)
		os.Exit(1)
	}

	var config []quickonf.Step
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Wrong file format!")
		fmt.Println(err)
		os.Exit(1)
	}

	q, err := quickonf.New(config)
	if err != nil {
		fmt.Println("Could not initialize quickonf")
		fmt.Println(err)
		os.Exit(1)
	}

	if list {
		q.List()
		return
	}

	steps := flag.Args()
	if len(steps) == 0 {
		q.Run(dryrun)
	} else {
		for i, s := range steps {
			s = "*" + s + "*"
			if _, err := path.Match(s, ""); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			steps[i] = s
		}
		q.Steps(steps, dryrun)
	}
}
