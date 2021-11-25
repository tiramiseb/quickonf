package state

import (
	"github.com/tiramiseb/quickonf/internal/output"
)

type Command interface {
	// Run the command and return true if it succeeds
	Run(*output.Group, variables, Options) bool
}
