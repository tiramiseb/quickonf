package state

import (
	"time"
)

type Options struct {
	DryRun bool
	Slow   bool
}

func slow(options Options) {
	if options.Slow {
		time.Sleep(500 * time.Millisecond)
	}
}