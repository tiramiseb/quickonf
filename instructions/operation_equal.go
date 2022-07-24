package instructions

import (
	"fmt"
)

type Equal struct {
	Left  string
	Right string
}

func (e *Equal) Compare(vars *Variables) bool {
	left := vars.TranslateVariables(e.Left)
	right := vars.TranslateVariables(e.Right)
	return left == right
}

func (e *Equal) String() string {
	return fmt.Sprintf("%s = %s", e.Left, e.Right)
}
