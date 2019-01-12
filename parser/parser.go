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
	expr := &Primitive{
		Kind:  p.tok,
		Value: p.lit,
	}
	p.next()
	return expr
}

func (p *Parser) parseIdent() Expr {
	expr := &Ident{Name: p.lit}
	p.next()
	return expr
}

func (p *Parser) parseCall() Expr {
	if p.tok != scanner.IDENT {
		log.Fatalf("Unexpected token %d", p.tok)
	}
	expr := &CallExpr{Func: &Ident{Name: p.lit}}
	p.next()
	for p.tok != scanner.RPAREN {
		expr.Args = append(expr.Args, p.parseExpr())
	}
	p.next()
	return expr
}

func (p *Parser) parseExpr() (expr Expr) {
	if p.tok == scanner.INT {
		expr = p.parseInt()
		return
	}
	if p.tok == scanner.LPAREN {
		p.next()
		expr = p.parseCall()
		return
	}
	if p.tok == scanner.IDENT {
		expr = p.parseIdent()
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
