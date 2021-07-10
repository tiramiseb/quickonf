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
	flag.StringVar(&conf, "config", "quickonf.yaml", "List all steps")
	flag.StringVar(&conf, "c", "quickonf.yaml", "List all steps (shorthand)")
	var list bool
	flag.BoolVar(&list, "list", false, "List all steps")
	flag.BoolVar(&list, "l", false, "List all steps (shorthand)")
	var dryrun bool
	flag.BoolVar(&dryrun, "dry-run", false, "Try all steps without modifying the system")
	flag.BoolVar(&dryrun, "r", false, "Try all steps without modifying the system (shorthand)")
	var output string
	flag.StringVar(&output, "output", "stdout", "Output format")
	flag.StringVar(&output, "o", "stdout", "Output format (shorthand)")

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

	q, err := quickonf.New(config, output)
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
	// Check the patterns before running steps
	for i, s := range steps {
		s = "*" + s + "*"
		if _, err := path.Match(s, ""); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		steps[i] = s
	}
	q.Run(steps, dryrun)
}
