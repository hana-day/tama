package parser

import (
	"github.com/hyusuk/tama/scanner"
	"log"
)

type File struct {
	Exprs []Expr // top-level expressions
	// Scope      *Scope          // package scope (this file only)
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

func (p *Parser) parseInt() Expr {
	return &Primitive{
		Kind:  p.tok,
		Value: p.lit,
	}
}

func (p *Parser) parseExpr() (expr Expr) {
	if p.tok == scanner.INT {
		expr = p.parseInt()
		p.next()
		return
	}
	log.Fatalf("Unexpected token %d", p.tok)
	return nil
}

func (p *Parser) parseExprs() (exprs []Expr) {
	for p.tok != scanner.EOF {
		exprs = append(exprs, p.parseExpr())
	}
	return
}

func (p *Parser) ParseFile() *File {
	return &File{
		Exprs: p.parseExprs(),
	}
}
