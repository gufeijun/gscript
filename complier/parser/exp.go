package parser

import (
	"gscript/complier/ast"
	. "gscript/complier/lexer"
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
func parseExpListBlock(l *Lexer) []ast.Exp {
	l.NextTokenKind(TOKEN_SEP_LPAREN)
	if l.Expect(TOKEN_SEP_RPAREN) {
		l.NextToken()
		return nil
	}
	exps := parseExpList(l)
	l.NextTokenKind(TOKEN_SEP_RPAREN)
	return exps
}

func parseExpList(l *Lexer) (exps []ast.Exp) {
	exps = append(exps, parseExp(l))
	for l.Expect(TOKEN_SEP_COMMA) {
		l.NextToken()
		exps = append(exps, parseExp(l))
	}
	return
}

func parseExp(l *Lexer) ast.Exp {
	return parseTerm12(l)
}

func parseTerm12(l *Lexer) ast.Exp {
	exp1 := parseTerm11(l)
	if !l.Expect(TOKEN_SEP_QMARK) { // ?
		return exp1
	}
	l.NextToken()
	exp2 := parseTerm11(l)
	l.NextTokenKind(TOKEN_SEP_COLON) // :
	return &ast.TernaryOpExp{exp1, exp2, parseTerm11(l)}
}

func parseTerm11(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_LOR}, parseTerm10)
}

func parseTerm10(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_LAND}, parseTerm9)
}

func parseTerm9(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_OR}, parseTerm8)
}

func parseTerm8(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_XOR}, parseTerm7)
}

func parseTerm7(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_AND}, parseTerm6)
}

func parseTerm6(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_EQ, TOKEN_OP_NE}, parseTerm5)
}

func parseTerm5(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_LE, TOKEN_OP_GE, TOKEN_OP_LT, TOKEN_OP_GT}, parseTerm4)
}

func parseTerm4(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_SHL, TOKEN_OP_SHR}, parseTerm3)
}

func parseTerm3(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_ADD, TOKEN_OP_SUB}, parseTerm2)
}

func parseTerm2(l *Lexer) ast.Exp {
	return _parseBinExp(l, []int{TOKEN_OP_DIV, TOKEN_OP_MUL, TOKEN_OP_MOD, TOKEN_OP_IDIV}, parseTerm1)
}

func parseTerm1(l *Lexer) ast.Exp {
	kind := l.LookAhead().Kind

	var unOp int
	if kind == TOKEN_OP_SUB {
		unOp = ast.UNOP_NEG
	} else if kind == TOKEN_OP_NOT {
		unOp = ast.UNOP_NOT
	} else if kind == TOKEN_OP_LNOT {
		unOp = ast.UNOP_LNOT
	} else {
		return parseTerm0(l)
	}
	l.NextToken()
	return &ast.UnOpExp{Op: unOp, Exp: parseTerm1(l)}
}

func parseTerm0(l *Lexer) ast.Exp {
	if l.Expect(TOKEN_OP_INC) {
		l.NextToken()
		return &ast.UnOpExp{ast.UNOP_INC, parseFactor(l)}
	} else if l.Expect(TOKEN_OP_DEC) {
		l.NextToken()
		return &ast.UnOpExp{ast.UNOP_DEC, parseFactor(l)}
	}
	exp := parseFactor(l)
	if l.Expect(TOKEN_OP_INC) {
		l.NextToken()
		return &ast.UnOpExp{ast.UNOP_INC_, exp}
	} else if l.Expect(TOKEN_OP_DEC) {
		l.NextToken()
		return &ast.UnOpExp{ast.UNOP_DEC_, exp}
	}
	return exp
}

// factor ::= '(' exp ')' | literal | nil | new ID ['(' [expList] ')']
// 		 | ID { ( '.' ID | '[' exp ']' | '(' [expList] ')' ) }
func parseFactor(l *Lexer) ast.Exp {
	switch l.LookAhead().Kind {
	case TOKEN_STRING:
		return parseStringLiteralExp(l)
	case TOKEN_NUMBER:
		return parseNumberLiteralExp(l)
	case TOKEN_SEP_LCURLY: // mapLiteral
		return parseMapLiteralExp(l)
	case TOKEN_SEP_LBRACK: // arrLiteral
		return parseArrLiteralExp(l)
	case TOKEN_KW_TRUE:
		return parseTrueExp(l)
	case TOKEN_KW_FALSE:
		return parseFalseExp(l)
	case TOKEN_KW_NIL:
		return parseNilExp(l)
	case TOKEN_KW_FUNC:
		return parseFuncLiteralExp(l)
	case TOKEN_KW_NEW:
		return parseNewObjectExp(l)
	case TOKEN_IDENTIFIER:
		return parseFuncCallOrAttrExp(l)
	case TOKEN_SEP_LPAREN: // (exp)
		l.NextToken()
		exp := parseExp(l)
		l.NextTokenKind(TOKEN_SEP_RPAREN)
		return exp
	default:
		panic(l.Line())
	}
}

