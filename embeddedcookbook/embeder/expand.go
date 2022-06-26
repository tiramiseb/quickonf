package main

import (
	"github.com/tiramiseb/quickonf/instructions"
)

func (e *embeder) expand(level int, exp *instructions.Expand) {
	e.write(level, "&Expand{")
	e.write(level+1, "Variable: \"%s\",", exp.Variable)
	e.write(level, "},")
}
