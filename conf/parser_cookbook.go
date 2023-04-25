package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) cookbook(line tokens) (next tokens) {
	if len(line) == 1 {
		p.groups = append(p.groups, &instructions.Group{
			Name: "Cookbook",
			Instructions: []instructions.Instruction{
				line[0].error("cookbook URI not provided"),
			},
		})
		p.checkResult.addError(line[0], CheckSeverityError, "cookbook URI not provided")
		return p.nextLine()
	}
	p.cookbooks = append(p.cookbooks, line[1].content)
	return p.nextLine()

}
