package output

// Output is an output for quickonf
type Output interface {
	StepTitle(string)                         // StepTitle writes a step title
	InstructionTitle(string)                  // InstructionTitle writes an instruction title
	Info(string)                              // Info writes an informational message
	Infof(format string, a ...interface{})    // Infof writes an informational message with formatting
	Success(string)                           // Success writes a successful message
	Successf(format string, a ...interface{}) // Success writes a successful message with formatting
	Alert(string)                             // Alert writes an alert message
	Error(error)                              // Error writes an error message

	ShowLoader() // ShowLoader displays the loader
	HideLoader() // HideLoader hides the loader

	ShowPercentage(int) // ShowPercentage displays the loader with percentage
	HidePercentage()    // HidePercentage hides the loader with percentage

	ShowXonY(int, int) // ShowXonY displays the loader with "X/Y" information
	HideXonY()         // HideXonY hides the loader with "X/Y" information

	Report() // Report writes the summary
}
