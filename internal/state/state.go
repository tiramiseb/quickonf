package state

import (
	"math/rand"
	"time"
)

// State is a global state, containing groups and optinos
type State struct {
	Options  Options
	Filtered bool
	Groups   []*Group
}

// Group is a list of successive commands
type Group struct {
	Name         string
	Instructions []Instruction
}

// Instruction is a single instruction
type Instruction interface {
	// Name returns the instruction name
	Name() string
	// Run the instruction and return true if it succeeds
	Run(Output, Variables, Options) bool
}

type Output interface {
	NewLine(name string) Output
	Info(message string)
	Infof(format string, a ...interface{})
	Success(message string)
	Successf(format string, a ...interface{})
	Error(message string)
	Errorf(format string, a ...interface{})
}

// Options is a list of options for a state
type Options struct {
	DryRun             bool
	Slow               bool
	NbConcurrentGroups int
}

func slow(options Options) {
	if options.Slow {
		time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
	}
}
