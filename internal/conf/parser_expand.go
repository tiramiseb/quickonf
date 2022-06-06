package conf

import "github.com/tiramiseb/quickonf/internal/instructions"

func (p *parser) expand(tok *token) instructions.Instruction {
	ins := &instructions.Expand{Variable: tok.content}
	return ins
}
