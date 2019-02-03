package parser

import (
	"github.com/hyusuk/tama/scanner"
	"github.com/hyusuk/tama/types"
	"strconv"
)

type File struct {
	Objs []types.Object // top-level expressions
}

type Parser struct {
	scanner scanner.Scanner
	tok     scanner.Token // Next token
	lit     string        // Next token literal
}

func (p *Parser) Init(src []byte) error {
	p.scanner.Init(src)
	if err := p.next(); err != nil {
		return err
	}
	return nil
}

func (p *Parser) next() error {
	var err error
	if p.tok, p.lit, err = p.scanner.Scan(); err != nil {
		return err
	}
	return nil
}

func (p *Parser) expect(tok scanner.Token) error {
	if p.tok != tok {
		return types.NewSyntaxError("expected token %d, but got %d", tok, p.tok)
	}
	if err := p.next(); err != nil {
		return err
	}
	return nil
}

func (p *Parser) parseFloat() (types.Object, error) {
	f, err := strconv.ParseFloat(p.lit, 64)
	if err != nil {
		return nil, types.NewSyntaxError("cannot parse %s as a number", p.lit)
	}
	n := types.Number(f)
	return n, p.next()
}

func (p *Parser) parseIdent() (types.Object, error) {
	sym := types.NewSymbol(p.lit)
	return sym, p.next()
}

func (p *Parser) parsePair() (types.Object, error) {
	if p.tok == scanner.RPAREN {
		return types.NilObject, p.next()
	}
	car, err := p.parseObject()
	if err != nil {
		return nil, err
	}
	cdr, err := p.parsePair()
	if err != nil {
		return nil, err
	}
	return types.Cons(car, cdr), nil
}

func (p *Parser) parseVector() (types.Object, error) {
	var v types.Vector
	for p.tok != scanner.RPAREN {
		o, err := p.parseObject()
		if err != nil {
			return nil, err
		}
		v = append(v, o)
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	return v, nil
}

func (p *Parser) parseString() (types.Object, error) {
	s := types.String(p.lit)
	if err := p.next(); err != nil {
		return nil, err
	}
	return s, nil
}

func (p *Parser) parseObject() (types.Object, error) {
	tok := p.tok
	switch tok {
	case scanner.NUMBER:
		return p.parseFloat()
	case scanner.LPAREN:
		if err := p.next(); err != nil {
			return nil, err
		}
		return p.parsePair()
	case scanner.VLPAREN:
		if err := p.next(); err != nil {
			return nil, err
		}
		return p.parseVector()
	case scanner.IDENT:
		return p.parseIdent()
	case scanner.QUOTE: // '(1 2 3) => (quote (1 2 3))
		if err := p.next(); err != nil {
			return nil, err
		}
		obj, err := p.parseObject()
		if err != nil {
			return nil, err
		}
		return types.List(types.NewSymbol("quote"), obj), nil
	case scanner.TRUE, scanner.FALSE:
		if err := p.next(); err != nil {
			return nil, err
		}
		if tok == scanner.TRUE {
			return types.Boolean(true), nil
		}
		return types.Boolean(false), nil
	case scanner.STRING:
		return p.parseString()
	default:
		return nil, types.NewSyntaxError("unexpected token %d", p.tok)

	}
}

func (p *Parser) parseObjects() ([]types.Object, error) {
	var objs []types.Object
	for p.tok != scanner.EOF {
		obj, err := p.parseObject()
		if err != nil {
			return objs, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func (p *Parser) ParseFile() (*File, error) {
	objs, err := p.parseObjects()
	if err != nil {
		return nil, err
	}
	return &File{Objs: objs}, nil
}
