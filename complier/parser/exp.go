package parser

import (
	"gscript/complier/ast"
	"gscript/complier/token"
)

/*
exp ::= | '(' exp ')' | literal | nil | new ID ['(' [expList] ')'] | ID | unOP exp
        | exp '--'
        | exp '++'
        | exp '.' ID
        | exp '[' exp ']'
        | exp binOP exp
        | exp '?' exp ':' exp
        | exp '(' [expList] ')'
*/

/*
exp    ::= term12
term12 ::= term11 [ '?' term11 ':' term11]	 // nesting of ternary operators is not allowed
term11 ::= term10 { '||' term10 }
term10 ::= term9  { '&&' term9 }
term9  ::= term8  { '|' term8 }
term8  ::= term7  { '^' term7 }
term7  ::= term6  { '&' term6 }
term6  ::= term5  { ('==' | '!=') term5 }
term5  ::= term4  { ( '>' | '<' | '>=' | '<=') term4 }
term4  ::= term3  { ( '<<' | '>>' ) term3 }
term3  ::= term2  { ( '+' | '-') term2 }
term2  ::= term1  { ( '/' | '*' | '%' | '//' ) term1 }
term1  ::= { ( '-' | '~  | '!' ) } term0
term0  ::= factor | ( '++' | '--') factor | factor ('++' | '--')
factor ::= '(' exp ')' | literal | nil | new ID ['(' [expList] ')']
		 | ID { ( '.' ID | '[' exp ']' | '(' [expList] ')' ) }
*/

// '(' [explist] ')
func parseExpListBlock(p *Parser) []ast.Exp {
	p.NextTokenKind(token.TOKEN_SEP_LPAREN)
	if p.Expect(token.TOKEN_SEP_RPAREN) {
		p.l.NextToken()
		return nil
	}
	exps := parseExpList(p)
	p.NextTokenKind(token.TOKEN_SEP_RPAREN)
	return exps
}

func parseExpList(p *Parser) (exps []ast.Exp) {
	exps = append(exps, parseExp(p))
	for p.Expect(token.TOKEN_SEP_COMMA) {
		p.l.NextToken()
		exps = append(exps, parseExp(p))
	}
	return
}

func parseExp(p *Parser) ast.Exp {
	return parseTerm12(p)
}

func parseTerm12(p *Parser) ast.Exp {
	exp1 := parseTerm11(p)
	if !p.Expect(token.TOKEN_SEP_QMARK) { // ?
		return exp1
	}
	p.l.NextToken()
	exp2 := parseTerm11(p)
	p.NextTokenKind(token.TOKEN_SEP_COLON) // :
	return &ast.TernaryOpExp{
		Exp1: exp1,
		Exp2: exp2,
		Exp3: parseTerm11(p),
	}
}

func parseTerm11(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_LOR}, parseTerm10)
}

func parseTerm10(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_LAND}, parseTerm9)
}

func parseTerm9(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_OR}, parseTerm8)
}

func parseTerm8(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_XOR}, parseTerm7)
}

func parseTerm7(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_AND}, parseTerm6)
}

func parseTerm6(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_EQ, token.TOKEN_OP_NE}, parseTerm5)
}

func parseTerm5(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_LE, token.TOKEN_OP_GE, token.TOKEN_OP_LT, token.TOKEN_OP_GT}, parseTerm4)
}

func parseTerm4(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_SHL, token.TOKEN_OP_SHR}, parseTerm3)
}

func parseTerm3(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_ADD, token.TOKEN_OP_SUB}, parseTerm2)
}

func parseTerm2(p *Parser) ast.Exp {
	return _parseBinExp(p, []int{token.TOKEN_OP_DIV, token.TOKEN_OP_MUL, token.TOKEN_OP_MOD, token.TOKEN_OP_IDIV}, parseTerm1)
}

