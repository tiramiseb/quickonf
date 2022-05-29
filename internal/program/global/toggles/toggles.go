package toggles

var (
	toggles   = map[string]bool{}
	listeners = map[string][]func(new bool){}
)

func Get(key string) bool {
	return toggles[key]
}

func Toggle(key string) {
	new := !toggles[key]
	toggles[key] = new
	for _, l := range listeners[key] {
		l(new)
	}
}

func Enable(key string) {
	toggles[key] = true
	for _, l := range listeners[key] {
		l(true)
	}
}
func Disable(key string) {
	toggles[key] = false
	for _, l := range listeners[key] {
		l(false)
	}
}
