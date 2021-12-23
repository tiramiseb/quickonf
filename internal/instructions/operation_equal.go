package instructions

import (
	"fmt"
	"strings"
)

type Equal struct {
	Left  string
	Right string
}

func (e *Equal) Compare(vars Variables) bool {
	left := vars.translateVariables(e.Left)
	right := vars.translateVariables(e.Right)

	// At first, only compare strings
	return strings.EqualFold(left, right)
}

func (e *Equal) String() string {
	return fmt.Sprintf("%s = %s", e.Left, e.Right)
}
