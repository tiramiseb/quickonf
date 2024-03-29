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

type docRecipe struct {
	Slug         string
	Name         string
	Doc          string
	VarsDoc      map[string]string
	Instructions string
}

func main() {

	// Commands
	for _, cmd := range commands.GetAll() {
		cmdYAML, err := yaml.Marshal(cmd)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile("data/commands/"+cmd.Name+".yaml", cmdYAML, 0o644); err != nil {
			panic(err)
		}
	}

	// Cookbook
	if err := embeddedcookbook.ForEach(func(recipe *instructions.Group) error {
		var instructions string

		for _, instr := range recipe.Instructions {
			instructions = instructions + "\n" + instr.String()
		}
		short := docRecipe{
			Slug:         slug.Make(recipe.Name),
			Name:         recipe.Name,
			Doc:          recipe.RecipeDoc,
			VarsDoc:      recipe.RecipeVarsDoc,
			Instructions: instructions,
		}
		cmdYAML, err := yaml.Marshal(short)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile("data/cookbook/"+slug.Make(short.Name)+".yaml", cmdYAML, 0o644); err != nil {
			panic(err)
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