func parseTerm1(p *Parser) ast.Exp {
	kind := p.l.LookAhead().Kind

	var unOp int
	if kind == token.TOKEN_OP_SUB {
		unOp = ast.UNOP_NEG
	} else if kind == token.TOKEN_OP_NOT {
		unOp = ast.UNOP_NOT
	} else if kind == token.TOKEN_OP_LNOT {
		unOp = ast.UNOP_LNOT
	} else {
		return parseTerm0(p)
	}
	p.l.NextToken()
	return &ast.UnOpExp{Op: unOp, Exp: parseTerm1(p)}
}

func parseTerm0(p *Parser) ast.Exp {
	if p.Expect(token.TOKEN_OP_INC) {
		p.l.NextToken()
		return &ast.UnOpExp{Op: ast.UNOP_INC, Exp: parseFactor(p)}
	} else if p.Expect(token.TOKEN_OP_DEC) {
		p.l.NextToken()
		return &ast.UnOpExp{Op: ast.UNOP_DEC, Exp: parseFactor(p)}
	}
	exp := parseFactor(p)
	if p.Expect(token.TOKEN_OP_INC) {
		p.l.NextToken()
		return &ast.UnOpExp{Op: ast.UNOP_INC_, Exp: exp}
	} else if p.Expect(token.TOKEN_OP_DEC) {
		p.l.NextToken()
		return &ast.UnOpExp{Op: ast.UNOP_DEC_, Exp: exp}
	}
	return exp
}

// factor ::= '(' exp ')' | literal | nil | new ID ['(' [expList] ')']
// 		 | ID { ( '.' ID | '[' exp ']' | '(' [expList] ')' ) }
func parseFactor(p *Parser) ast.Exp {
	switch ahead := p.l.LookAhead(); ahead.Kind {
	case token.TOKEN_STRING:
		return parseStringLiteralExp(p)
	case token.TOKEN_NUMBER:
		return parseNumberLiteralExp(p)
	case token.TOKEN_SEP_LCURLY: // mapLiteral
		return parseMapLiteralExp(p)
	case token.TOKEN_SEP_LBRACK: // arrLiteral
		return parseArrLiteralExp(p)
	case token.TOKEN_KW_TRUE:
		return parseTrueExp(p)
	case token.TOKEN_KW_FALSE:
		return parseFalseExp(p)
	case token.TOKEN_KW_NIL:
		return parseNilExp(p)
	case token.TOKEN_KW_FUNC:
		return parseFuncLiteralExp(p)
	case token.TOKEN_KW_NEW:
		return parseNewObjectExp(p)
	case token.TOKEN_IDENTIFIER:
		return parseFuncCallOrAttrExp(p)
	case token.TOKEN_SEP_LPAREN: // (exp)
		p.l.NextToken()
		exp := parseExp(p)
		p.NextTokenKind(token.TOKEN_SEP_RPAREN)
		return exp
	default:
		p.exit("unexpected token '%s' for parsing factor", ahead.Content)
	}
	return nil
}

func parseFuncLiteralExp(p *Parser) *ast.FuncLiteralExp {
	p.l.NextToken()
	return &ast.FuncLiteralExp{FuncLiteral: p.parseFuncLiteral()}
}

func parseNilExp(p *Parser) *ast.NilExp {
	p.l.NextToken()
	return &ast.NilExp{}
}

func parseFalseExp(p *Parser) *ast.FalseExp {
	p.l.NextToken()
	return &ast.FalseExp{}
}

func parseTrueExp(p *Parser) *ast.TrueExp {
	p.l.NextToken()
	return &ast.TrueExp{}
}

func parseNumberLiteralExp(p *Parser) *ast.NumberLiteralExp {
	return &ast.NumberLiteralExp{Value: p.l.NextToken().Value}
}

func parseStringLiteralExp(p *Parser) *ast.StringLiteralExp {
	return &ast.StringLiteralExp{Value: p.l.NextToken().Value.(string)}
}

