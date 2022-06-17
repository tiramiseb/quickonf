package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) expand(tok *token) instructions.Instruction {
	ins := &instructions.Expand{Variable: tok.content}
	return ins
}
