package parser

import (
	"gscript/complier/ast"
	. "gscript/complier/lexer"
)

func parseStmt(l *Lexer) ast.Stmt {
	var stmt ast.Stmt
	switch l.LookAhead().Kind {
	case TOKEN_KW_CONST, TOKEN_KW_LET:
		stmt = parseVarDeclStmt(l) // tested
	case TOKEN_IDENTIFIER:
		stmt = parseVarOpOrLabel(l) // tested
	case TOKEN_KW_FUNC:
		stmt = parseFunc(l) // tested
	case TOKEN_KW_BREAK:
		stmt = parseBreakStmt(l) // tested
	case TOKEN_KW_CONTINUE:
		stmt = parseContinueStmt(l) // tested
	case TOKEN_KW_RETURN:
		stmt = parseReturnStmt(l) // tested
	case TOKEN_KW_GOTO:
		stmt = parseGotoStmt(l) // tested
	case TOKEN_KW_FALLTHROUGH:
		stmt = parseFallthroughStmt(l) // tested
	case TOKEN_KW_LOOP:
		stmt = parseLoopStmt(l) // tested
	case TOKEN_KW_WHILE:
		stmt = parseWhileStmt(l) // tested
	case TOKEN_KW_FOR:
		stmt = parseForStmt(l) // tested
	case TOKEN_KW_IF:
		stmt = parseIfStmt(l) // tested
	case TOKEN_KW_CLASS:
		stmt = parseClassStmt(l)
	case TOKEN_KW_ENUM:
		stmt = parseEnumStmt(l) // tested
	case TOKEN_KW_SWITCH:
		stmt = parseSwitchStmt(l) // tested
	case TOKEN_OP_INC, TOKEN_OP_DEC:
		stmt = parseIncOrDecVar(l) // tested
	default:
		panic(l.Line())
	}
	l.ConsumeIf(TOKEN_SEP_SEMI)
	return stmt
}

// varDeclStmt ::= (const|let) id {,id} (:=|=) exp {,exp} ;
func parseVarDeclStmt(l *Lexer) (stmt *ast.VarDeclStmt) {
	stmt = new(ast.VarDeclStmt)
	stmt.Const = l.NextToken().Kind == TOKEN_KW_CONST // const or let
	stmt.Lefts = parseNameList(l)                     // id {,id}
	if !l.Expect(TOKEN_OP_CLONE) && !l.Expect(TOKEN_OP_ASSIGN) {
		stmt.Rights = make([]ast.Exp, len(stmt.Lefts))
		for i := range stmt.Rights {
			stmt.Rights[i] = &ast.NilExp{}
		}
		l.NextTokenKind(TOKEN_SEP_SEMI) // ;
		return
	}
	stmt.DeepEq = l.NextToken().Kind == TOKEN_OP_CLONE // = or :=
	stmt.Rights = parseExpList(l)                      // exp {,exp}
	return
}

// stmt ::= | varAssign	;						# case1: var {, var} assignOp exp {, exp}
//          | varIncOrDec ;						# case2: var++ or var--
//          | var '(' {explist} ')' callTail ;	# case3: function call
//          | ID ':'							# case4: label
func parseVarOpOrLabel(l *Lexer) ast.Stmt {
	name := l.NextToken().Content
	if l.ConsumeIf(TOKEN_SEP_COLON) { // case4
		return &ast.LabelStmt{name}
	}
	v := _parseVar(l, name)
	switch l.LookAhead().Kind {
	case TOKEN_OP_INC, TOKEN_OP_DEC: // case2
		stmt := _incOrDec2Assign(l.NextToken().Kind, v)
		return stmt
	case TOKEN_SEP_LPAREN: // case3
		stmt := new(ast.NamedFuncCallStmt)
		stmt.Var = v
		stmt.Args = parseExpListBlock(l)
		stmt.CallTails = parseCallTails(l)
		return stmt
	default: // case 1
		return _parseAssignStmt(l, v)
	}
}

