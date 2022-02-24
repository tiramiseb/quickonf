package instructions

// Instruction is a single instruction
type Instruction interface {
	// Name returns the instruction name
	Name() string
	// Run the instruction and return the check reports
	Run(Variables, chan bool) ([]*CheckReport, bool)
	// Reset everything, to have it as it has never run
	Reset()
}
