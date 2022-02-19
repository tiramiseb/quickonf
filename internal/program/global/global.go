package global

type global struct {
	values map[string]bool
}

var Global = &global{
	values: map[string]bool{},
}

func (g *global) Set(key string, value bool) {
	g.values[key] = value
}

func (g *global) Get(key string) bool {
	return g.values[key]
}
