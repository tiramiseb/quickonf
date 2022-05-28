package groups

import (
	"sync"

	"github.com/tiramiseb/quickonf/internal/instructions"
)

var (
	displayed   []*instructions.Group
	displayedMu sync.Mutex
)

func GetDisplayed() []*instructions.Group {
	displayedMu.Lock()
	defer displayedMu.Unlock()
	return displayed
}

func CountDisplayed() int {
	displayedMu.Lock()
	defer displayedMu.Unlock()
	return len(displayed)
}
