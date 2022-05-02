package parser

import (
	"fmt"
	"gscript/complier/ast"
	. "gscript/complier/lexer"
	"os"
)

func Parse(l *Lexer) *ast.Program {
	program := parseProgram(l)
	if !l.Expect(TOKEN_EOF) {
		fmt.Println("statement after export is not allowed!!!")
		os.Exit(0)
	}
	return program
}

func parseProgram(l *Lexer) *ast.Program {
	return &ast.Program{
		File:       l.SrcFile(),
		Imports:    parseImports(l),
		BlockStmts: parseBlockStmts(l),
		Export:     parseExport(l),
	}
}

func parseImports(l *Lexer) []ast.Import {
	var imports []ast.Import

	for l.Expect(TOKEN_KW_IMPORT) {
		imports = append(imports, parseImport(l))
		l.ConsumeIf(TOKEN_SEP_SEMI)
	}
	return imports
}

// import net,http as n,h
// import "./localPackage"
func parseImport(l *Lexer) (ipt ast.Import) {
	l.NextToken()
	for {
		ahead := l.LookAhead()
		stdLib := false
		if ahead.Kind == TOKEN_IDENTIFIER {
			stdLib = true
		} else if ahead.Kind != TOKEN_STRING {
			panic(l.Line())
		}
		ipt.Libs = append(ipt.Libs, ast.Lib{
			Stdlib: stdLib,
			Path:   ahead.Value.(string),
		})
		l.NextToken()
		if !l.Expect(TOKEN_SEP_COMMA) {
			break
		}
		l.NextToken()
	}
	if !l.Expect(TOKEN_KW_AS) {
		return
	}
	l.NextToken()
	for i := 0; ; i++ {
		if !l.Expect(TOKEN_IDENTIFIER) {
			panic(l.Line())
		}
		ipt.Libs[i].Alia = l.NextToken().Content
		if !l.Expect(TOKEN_SEP_COMMA) {
			return
		}
		l.NextToken()
	}
}

func parseBlockStmts(l *Lexer) []ast.BlockStmt {
	var blockStmts []ast.BlockStmt
	for {
		switch l.LookAhead().Kind {
		case TOKEN_SEP_LCURLY:
			blockStmts = append(blockStmts, parseBlock(l))
		case TOKEN_KW_EXPORT, TOKEN_EOF, TOKEN_KW_CASE, TOKEN_KW_DEFAULT, TOKEN_SEP_RCURLY:
			return blockStmts
		default:
			blockStmts = append(blockStmts, parseStmt(l))
		}
	}
}

func parseBlock(l *Lexer) (block ast.Block) {
	l.NextTokenKind(TOKEN_SEP_LCURLY)
	block.Blocks = parseBlockStmts(l)
	l.NextTokenKind(TOKEN_SEP_RCURLY)
	return block
}

func parseExport(l *Lexer) (ept ast.Export) {
	if !l.ConsumeIf(TOKEN_KW_EXPORT) {
		return
	}
	ept.Exp = parseExp(l)
	l.ConsumeIf(TOKEN_SEP_SEMI)
	return
}
