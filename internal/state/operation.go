package state

type Operation interface {
	Compare(vars Variables) bool
	String() string
}