func parseFuncCallOrAttrExp(p *Parser) ast.Exp {
	var exp ast.Exp
	exp = &ast.NameExp{
		Line: p.l.Line(),
		Name: p.l.NextToken().Content,
	}
	for {
		switch p.l.LookAhead().Kind {
		case token.TOKEN_SEP_DOT: // access attribute
			p.l.NextToken()
			exp = &ast.BinOpExp{
				Exp1: exp,
				Exp2: &ast.StringLiteralExp{
					Value: p.NextTokenKind(token.TOKEN_IDENTIFIER).Content,
				},
				BinOp: ast.BINOP_ATTR,
			}
		case token.TOKEN_SEP_LBRACK: // access map
			p.l.NextToken()
			exp = &ast.BinOpExp{
				Exp1:  exp,
				Exp2:  parseExp(p),
				BinOp: ast.BINOP_ATTR,
			}
			p.NextTokenKind(token.TOKEN_SEP_RBRACK)
		case token.TOKEN_SEP_LPAREN: // function call
			exp = &ast.FuncCallExp{Func: exp, Args: parseExpListBlock(p)}
		default:
			return exp
		}
	}
}

func parseNewObjectExp(p *Parser) *ast.NewObjectExp {
	p.l.NextToken()
	exp := &ast.NewObjectExp{}
	exp.Name = p.NextTokenKind(token.TOKEN_IDENTIFIER).Content
	exp.Line = p.l.Line()
	if !p.Expect(token.TOKEN_SEP_LPAREN) {
		return exp
	}
	exp.Args = parseExpListBlock(p)
	return exp
}

func parseArrLiteralExp(p *Parser) *ast.ArrLiteralExp {
	p.l.NextToken()
	exp := &ast.ArrLiteralExp{}
	for !p.Expect(token.TOKEN_SEP_RBRACK) {
		exp.Vals = append(exp.Vals, parseExp(p))
		if p.Expect(token.TOKEN_SEP_COMMA) || p.Expect(token.TOKEN_SEP_SEMI) {
			p.l.NextToken()
		}
	}
	p.NextTokenKind(token.TOKEN_SEP_RBRACK)
	return exp
}

func parseMapLiteralExp(p *Parser) *ast.MapLiteralExp {
	p.l.NextToken()
	exp := &ast.MapLiteralExp{}
loop:
	for {
		var key interface{}
		switch p.l.LookAhead().Kind {
		case token.TOKEN_KW_TRUE:
			key = true
		case token.TOKEN_KW_FALSE:
			key = false
		case token.TOKEN_STRING, token.TOKEN_NUMBER:
			key = p.l.LookAhead().Value
		case token.TOKEN_IDENTIFIER:
			key = p.l.LookAhead().Content
		default:
			break loop
		}
		p.l.NextToken()
		exp.Keys = append(exp.Keys, key)

		var val ast.Exp
		p.NextTokenKind(token.TOKEN_SEP_COLON)
		val = parseExp(p)
		exp.Vals = append(exp.Vals, val)
		if p.Expect(token.TOKEN_SEP_COMMA) || p.Expect(token.TOKEN_SEP_SEMI) {
			p.l.NextToken()
		} else if !p.Expect(token.TOKEN_SEP_RCURLY) {
			p.exit("expect '}', but got %s", p.l.NextToken().Content)
		}
	}
	p.NextTokenKind(token.TOKEN_SEP_RCURLY)
	return exp
}

func _parseBinExp(p *Parser, expect []int, cb func(*Parser) ast.Exp) ast.Exp {
	exp := cb(p)
	for {
		flag := func() bool {
			for _, kind := range expect {
				if p.Expect(kind) {
					return true
				}
			}
			return false
		}()
		if !flag {
			return exp
		}
		binExp := &ast.BinOpExp{
			Exp1:  exp,
			BinOp: p.l.NextToken().Kind,
			Exp2:  cb(p),
		}
		exp = binExp
	}
}
