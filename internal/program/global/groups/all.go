package groups

import "github.com/tiramiseb/quickonf/internal/instructions"

var all []*instructions.Group

func GetAll() []*instructions.Group {
	return all
}

func CountAll() int {
	return len(all)
}
