package arg

import "errors"

type Parser struct {
	args []string
}

func NewParser(args []string) *Parser {
	return &Parser{args: args}
}

func (p *Parser) Pop() (string, error) {
	if len(p.args) == 0 {
		return "", errors.New("args have zero length")
	}
	popped := p.args[0]
	p.args = p.args[1:]
	return popped, nil
}

func (p *Parser) Remaining() (remaining []string) {
	return p.args
}
