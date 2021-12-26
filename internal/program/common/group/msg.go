package group

import "github.com/tiramiseb/quickonf/internal/instructions"

type MsgType int

const (
	CheckTrigger MsgType = iota
	CheckDone
	ApplyChange
	ApplySuccess
	ApplyFail
)

type Msg struct {
	Gidx  int
	Group *instructions.Group
	Type  MsgType
}
