package messages

type Apply struct{}

type ApplyAll struct{}

type ConfirmApplyAll struct{}

type Details struct{}

type Filter struct{}

type Help struct{}

type NewSignal struct{}

type Recheck struct{}

type ToggleStatus struct {
	Name   string
	Status bool
}
