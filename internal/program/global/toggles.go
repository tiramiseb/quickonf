package global

type toggles map[string]bool
type togglesListeners map[string][]func(new bool)

var (
	Toggles = toggles{
		"details":        false,
		"filter":         true,
		"focusOnDetails": false,
		"help":           false,
		"helpIntro":      true,
		"helpLanguage":   false,
	}
	TogglesListeners = togglesListeners{}
)

// Toggle returns the new value
func (t toggles) Toggle(key string) {
	new := !t[key]
	t[key] = new
	for _, l := range TogglesListeners[key] {
		l(new)
	}
}

func (t toggles) Enable(key string) {
	t[key] = true
	for _, l := range TogglesListeners[key] {
		l(true)
	}
}
func (t toggles) Disable(key string) {
	t[key] = false
	for _, l := range TogglesListeners[key] {
		l(false)
	}
}