// has parsed first var
func _parseAssignStmt(l *Lexer, v ast.Var) *ast.VarAssignStmt {
	var stmt ast.VarAssignStmt
	stmt.Lefts = append(stmt.Lefts, v)
	for l.ConsumeIf(TOKEN_SEP_COMMA) {
		stmt.Lefts = append(stmt.Lefts, parseVar(l))
	}
	token := l.NextToken()
	if !_isAsignOp(token.Kind) {
		panic(l.Line())
	}
	stmt.AssignOp = token.Kind
	stmt.Rights = parseExpList(l)
	return &stmt
}

// i++ i-- ==> i+=1 i-=1
func _incOrDec2Assign(kind int, v ast.Var) (stmt *ast.VarAssignStmt) {
	op := ast.ASIGN_OP_SUBEQ
	if kind == TOKEN_OP_INC {
		op = ast.ASIGN_OP_ADDEQ
	}
	return &ast.VarAssignStmt{
		AssignOp: op,
		Lefts:    []ast.Var{v},
		Rights:   []ast.Exp{&ast.NumberLiteralExp{int64(1)}},
	}
}

// ++i ==> i+=1    --i ==> i-=1
func parseIncOrDecVar(l *Lexer) (stmt *ast.VarAssignStmt) {
	return _incOrDec2Assign(l.NextToken().Kind, parseVar(l))
}

func _parseVar(l *Lexer, prefix string) (v ast.Var) {
	v.Prefix = prefix
	v.Attrs = parseAttrTail(l)
	return
}

func parseVar(l *Lexer) (v ast.Var) {
	if !l.Expect(TOKEN_IDENTIFIER) {
		panic(l.Line())
	}
	return _parseVar(l, l.NextToken().Content)
}

func parseAttrTail(l *Lexer) (exps []ast.Exp) {
	for {
		switch l.LookAhead().Kind {
		case TOKEN_SEP_DOT:
			l.NextToken()
			// a.b ==> a["b"]
			token := l.NextTokenKind(TOKEN_IDENTIFIER)
			exps = append(exps, &ast.StringLiteralExp{token.Content})
		case TOKEN_SEP_LBRACK: //[
			l.NextToken()
			exps = append(exps, parseExp(l))
			l.NextTokenKind(TOKEN_SEP_RBRACK)
		default:
			return
		}
	}
}

func parseSwitchStmt(l *Lexer) (stmt *ast.SwitchStmt) {
	stmt = new(ast.SwitchStmt)
	l.NextToken()                 // switch
	stmt.Value = parseExpBlock(l) // (condition)
	l.ConsumeIf(TOKEN_SEP_SEMI)
	l.NextTokenKind(TOKEN_SEP_LCURLY) // {

	for l.LookAhead().Kind == TOKEN_KW_CASE {
		l.NextToken() // case
		stmt.Cases = append(stmt.Cases, parseConstLiterals(l))
		l.NextTokenKind(TOKEN_SEP_COLON) // :
		stmt.Blocks = append(stmt.Blocks, parseBlockStmts(l))
	}

	if l.Expect(TOKEN_KW_DEFAULT) { // default
		l.NextToken()
		l.NextTokenKind(TOKEN_SEP_COLON)
		stmt.Default = parseBlockStmts(l)
	}
	l.NextTokenKind(TOKEN_SEP_RCURLY) // }
	return
}

func isConstLiteral(token *Token) bool {
	return token.Kind == TOKEN_KW_TRUE || token.Kind == TOKEN_KW_FALSE ||
		token.Kind == TOKEN_STRING || token.Kind == TOKEN_NUMBER
}

func parseConstLiterals(l *Lexer) (literals []interface{}) {
	literals = append(literals, parseConstLiteral(l))
	for l.LookAhead().Kind == TOKEN_SEP_COMMA {
		l.NextToken()
		literals = append(literals, parseConstLiteral(l))
	}
	return
}

