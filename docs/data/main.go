package main

//go:generate go run .

import (
	"os"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"

	"github.com/tiramiseb/quickonf/commands"
	"github.com/tiramiseb/quickonf/embeddedcookbook"
	"github.com/tiramiseb/quickonf/instructions"
)

type shortRecipe struct {
	Name    string
	Doc     string
	VarsDoc map[string]string
}

func main() {

	// Commands
	for _, cmd := range commands.GetAll() {
		cmdYAML, err := yaml.Marshal(cmd)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile("commands/"+cmd.Name+".yaml", cmdYAML, 0o644); err != nil {
			panic(err)
		}
	}

	// Cookbook
	if err := embeddedcookbook.ForEach(func(recipe *instructions.Group) error {
		short := shortRecipe{
			Name:    recipe.Name,
			Doc:     recipe.RecipeDoc,
			VarsDoc: recipe.RecipeVarsDoc,
		}
		cmdYAML, err := yaml.Marshal(short)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile("cookbook/"+slug.Make(short.Name)+".yaml", cmdYAML, 0o644); err != nil {
			panic(err)
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
