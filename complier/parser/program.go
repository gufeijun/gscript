package parser

import (
	"gscript/complier/ast"
	"gscript/complier/token"
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

	for p.Expect(token.TOKEN_KW_IMPORT) {
		imports = append(imports, p.parseImport())
		p.ConsumeIf(token.TOKEN_SEP_SEMI)
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
		if ahead.Kind == token.TOKEN_IDENTIFIER {
			stdLib = true
		} else if ahead.Kind != token.TOKEN_STRING {
			p.exit("expect path(string) after keyword \"import\"")
		}
		ipt.Libs = append(ipt.Libs, ast.Lib{
			Stdlib: stdLib,
			Path:   ahead.Value.(string),
		})
		p.l.NextToken()
		if !p.Expect(token.TOKEN_SEP_COMMA) {
			break
		}
		p.l.NextToken()
	}
	if !p.Expect(token.TOKEN_KW_AS) {
		return
	}
	p.l.NextToken()
	for i := 0; ; i++ {
		if !p.Expect(token.TOKEN_IDENTIFIER) {
			p.exit("expect alias(identifier) for imported package '%s' after keyword 'as'", ipt.Libs[i].Path)
		}
		ipt.Libs[i].Alias = p.l.NextToken().Content
		if !p.Expect(token.TOKEN_SEP_COMMA) {
			return
		}
		p.l.NextToken()
	}
}

func (p *Parser) parseBlockStmts(atTop bool) []ast.BlockStmt {
	var blockStmts []ast.BlockStmt
	for {
		switch p.l.LookAhead().Kind {
		case token.TOKEN_SEP_LCURLY:
			blockStmts = append(blockStmts, p.parseBlock())
		case token.TOKEN_KW_EXPORT, token.TOKEN_EOF, token.TOKEN_KW_CASE, token.TOKEN_KW_DEFAULT, token.TOKEN_SEP_RCURLY:
			return blockStmts
		default:
			blockStmts = append(blockStmts, p.parseStmt(atTop))
		}
	}
}

func (p *Parser) parseBlock() (block ast.Block) {
	p.NextTokenKind(token.TOKEN_SEP_LCURLY)
	block.Blocks = p.parseBlockStmts(false)
	p.NextTokenKind(token.TOKEN_SEP_RCURLY)
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	return block
}

func (p *Parser) parseExport() (ept ast.Export) {
	if !p.ConsumeIf(token.TOKEN_KW_EXPORT) {
		return
	}
	ept.Exp = parseExp(p)
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	return
}
