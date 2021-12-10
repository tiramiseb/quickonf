package state

type If struct {
	Operation    Operation
	Instructions []Instruction
}

func (i *If) Name() string {
	return "if"
}

func (i *If) Run(out Output, vars Variables, options Options) bool {
	success := i.Operation.Compare(vars)
	if !success {
		slow(options)
		out.Infof(`"%s" is false, not running commands...`, i.Operation.String())
		return true
	}
	out.Infof(`"%s" is true, running commands...`, i.Operation.String())
	for _, ins := range i.Instructions {
		subout := out.NewLine(ins.Name())
		if !ins.Run(subout, vars, options) {
			return false
		}
	}
	return true
}
