package instructions

import (
	"fmt"
	"strings"
)

type Different struct {
	Left  string
	Right string
}

func (d *Different) Compare(vars Variables) bool {
	left := vars.translateVariables(d.Left)
	right := vars.translateVariables(d.Right)

	// At first, only compare strings
	return !strings.EqualFold(left, right)
}

func (d *Different) String() string {
	return fmt.Sprintf("%s != %s", d.Left, d.Right)
}
