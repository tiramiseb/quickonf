package instructions

import (
	"sort"

	"github.com/tiramiseb/quickonf/internal/commands"
)

// GroupOutput makes a new command output when necessary
type GroupOutput interface {
	NewCommandOutput(name string) commands.Output
}

// Group is a list of successive commands
type Group struct {
	Name         string
	Priority     int
	Instructions []Instruction

	Reports []*CheckReport

	alreadyApplied bool
}

// Check runs the group checks and returns its success status
func (g *Group) Check(signalTarget chan bool) bool {
	vars := NewVariablesSet()
	for _, ins := range g.Instructions {
		r, ok := ins.RunCheck(vars, signalTarget)
		g.Reports = append(g.Reports, r...)
		if !ok {
			return false
		}
	}
	return true
}

// Status returns status of the group (according to statuses of its instructions)
func (g *Group) Status() commands.Status {
	var hasInfo, hasRunning, hasSuccess bool
	for _, r := range g.Reports {
		status, _ := r.GetStatusAndMessage()
		switch status {
		case commands.StatusInfo:
			hasInfo = true
		case commands.StatusRunning:
			hasRunning = true
		case commands.StatusSuccess:
			hasSuccess = true
		case commands.StatusError:
			return commands.StatusError
		}
	}
	if hasRunning {
		return commands.StatusRunning
	}
	if hasInfo {
		return commands.StatusInfo
	}
	if hasSuccess {
		return commands.StatusSuccess
	}
	return commands.StatusNotRun
}

// HasApply checks if the group has at lease one instruction to apply
func (g *Group) HasApply() bool {
	for _, r := range g.Reports {
		if r.HasApply() {
			return true
		}
	}
	return false
}

// Reset instructs to reset so that it can re-run later
func (g *Group) Reset() {
	for _, ins := range g.Instructions {
		ins.Reset()
		g.Reports = nil
	}
}

// Apply applies modifications for this group
func (g *Group) Apply() {
	if g.alreadyApplied {
		return
	}
	g.alreadyApplied = true
	for _, report := range g.Reports {
		if !report.Apply() {
			return
		}
	}
}

func SortGroups(groups []*Group) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Priority > groups[j].Priority
	})
}