func parseConstLiteral(l *Lexer) interface{} {
	token := l.NextToken()
	if !isConstLiteral(token) {
		panic(l.Line())
	}
	return token.Value
}

func parseEnumStmt(l *Lexer) (stmt ast.EnumStmt) {
	l.NextToken()
	l.NextTokenKind(TOKEN_SEP_LCURLY)
	var enum int64
	var ok bool
	for {
		if l.ConsumeIf(TOKEN_SEP_RCURLY) {
			return
		}
		if !l.Expect(TOKEN_IDENTIFIER) {
			panic(l.Line())
		}
		stmt.Names = append(stmt.Names, l.NextToken().Content)
		if l.ConsumeIf(TOKEN_OP_ASSIGN) {
			if enum, ok = l.NextToken().Value.(int64); !ok {
				panic(l.Line())
			}
		}
		stmt.Values = append(stmt.Values, enum)
		enum++
		if l.ConsumeIf(TOKEN_SEP_RCURLY) {
			return
		}
		if l.Expect(TOKEN_SEP_COMMA) || l.Expect(TOKEN_SEP_SEMI) {
			l.NextToken()
		}
	}
}

func parseClassStmt(l *Lexer) (stmt *ast.ClassStmt) {
	stmt = new(ast.ClassStmt)
	l.NextToken()
	stmt.Name = l.NextTokenKind(TOKEN_IDENTIFIER).Content
	l.ConsumeIf(TOKEN_SEP_SEMI)
	l.NextTokenKind(TOKEN_SEP_LCURLY)

	for l.Expect(TOKEN_IDENTIFIER) {
		attr := l.NextToken().Content
		if l.Expect(TOKEN_SEP_LPAREN) {
			stmt.AttrName = append(stmt.AttrName, attr)
			stmt.AttrValue = append(stmt.AttrValue, &ast.FuncLiteralExp{parseFuncLiteral(l)})
		} else if l.ConsumeIf(TOKEN_OP_ASSIGN) {
			value := parseExp(l)
			stmt.AttrName = append(stmt.AttrName, attr)
			stmt.AttrValue = append(stmt.AttrValue, value)
		}
		if l.Expect(TOKEN_SEP_SEMI) || l.Expect(TOKEN_SEP_COMMA) {
			l.NextToken()
		}
	}

	l.NextTokenKind(TOKEN_SEP_RCURLY)
	return
}

func parseClassBody(l *Lexer) (attr string, value ast.Exp) {
	attr = l.NextToken().Content
	if l.ConsumeIf(TOKEN_OP_ASSIGN) {
		value = parseExp(l)
	}
	return
}

// single stmt to block
func parseBlockStmt(l *Lexer) (stmt ast.Block) {
	switch l.LookAhead().Kind {
	case TOKEN_SEP_LCURLY:
		return parseBlock(l)
	default:
		stmt.Blocks = append(stmt.Blocks, parseStmt(l))
		return
	}
}

func parseIfStmt(l *Lexer) (stmt ast.IfStmt) {
	l.NextToken() // if
	stmt.Conditions = append(stmt.Conditions, parseExpBlock(l))
	l.ConsumeIf(TOKEN_SEP_SEMI)
	stmt.Blocks = append(stmt.Blocks, parseBlockStmt(l)) // block | stmt

	for l.ConsumeIf(TOKEN_KW_ELIF) { // elif
		stmt.Conditions = append(stmt.Conditions, parseExpBlock(l))
		l.ConsumeIf(TOKEN_SEP_SEMI)
		stmt.Blocks = append(stmt.Blocks, parseBlockStmt(l))
	}

	// else ==> elif(true)
	if l.ConsumeIf(TOKEN_KW_ELSE) { // else
		stmt.Conditions = append(stmt.Conditions, &ast.TrueExp{})
		stmt.Blocks = append(stmt.Blocks, parseBlockStmt(l))
	}

	return
}

