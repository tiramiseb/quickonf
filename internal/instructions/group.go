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

	Reports []CheckReport
}

// Run runs the group checks and returns its success status
func (g *Group) Run() bool {
	vars := NewVariablesSet()
	for _, ins := range g.Instructions {
		r, ok := ins.Run(vars)
		g.Reports = append(g.Reports, r...)
		if !ok {
			return false
		}
	}
	return true
}

// HasApply checks if the group has at lease one instruction to apply
func (g *Group) HasApply() bool {
	for _, r := range g.Reports {
		if r.Apply != nil {
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
func (g *Group) Apply(out GroupOutput) bool {
	for _, report := range g.Reports {
		if !report.Apply.Run(out.NewCommandOutput(report.Apply.Name)) {
			return false
		}
	}
	return true
}

func SortGroups(groups []*Group) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Priority > groups[j].Priority
	})
}
