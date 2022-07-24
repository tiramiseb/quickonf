package conf

func (p *parser) cookbook(line tokens) (next tokens) {
	if len(line) == 1 {
		p.errs = append(p.errs, line[0].errorf(`cookbook URI not provided`))
		return p.nextLine()
	}
	p.cookbooks = append(p.cookbooks, line[1].content)
	return p.nextLine()

}
