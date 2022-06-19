package instructions

import (
	"github.com/tiramiseb/quickonf/commands"
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

	previous *Group
	next     *Group
}

// Check runs the group checks and returns its success status
func (g *Group) Check(signalTarget chan bool, reset bool) bool {
	// TODO Is the return value really used?
	vars := NewVariablesSet()
	g.Reports = g.Reports[:0]
	if signalTarget != nil {
		signalTarget <- true
	}
	for _, ins := range g.Instructions {
		if reset {
			ins.Reset()
		}
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

// Get the n-th previous group
func (g *Group) Previous(n int, includeSuccess bool) *Group {
	switch {
	case !includeSuccess && g.Status() == commands.StatusSuccess:
		// Success unwanted and I am in success: pass to the previous one
		if g.previous == nil {
			if g.next == nil {
				// I am the only one
				return g
			}
			// I am the first, return the next
			return g.next.nextByDefault()
		}
		return g.previous.Previous(n, false)
	case n == 0 || g.previous == nil:
		// It is me you are looking for!
		return g
	default:
		return g.previous.Previous(n-1, includeSuccess)
	}
}

// previousByDefault returns the previous non-success group, because there was no next one
func (g *Group) previousByDefault() *Group {
	switch {
	case g.Status() != commands.StatusSuccess:
		return g
	case g.previous == nil:
		return nil
	default:
		return g.previous.previousByDefault()
	}
}

// Get the n-th next group
func (g *Group) Next(n int, includeSuccess bool) *Group {
	switch {
	case !includeSuccess && g.Status() == commands.StatusSuccess:
		// Success unwanted and I am in success: pass to the next one
		if g.next == nil {
			if g.previous == nil {
				// I am the only one
				return g
			}
			// I am the last, return the previous
			return g.previous.previousByDefault()
		}
		return g.next.Next(n, false)
	case n == 0 || g.next == nil:
		// It is me you are looking for!
		return g
	default:
		return g.next.Next(n-1, includeSuccess)
	}
}

// nextByDefault returns the next non-success group, because there was no previous one
func (g *Group) nextByDefault() *Group {
	switch {
	case g.Status() != commands.StatusSuccess:
		return g
	case g.next == nil:
		return nil
	default:
		return g.next.nextByDefault()
	}
}