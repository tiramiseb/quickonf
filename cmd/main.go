package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tiramiseb/quickonf/internal/quickonf"
)

const conf = "quickonf.yaml"

func main() {
	var list bool
	var dryrun bool
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
	q.Run(dryrun)
}
