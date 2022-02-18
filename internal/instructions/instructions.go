package instructions

import "github.com/tiramiseb/quickonf/internal/commands"

// CheckReport is a single report after checking is something must be applied
type CheckReport struct {
	Name    string
	Status  commands.Status
	Message string
	Apply   *commands.Apply
}

// Instruction is a single instruction
type Instruction interface {
	// Name returns the instruction name
	Name() string
	// Run the instruction and return the check reports
	Run(Variables) ([]CheckReport, bool)
	// Reset everything, to have it as it has never run
	Reset()
}
