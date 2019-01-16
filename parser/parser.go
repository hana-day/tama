package parser

import (
	"github.com/hyusuk/tama/scanner"
	"github.com/hyusuk/tama/types"
	"log"
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

func (p *Parser) Init(src []byte) {
	p.scanner.Init(src)
	p.next()
}

func (p *Parser) next() {
	p.tok, p.lit = p.scanner.Scan()
}

func (p *Parser) expect(tok scanner.Token) {
	if p.tok != tok {
		log.Fatalf("expected token %d, but got %d", tok, p.tok)
	}
	p.next()
}

func (p *Parser) parseInt() types.Object {
	f, _ := strconv.ParseFloat(p.lit, 64)
	n := types.Number(f)
	p.next()
	return n
}

func (p *Parser) parseIdent() types.Object {
	sym := &types.Symbol{Name: types.String(p.lit)}
	p.next()
	return sym
}

func (p *Parser) parsePair() types.Object {
	if p.tok == scanner.RPAREN {
		p.next()
		return types.NilObject
	}
	obj := p.parseObject()
	return types.Cons(obj, p.parsePair())
}

func (p *Parser) parseObject() types.Object {
	if p.tok == scanner.INT {
		return p.parseInt()
	}
	if p.tok == scanner.LPAREN {
		p.next()
		return p.parsePair()
	}
	if p.tok == scanner.IDENT {
		return p.parseIdent()
	}
	log.Fatalf("Unexpected token %d", p.tok)
	return nil
}

func (p *Parser) parseObjects() (objs []types.Object) {
	for p.tok != scanner.EOF {
		objs = append(objs, p.parseObject())
	}
	return
}

func (p *Parser) ParseFile() *File {
	return &File{
		Objs: p.parseObjects(),
	}
}
