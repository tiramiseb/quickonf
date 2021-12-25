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

	Applys  []commands.Apply
	Reports []CheckReport
}

// Run runs the group checks and returns its success status
func (g *Group) Run() bool {
	vars := NewVariablesSet()
	for _, ins := range g.Instructions {
		a, r, ok := ins.Run(vars)
		g.Applys = append(g.Applys, a...)
		g.Reports = append(g.Reports, r...)
		if !ok {
			return false
		}
	}
	return true
}

// Reset instructs to reset so tat it can re-run later
func (g *Group) Reset() {
	for _, ins := range g.Instructions {
		ins.Reset()
		g.Applys = nil
		g.Reports = nil
	}
}

// Apply applies modifications for this group
func (g *Group) Apply(out GroupOutput) bool {
	for _, apply := range g.Applys {
		if !apply.Run(out.NewCommandOutput(apply.Name)) {
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
