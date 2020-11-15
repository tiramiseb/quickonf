package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/tiramiseb/quickonf/internal/quickonf"
)

const conf = "quickonf.yaml"

func main() {
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

	q.Run()
}
