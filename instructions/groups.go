package instructions

import (
	"sort"
	"sync"

	"github.com/tiramiseb/quickonf/commands"
)

type Groups struct {
	groups []*Group

	maxNameLength int
	count         int

	initialChecksOnce sync.Once
}

func NewGroups(groups []*Group) *Groups {
	// Sort groups by priority
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Priority > groups[j].Priority
	})
	var maxNameLength int
	var prevGroup *Group
	for _, g := range groups {
		if prevGroup != nil {
			prevGroup.next = g
			g.previous = prevGroup
		}
		l := len(g.Name)
		if l > maxNameLength {
			maxNameLength = l
		}
		prevGroup = g
	}
	return &Groups{
		groups:        groups,
		maxNameLength: maxNameLength,
		count:         len(groups),
	}
}

func (g *Groups) InitialChecks(signalTarget chan bool) {
	g.initialChecksOnce.Do(func() {
		// Groups are already ordered by priority, just take them one after another, and wait when priority changes.
		if len(g.groups) == 0 {
			return
		}
		currentPriority := g.groups[0].Priority
		var wg sync.WaitGroup
		for _, g := range g.groups {
			if g.Priority != currentPriority {
				wg.Wait()
				currentPriority = g.Priority
				wg = sync.WaitGroup{}
			}
			wg.Add(1)
			thisGroup := g
			go func() {
				thisGroup.Check(signalTarget, false)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func (g *Groups) All() []*Group {
	return g.groups
}

func (g *Groups) ApplyAll() {
	currentPriority := g.groups[0].Priority
	var wg sync.WaitGroup
	for _, g := range g.groups {
		if g.Priority != currentPriority {
			wg.Wait()
			currentPriority = g.Priority
			wg = sync.WaitGroup{}
		}
		if g.Status() == commands.StatusSuccess || g.Status() == commands.StatusError {
			continue
		}
		wg.Add(1)
		thisGroup := g
		go func() {
			thisGroup.Apply()
			wg.Done()
		}()
	}
	wg.Wait()
}

func (g *Groups) FirstGroup() *Group {
	if len(g.groups) == 0 {
		return nil
	}
	return g.groups[0]
}

func (g *Groups) MaxNameLength() int {
	return g.maxNameLength
}

func (g *Groups) Count() int {
	return g.count
}
