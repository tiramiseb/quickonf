package helper

import "sync"

type SearchableStrings struct {
	mutex sync.RWMutex
	elems []string
}

func (s *SearchableStrings) Add(item string) {
	s.mutex.Lock()
	s.elems = append(s.elems, item)
	s.mutex.Unlock()
}

func (s *SearchableStrings) Contains(item string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, elem := range s.elems {
		if elem == item {
			return true
		}
	}
	return false
}
