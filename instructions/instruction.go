package instructions

// Instruction is a single instruction
type Instruction interface {
	// Name returns the instruction name
	Name() string
	// RunCheck runs the instruction's check part and return the check reports
	RunCheck(Variables, chan bool) ([]*CheckReport, bool)
	// Reset everything, to have it as it has never run
	Reset()
}
