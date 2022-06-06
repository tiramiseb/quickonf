package instructions

import (
	"fmt"
)

type Different struct {
	Left  string
	Right string
}

func (d *Different) Compare(vars Variables) bool {
	left := vars.translateVariables(d.Left)
	right := vars.translateVariables(d.Right)

	return left != right

}

func (d *Different) String() string {
	return fmt.Sprintf("%s != %s", d.Left, d.Right)
}
