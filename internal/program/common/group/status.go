package group

type Status int

const (
	StatusWaiting Status = iota
	StatusRunning
	StatusFailed
	StatusSucceeded
)
