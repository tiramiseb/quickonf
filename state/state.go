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

func (s *State) Run(options Options) {
	nb := len(s.Groups)
	if s.Filtered {
		if nb > 1 {
			output.SetTitle(fmt.Sprintf("Applying %d steps", nb))
		} else {
			output.SetTitle(fmt.Sprintf("Applying %d step", nb))

		}
	} else {
		output.SetTitle("Applying all steps")
	}
	output.Start(nb)
	defer output.End()
	limit := semaphore.NewWeighted(8)
	limitCtx := context.Background()
	for _, group := range s.Groups {
		limit.Acquire(limitCtx, 1)
		gr := group
		gr.variables = newVariablesSet()
		go func() {
			gr.Run(options)
			limit.Release(1)
		}()
	}
	limit.Acquire(limitCtx, 8)
}
