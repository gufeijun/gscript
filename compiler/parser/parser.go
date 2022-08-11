package parser

import (
	"fmt"
	"gscript/compiler/ast"
	"gscript/compiler/lexer"
	"gscript/compiler/token"
	"os"
)

type Parser struct {
	l          *lexer.Lexer
	EnumStmts  []*ast.EnumStmt
	ClassStmts []*ast.ClassStmt
	FuncDefs   []*ast.FuncDefStmt
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{
		l: l,
	}
}

func (p *Parser) Parse() *ast.Program {
	program := p.parseProgram()
	if !p.Expect(token.TOKEN_EOF) {
		fmt.Println("statement after export is not allowed!!!")
		os.Exit(0)
	}
	return program
}

func (p *Parser) NextTokenKind(kind int) *token.Token {
	t := p.l.NextToken()
	if t.Kind != kind {
		p.exit("expect token '%s', but got '%s'", token.TokenDescs[kind], t.Content)
	}
	return t

}

func (p *Parser) ConsumeIf(kind int) bool {
	if p.Expect(kind) {
		p.l.NextToken()
		return true
	}
	return false
}

func (p *Parser) Expect(kind int) bool {
	return p.l.LookAhead().Kind == kind
}

func (p *Parser) exit(format string, args ...interface{}) {
	errMsg := fmt.Sprintf(format, args...)
	fmt.Printf("Parser error: [%s:%d:%d]\n", p.l.SrcFile(), p.l.Line(), p.l.Column())
	fmt.Printf("\t%s\n", errMsg)
	os.Exit(0)
}
