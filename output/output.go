package output

type Output interface {
	StepTitle(string)
	ModuleName(string)
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
