package main

//go:generate go run . doc

import (
	"os"

	"github.com/tiramiseb/quickonf/internal/commands"
	"gopkg.in/yaml.v3"
)

func makeDocCommands() {
	for _, cmd := range commands.GetAll() {
		cmdYAML, err := yaml.Marshal(cmd)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile("data/commands/"+cmd.Name+".yaml", cmdYAML, 0644); err != nil {
			panic(err)
		}
	}
}
