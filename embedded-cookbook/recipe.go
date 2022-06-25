package main

import (
	"github.com/tiramiseb/quickonf/instructions"
)

func (e *embeder) recipe(level int, rec *instructions.Recipe) {
	e.write(level, "&Recipe{")
	e.write(level+1, "RecipeName: \"%s\",", rec.RecipeName)
	if len(rec.Variables) > 0 {
		e.write(level+1, "Variables: map[string]string{")
		for k, v := range rec.Variables {
			e.write(level+2, "`%s`: `%s`,", k, v)
		}
		e.write(level+1, "},")
	}
	e.write(level, "},")
}