// (exp)
func parseExpBlock(l *Lexer) (exp ast.Exp) {
	l.NextTokenKind(TOKEN_SEP_LPAREN) // (
	exp = parseExp(l)                 // condition
	l.NextTokenKind(TOKEN_SEP_RPAREN) // )
	return
}

func _isAsignOp(kind int) bool {
	return kind > TOKEN_ASIGN_START && kind < TOKEN_ASIGN_END
}

func parseVarAssignStmt(l *Lexer, prefix string) (stmt *ast.VarAssignStmt) {
	return _parseAssignStmt(l, _parseVar(l, prefix))
}

func parseForStmt(l *Lexer) (stmt *ast.ForStmt) {
	l.NextToken()
	l.NextTokenKind(TOKEN_SEP_LPAREN) // (

	stmt = new(ast.ForStmt)
	switch l.LookAhead().Kind {
	case TOKEN_KW_CONST, TOKEN_KW_LET:
		stmt.DeclStmt = parseVarDeclStmt(l)
	case TOKEN_IDENTIFIER:
		stmt.AsgnStmt = parseVarAssignStmt(l, l.NextToken().Content)
	default:
		panic(l.Line())
	}

	l.NextTokenKind(TOKEN_SEP_SEMI) // ;
	stmt.Condition = parseExp(l)
	l.NextTokenKind(TOKEN_SEP_SEMI) // ;
	stmt.ForTail = parseForTail(l)
	l.NextTokenKind(TOKEN_SEP_RPAREN) // )
	l.ConsumeIf(TOKEN_SEP_SEMI)
	stmt.Block = parseBlockStmt(l)
	return
}

// forTail ::= varAssign | varIncOrDec | incOrDecVar
func parseForTail(l *Lexer) *ast.VarAssignStmt {
	if l.Expect(TOKEN_OP_INC) || l.Expect(TOKEN_OP_DEC) { // incOrDecVar
		return parseIncOrDecVar(l)
	}
	v := parseVar(l)
	if l.Expect(TOKEN_OP_INC) || l.Expect(TOKEN_OP_DEC) { // varIncOrDec
		return _incOrDec2Assign(l.NextToken().Kind, v)
	}
	return _parseAssignStmt(l, v) // varAssign
}

func parseLoopStmt(l *Lexer) (stmt *ast.LoopStmt) {
	l.NextToken()
	l.NextTokenKind(TOKEN_SEP_LPAREN)
	l.NextTokenKind(TOKEN_KW_LET)

	stmt = new(ast.LoopStmt)
	firstName := l.NextTokenKind(TOKEN_IDENTIFIER).Content
	if l.Expect(TOKEN_SEP_COMMA) {
		stmt.Key = firstName
		l.NextToken()
		stmt.Val = l.NextTokenKind(TOKEN_IDENTIFIER).Content
	} else {
		stmt.Val = firstName
	}
	l.NextTokenKind(TOKEN_SEP_COLON) // :

	stmt.Iterator = parseExp(l)

	l.NextTokenKind(TOKEN_SEP_RPAREN) // )
	l.ConsumeIf(TOKEN_SEP_SEMI)

	stmt.Block = parseBlockStmt(l)

	return
}

func parseWhileStmt(l *Lexer) (stmt *ast.WhileStmt) {
	l.NextToken()
	stmt = new(ast.WhileStmt)
	stmt.Condition = parseExpBlock(l)
	l.ConsumeIf(TOKEN_SEP_SEMI)
	stmt.Block = parseBlockStmt(l)
	return
}

func parseFallthroughStmt(l *Lexer) (stmt ast.FallthroughStmt) {
	l.NextToken()
	return
}

func parseGotoStmt(l *Lexer) (stmt ast.GotoStmt) {
	l.NextToken()
	if !l.Expect(TOKEN_IDENTIFIER) {
		panic("")
	}
	stmt.Label = l.NextToken().Content
	return
}

