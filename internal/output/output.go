package output

type Output interface {
	StepTitle(string)
	InstructionTitle(string)
	Info(string)
	Success(string)
	Alert(string)
	Error(error)

	ShowLoader()
	HideLoader()

	ShowPercentage(int)
	HidePercentage()

	Report()
}
