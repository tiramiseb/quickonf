package output

import "fmt"

type Prog struct{}

// NewProg returns a new output writer which writes on stdout for programmatic use
func NewProg() Output {
	return &Prog{}
}

// StepTitle writes a step title
func (p *Prog) StepTitle(str string) {
	fmt.Println("STEP:" + str)
}

// InstructionTitle writes an instruction title
func (p *Prog) InstructionTitle(str string) {
	fmt.Println("INSTRUCTION:" + str)
}

// Info writes an informational message
func (p *Prog) Info(str string) {
	fmt.Println("INFO:" + str)
}

// Infof writes an informational message with format
func (p *Prog) Infof(format string, a ...interface{}) {
	fmt.Printf("INFO:"+format+"\n", a...)
}

// Success writes a successful message
func (p *Prog) Success(str string) {
	fmt.Println("SUCCESS:" + str)
}

// Successf writes a successful message with format
func (p *Prog) Successf(format string, a ...interface{}) {
	fmt.Printf("SUCCESS:"+format+"\n", a...)
}

// Alert writes an alert message
func (p *Prog) Alert(str string) {
	fmt.Println("ALERT:" + str)
}

// Alertf writes an alert message with format
func (p *Prog) Alertf(format string, a ...interface{}) {
	fmt.Printf("ALERT:"+format+"\n", a...)
}

// Error writes an error message
func (p *Prog) Error(err error) {
	fmt.Println("ERROR:" + err.Error())
}

// ShowLoader instruct to display a loader
func (p *Prog) ShowLoader() {
	fmt.Println("LOADER")
}

// HideLoader instructs to hide the loader
func (p *Prog) HideLoader() {
	fmt.Println("HIDE LOADER")
}

// ShowPercentage instructs to display a loader with given percentage
func (p *Prog) ShowPercentage(pc int) {
	fmt.Printf("PERCENTAGE:%d\n", pc)
}

// HidePercentage instructs to hide the loader with percentage
func (p *Prog) HidePercentage() {
	fmt.Println("HIDE PERCENTAGE")
}

// ShowXonY instructs to display a loader with with "X/Y" information
func (p *Prog) ShowXonY(x, y int) {
	fmt.Printf("XY:%d:%d\n", x, y)
}

// HideXonY instructs to hide the loader with "X/Y" information
func (p *Prog) HideXonY() {
	fmt.Println("HIDE XY")
}

// Report does not write anything
func (p *Prog) Report() {}
