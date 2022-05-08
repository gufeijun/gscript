package parser

import (
	"fmt"
	"gscript/complier/ast"
	"gscript/complier/lexer"
	"os"
)

type Parser struct {
	l          *lexer.Lexer
	EnumStmts  []*ast.EnumStmt
	ClassStmts []*ast.ClassStmt
	Labels     []*ast.LabelStmt
	FuncDefs   []*ast.FuncDefStmt
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{
		l: l,
	}
}

func (p *Parser) Parse() *ast.Program {
	program := p.parseProgram()
	if !p.l.Expect(lexer.TOKEN_EOF) {
		fmt.Println("statement after export is not allowed!!!")
		os.Exit(0)
	}
	return program
}
