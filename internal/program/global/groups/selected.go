package groups

import "github.com/tiramiseb/quickonf/internal/instructions"

func GetSelectedIndex() int {
	return selected
}

func IncrementSelected(step int) {
	selected += step
	if selected >= len(displayed) {
		selected = len(displayed) - 1
	}
}

func DecrementSelected(step int) {
	selected -= step
	if selected < 0 {
		selected = 0
	}
}

func SelectFirst() {
	selected = 0
}

func SelectLast() {
	selected = len(displayed) - 1
}

func GetSelected() *instructions.Group {
	if len(displayed) == 0 {
		return nil
	}
	return displayed[selected]
}
