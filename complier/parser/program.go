package parser

import (
	"gscript/complier/ast"
	. "gscript/complier/lexer"
)

func (p *Parser) parseProgram() *ast.Program {
	return &ast.Program{
		File:       p.l.SrcFile(),
		Imports:    p.parseImports(),
		BlockStmts: p.parseBlockStmts(true),
		Export:     p.parseExport(),
	}
}

func (p *Parser) parseImports() []ast.Import {
	var imports []ast.Import

	for p.l.Expect(TOKEN_KW_IMPORT) {
		imports = append(imports, p.parseImport())
		p.l.ConsumeIf(TOKEN_SEP_SEMI)
	}
	return imports
}

// import net,http as n,h
// import "./localPackage"
func (p *Parser) parseImport() (ipt ast.Import) {
	p.l.NextToken()
	for {
		ahead := p.l.LookAhead()
		stdLib := false
		if ahead.Kind == TOKEN_IDENTIFIER {
			stdLib = true
		} else if ahead.Kind != TOKEN_STRING {
			panic(p.l.Line())
		}
		ipt.Libs = append(ipt.Libs, ast.Lib{
			Stdlib: stdLib,
			Path:   ahead.Value.(string),
		})
		p.l.NextToken()
		if !p.l.Expect(TOKEN_SEP_COMMA) {
			break
		}
		p.l.NextToken()
	}
	if !p.l.Expect(TOKEN_KW_AS) {
		return
	}
	p.l.NextToken()
	for i := 0; ; i++ {
		if !p.l.Expect(TOKEN_IDENTIFIER) {
			panic(p.l.Line())
		}
		ipt.Libs[i].Alias = p.l.NextToken().Content
		if !p.l.Expect(TOKEN_SEP_COMMA) {
			return
		}
		p.l.NextToken()
	}
}

func (p *Parser) parseBlockStmts(atTop bool) []ast.BlockStmt {
	var blockStmts []ast.BlockStmt
	for {
		switch p.l.LookAhead().Kind {
		case TOKEN_SEP_LCURLY:
			blockStmts = append(blockStmts, p.parseBlock())
		case TOKEN_KW_EXPORT, TOKEN_EOF, TOKEN_KW_CASE, TOKEN_KW_DEFAULT, TOKEN_SEP_RCURLY:
			return blockStmts
		default:
			blockStmts = append(blockStmts, p.parseStmt(atTop))
		}
	}
}

func (p *Parser) parseBlock() (block ast.Block) {
	p.l.NextTokenKind(TOKEN_SEP_LCURLY)
	block.Blocks = p.parseBlockStmts(false)
	p.l.NextTokenKind(TOKEN_SEP_RCURLY)
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	return block
}

func (p *Parser) parseExport() (ept ast.Export) {
	if !p.l.ConsumeIf(TOKEN_KW_EXPORT) {
		return
	}
	ept.Exp = parseExp(p)
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	return
}