func parseNameList(l *Lexer) (names []string) {
	names = append(names, parseName(l))
	for l.Expect(TOKEN_SEP_COMMA) {
		l.NextToken()
		names = append(names, parseName(l))
	}
	return names
}

func parseName(l *Lexer) string {
	if !l.Expect(TOKEN_IDENTIFIER) {
		panic(l.Line())
	}
	return l.NextToken().Content
}

// stmt ::= func ID funcBody ';          			   	// function definition
//	      | funcLiteral '(' [expList] ')' callTail ';'  // anonymous function call
func parseFunc(l *Lexer) ast.Stmt {
	l.NextToken() // func
	switch l.LookAhead().Kind {
	case TOKEN_IDENTIFIER:
		return parseFuncDefStmt(l)
	case TOKEN_SEP_LPAREN:
		return parseAnonymousFuncCallStmt(l)
	default:
		panic(l.Line())
	}
}

func parseFuncDefStmt(l *Lexer) (stmt *ast.FuncDefStmt) {
	stmt = new(ast.FuncDefStmt)
	stmt.Name = l.NextToken().Content
	stmt.FuncLiteral = parseFuncLiteral(l)
	return
}

func parseFuncLiteral(l *Lexer) (literal ast.FuncLiteral) {
	var defaultValue bool
	l.NextTokenKind(TOKEN_SEP_LPAREN) // (

	literal.Parameters, defaultValue = parseParameters(l)

	if l.ConsumeIf(TOKEN_SEP_VARARG) { // ...
		if defaultValue {
			// TODO error
			panic(l.Line())
		}
		literal.VarArg = l.NextTokenKind(TOKEN_IDENTIFIER).Content
	}
	l.NextTokenKind(TOKEN_SEP_RPAREN) // )

	literal.Block = parseBlock(l)
	return
}

func parseParameters(l *Lexer) (pars []ast.Parameter, defaultValue bool) {
	for l.Expect(TOKEN_IDENTIFIER) {
		par := ast.Parameter{}
		par.Name = l.NextToken().Content
		if l.Expect(TOKEN_OP_ASSIGN) {
			defaultValue = true
			l.NextToken()
			par.Default = parseExp(l)
		} else if defaultValue {
			// TODO error
			panic(l.Line())
		}
		pars = append(pars, par)
		if l.Expect(TOKEN_SEP_COMMA) {
			l.NextToken()
		}
	}
	return
}

func parseAnonymousFuncCallStmt(l *Lexer) (stmt *ast.AnonymousFuncCallStmt) {
	stmt = new(ast.AnonymousFuncCallStmt)
	stmt.FuncLiteral = parseFuncLiteral(l)
	stmt.CallArgs = parseExpListBlock(l)
	stmt.CallTails = parseCallTails(l)
	return
}

// . or [ or (
func expectTail(token *Token) bool {
	return token.Kind == TOKEN_SEP_DOT || token.Kind == TOKEN_SEP_LBRACK || token.Kind == TOKEN_SEP_LPAREN
}

func parseCallTails(l *Lexer) (tails []ast.CallTail) {
	for expectTail(l.LookAhead()) {
		var tail ast.CallTail
		tail.Attrs = parseAttrTail(l)
		tail.Args = parseExpListBlock(l)
		tails = append(tails, tail)
	}
	return
}

func parseBreakStmt(l *Lexer) (stmt ast.BreakStmt) {
	l.NextToken()
	return
}

func parseContinueStmt(l *Lexer) (stmt ast.ContinueStmt) {
	l.NextToken()
	return
}

func parseReturnStmt(l *Lexer) (stmt *ast.ReturnStmt) {
	stmt = new(ast.ReturnStmt)
	l.NextToken()
	if l.Expect(TOKEN_SEP_RCURLY) || l.Expect(TOKEN_SEP_SEMI) {
		return
	}
	stmt.Args = parseExpList(l)
	return
}
