package state

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"

	"github.com/tiramiseb/quickonf/internal/output"
)

type State struct {
	Filtered bool
	Groups   []*Group
}

func (s *State) Run() {
	if s.Filtered {
		output.SetTitle(fmt.Sprintf("Applying %d steps", len(s.Groups)))
	} else {
		output.SetTitle("Applying all steps")
	}
	output.Start(len(s.Groups))
	defer output.End()
	limit := semaphore.NewWeighted(8)
	limitCtx := context.Background()
	for _, group := range s.Groups {
		limit.Acquire(limitCtx, 1)
		gr := group
		gr.variables = newVariablesSet()
		go func() {
			gr.Run()
			limit.Release(1)
		}()
	}
	limit.Acquire(limitCtx, 8)
}
