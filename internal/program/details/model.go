package details

type Model struct{}

func New() *Model {
	return &Model{}
}

func (m *Model) View() string {
	return "details"
}
