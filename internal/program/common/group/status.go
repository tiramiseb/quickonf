package group

type Status int

const (
	StatusWaiting Status = iota
	StatusInfo
	StatusRunning
	StatusFailed
	StatusSucceeded
)
