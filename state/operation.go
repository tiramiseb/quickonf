package state

type Operation interface {
	Compare(vars variables) bool
	String() string
}