func parseFuncLiteralExp(l *Lexer) *ast.FuncLiteralExp {
	l.NextToken()
	return &ast.FuncLiteralExp{parseFuncLiteral(l)}
}

func parseNilExp(l *Lexer) *ast.NilExp {
	l.NextToken()
	return &ast.NilExp{}
}

func parseFalseExp(l *Lexer) *ast.FalseExp {
	l.NextToken()
	return &ast.FalseExp{}
}

func parseTrueExp(l *Lexer) *ast.TrueExp {
	l.NextToken()
	return &ast.TrueExp{}
}

func parseNumberLiteralExp(l *Lexer) *ast.NumberLiteralExp {
	return &ast.NumberLiteralExp{l.NextToken().Value}
}

func parseStringLiteralExp(l *Lexer) *ast.StringLiteralExp {
	return &ast.StringLiteralExp{l.NextToken().Value.(string)}
}

func parseFuncCallOrAttrExp(l *Lexer) ast.Exp {
	var exp ast.Exp
	exp = &ast.NameExp{l.NextToken().Content}
	for {
		switch l.LookAhead().Kind {
		case TOKEN_SEP_DOT: // access attribute
			l.NextToken()
			exp = &ast.BinOpExp{
				Exp1:  exp,
				Exp2:  &ast.StringLiteralExp{l.NextTokenKind(TOKEN_IDENTIFIER).Content},
				BinOp: ast.BINOP_ATTR,
			}
		case TOKEN_SEP_LBRACK: // access map
			l.NextToken()
			exp = &ast.BinOpExp{
				Exp1:  exp,
				Exp2:  parseExp(l),
				BinOp: ast.BINOP_ATTR,
			}
			l.NextTokenKind(TOKEN_SEP_RBRACK)
		case TOKEN_SEP_LPAREN: // function call
			exp = &ast.FuncCallExp{exp, parseExpListBlock(l)}
		default:
			return exp
		}
	}
}

func parseNewObjectExp(l *Lexer) *ast.NewObjectExp {
	l.NextToken()
	exp := &ast.NewObjectExp{}
	exp.Name = l.NextTokenKind(TOKEN_IDENTIFIER).Content
	if !l.Expect(TOKEN_SEP_LPAREN) {
		return exp
	}
	exp.Args = parseExpListBlock(l)
	return exp
}

func parseArrLiteralExp(l *Lexer) *ast.ArrLiteralExp {
	l.NextToken()
	exp := &ast.ArrLiteralExp{}
	for !l.Expect(TOKEN_SEP_RBRACK) {
		exp.Vals = append(exp.Vals, parseExp(l))
		if l.Expect(TOKEN_SEP_COMMA) || l.Expect(TOKEN_SEP_SEMI) {
			l.NextToken()
		}
	}
	l.NextTokenKind(TOKEN_SEP_RBRACK)
	return exp
}

func parseMapLiteralExp(l *Lexer) *ast.MapLiteralExp {
	l.NextToken()
	exp := &ast.MapLiteralExp{}
loop:
	for {
		var key interface{}
		switch l.LookAhead().Kind {
		case TOKEN_KW_TRUE:
			key = true
		case TOKEN_KW_FALSE:
			key = false
		case TOKEN_STRING, TOKEN_NUMBER:
			key = l.LookAhead().Value
		case TOKEN_IDENTIFIER:
			key = l.LookAhead().Content
		default:
			break loop
		}
		l.NextToken()
		exp.Keys = append(exp.Keys, key)

		var val ast.Exp
		if l.Expect(TOKEN_SEP_COLON) {
			l.NextToken()
			val = parseExp(l)
		} else {
			val = &ast.NilExp{}
		}
		exp.Vals = append(exp.Vals, val)
		if l.Expect(TOKEN_SEP_COMMA) || l.Expect(TOKEN_SEP_SEMI) {
			l.NextToken()
		} else if !l.Expect(TOKEN_SEP_RCURLY) {
			panic(l.Line())
		}
	}
	l.NextTokenKind(TOKEN_SEP_RCURLY)
	return exp
}

func _parseBinExp(l *Lexer, expect []int, cb func(*Lexer) ast.Exp) ast.Exp {
	exp := cb(l)
	for {
		flag := func() bool {
			for _, kind := range expect {
				if l.Expect(kind) {
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
			BinOp: l.NextToken().Kind,
			Exp2:  cb(l),
		}
		exp = binExp
	}
}
