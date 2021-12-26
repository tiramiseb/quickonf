package instructions

import "github.com/tiramiseb/quickonf/internal/commands"

// CheckReport is a single report after checking is something must be applied
type CheckReport struct {
	Name    string
	Status  commands.Status
	Message string
}

// Instruction is a single instruction
type Instruction interface {
	// Name returns the instruction name
	Name() string
	// Run the instruction and return the apply functions (or nil), check reports (even if there is nothing to apply) and true if it succeeds
	Run(Variables) ([]commands.Apply, []CheckReport, bool)
	// Reset everything, to have it as it has never run
	Reset()
}
