package messages

type ToggleAction int

const (
	ToggleActionDisable = -1
	ToggleActionToggle  = 0
	ToggleActionEnable  = 1
)

type Toggle struct {
	Name   string
	Action ToggleAction
}

type ToggleStatus struct {
	Name   string
	Status bool
}
