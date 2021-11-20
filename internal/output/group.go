package output

import (
	"fmt"
	"strings"
)

var (
	nbWaitingGroups  int
	nbRunningGroups  int
	nbFinishedGroups int
	nbFailedGroups   int
)

var groups []*Group

type Group struct {
	name         string
	closed       bool
	instructions []*Instruction

	toDraw []string
}

func NewGroup(name string) *Group {
	nbWaitingGroups--
	nbRunningGroups++
	s := &Group{name: name}
	groups = append(groups, s)
	s.preRedraw()
	redraw()
	return s
}

func (s *Group) NewInstruction(name string) *Instruction {
	instr := &Instruction{group: s, name: name}
	s.instructions = append(s.instructions, instr)
	return instr
}

func (s *Group) preRedraw() {
	name := s.name
	spaces := width - len(name) - 3 // 3 because space + prefix + space
	if spaces < 0 {
		name = name[:len(name)+spaces]
		spaces = 0
	}
	prefix := prefixRunning
	if s.closed {
		prefix = prefixSuccess
		s.toDraw = make([]string, 1)
	} else {
		s.toDraw = make([]string, len(s.instructions)+1)
	}
	s.toDraw[0] = fmt.Sprintf("%s %s %s%s%s", bgBlue, prefix, name, strings.Repeat(" ", spaces), reset)
	if !s.closed {
		for i, instr := range s.instructions {
			s.toDraw[i+1] = instr.draw()
		}
	}
}

func (s *Group) Fail() {
	nbRunningGroups--
	nbFailedGroups++
	s.preRedraw()
	redraw()
}

func (s *Group) Close() {
	s.closed = true
	nbRunningGroups--
	nbFinishedGroups++
	s.preRedraw()
	redraw()
}
