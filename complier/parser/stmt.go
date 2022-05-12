package parser

import (
	"gscript/complier/ast"
	. "gscript/complier/lexer"
)

func (p *Parser) parseStmt() ast.Stmt {
	var stmt ast.Stmt
	switch p.l.LookAhead().Kind {
	case TOKEN_KW_CONST, TOKEN_KW_LET:
		stmt = p.parseVarDeclStmt()
	case TOKEN_IDENTIFIER:
		stmt = p.parseVarOpOrLabel()
	case TOKEN_KW_FUNC:
		stmt = p.parseFunc()
	case TOKEN_KW_BREAK:
		stmt = p.parseBreakStmt()
	case TOKEN_KW_CONTINUE:
		stmt = p.parseContinueStmt()
	case TOKEN_KW_RETURN:
		stmt = p.parseReturnStmt()
	case TOKEN_KW_GOTO:
		stmt = p.parseGotoStmt()
	case TOKEN_KW_FALLTHROUGH:
		stmt = p.parseFallthroughStmt()
	case TOKEN_KW_LOOP:
		stmt = p.parseLoopStmt()
	case TOKEN_KW_WHILE:
		stmt = p.parseWhileStmt()
	case TOKEN_KW_FOR:
		stmt = p.parseForStmt()
	case TOKEN_KW_IF:
		stmt = p.parseIfStmt()
	case TOKEN_KW_CLASS:
		stmt = p.parseClassStmt()
	case TOKEN_KW_ENUM:
		stmt = p.parseEnumStmt()
	case TOKEN_KW_SWITCH:
		stmt = p.parseSwitchStmt()
	case TOKEN_OP_INC, TOKEN_OP_DEC:
		stmt = p.parseIncOrDecVar()
	default:
		panic(p.l.Line())
	}
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	return stmt
}

// varDeclStmt ::= (const|let) id {,id} = exp {,exp} ;
func (p *Parser) parseVarDeclStmt() (stmt *ast.VarDeclStmt) {
	stmt = new(ast.VarDeclStmt)
	stmt.Const = p.l.NextToken().Kind == TOKEN_KW_CONST // const or let
	stmt.Lefts = p.parseNameList()                      // id {,id}
	if !p.l.Expect(TOKEN_OP_ASSIGN) {
		if stmt.Const {
			panic("const variable should have initial value") // TODO
		}
		stmt.Rights = make([]ast.Exp, len(stmt.Lefts))
		for i := range stmt.Rights {
			stmt.Rights[i] = &ast.NilExp{}
		}
		return
	}
	p.l.NextToken()
	stmt.Rights = parseExpList(p) // exp {,exp}
	if len(stmt.Lefts) < len(stmt.Rights) {
		panic("") // TODO
	}
	return
}

// stmt ::= | varAssign	;						# case1: var {, var} assignOp exp {, exp}
//          | varIncOrDec ;						# case2: var++ or var--
//          | var '(' {explist} ')' callTail ;	# case3: function call
//          | ID ':'							# case4: label
func (p *Parser) parseVarOpOrLabel() ast.Stmt {
	name := p.l.NextToken().Content
	if p.l.ConsumeIf(TOKEN_SEP_COLON) { // case4
		label := &ast.LabelStmt{Name: name}
		return label
	}
	v := p._parseVar(name)
	switch p.l.LookAhead().Kind {
	case TOKEN_OP_INC, TOKEN_OP_DEC: // case2
		stmt := _incOrDec2Assign(p.l.NextToken().Kind, v)
		return stmt
	case TOKEN_SEP_LPAREN: // case3
		stmt := new(ast.NamedFuncCallStmt)
		stmt.Prefix = v.Prefix
		stmt.CallTails = append(stmt.CallTails, ast.CallTail{
			Attrs: v.Attrs,
			Args:  parseExpListBlock(p),
		})
		stmt.CallTails = append(stmt.CallTails, p.parseCallTails()...)
		return stmt
	default: // case 1
		return p._parseAssignStmt(v)
	}
}

// has parsed first var
func (p *Parser) _parseAssignStmt(v ast.Var) *ast.VarAssignStmt {
	var stmt ast.VarAssignStmt
	stmt.Lefts = append(stmt.Lefts, v)
	for p.l.ConsumeIf(TOKEN_SEP_COMMA) {
		stmt.Lefts = append(stmt.Lefts, p.parseVar())
	}
	token := p.l.NextToken()
	if !_isAsignOp(token.Kind) {
		panic(p.l.Line())
	}
	stmt.AssignOp = token.Kind
	stmt.Rights = parseExpList(p)
	if len(stmt.Lefts) < len(stmt.Rights) {
		panic("") // TODO
	}
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
		Rights:   []ast.Exp{&ast.NumberLiteralExp{Value: int64(1)}},
	}
}

// ++i ==> i+=1    --i ==> i-=1
func (p *Parser) parseIncOrDecVar() (stmt *ast.VarAssignStmt) {
	return _incOrDec2Assign(p.l.NextToken().Kind, p.parseVar())
}

func (p *Parser) _parseVar(prefix string) (v ast.Var) {
	v.Prefix = prefix
	v.Attrs = p.parseAttrTail()
	return
}

func (p *Parser) parseVar() (v ast.Var) {
	if !p.l.Expect(TOKEN_IDENTIFIER) {
		panic(p.l.Line())
	}
	return p._parseVar(p.l.NextToken().Content)
}

func (p *Parser) parseAttrTail() (exps []ast.Exp) {
	for {
		switch p.l.LookAhead().Kind {
		case TOKEN_SEP_DOT:
			p.l.NextToken()
			// a.b ==> a["b"]
			token := p.l.NextTokenKind(TOKEN_IDENTIFIER)
			exps = append(exps, &ast.StringLiteralExp{Value: token.Content})
		case TOKEN_SEP_LBRACK: //[
			p.l.NextToken()
			exps = append(exps, parseExp(p))
			p.l.NextTokenKind(TOKEN_SEP_RBRACK)
		default:
			return
		}
	}
}

func (p *Parser) parseSwitchStmt() (stmt *ast.SwitchStmt) {
	stmt = new(ast.SwitchStmt)
	p.l.NextToken()                // switch
	stmt.Value = p.parseExpBlock() // (condition)
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	p.l.NextTokenKind(TOKEN_SEP_LCURLY) // {

	for p.l.ConsumeIf(TOKEN_KW_CASE) {
		stmt.Cases = append(stmt.Cases, parseExpList(p))
		p.l.NextTokenKind(TOKEN_SEP_COLON) // :
		stmt.Blocks = append(stmt.Blocks, p.parseBlockStmts())
	}

	if p.l.ConsumeIf(TOKEN_KW_DEFAULT) { // default
		p.l.NextTokenKind(TOKEN_SEP_COLON)
		stmt.Default = p.parseBlockStmts()
	}
	p.l.NextTokenKind(TOKEN_SEP_RCURLY) // }
	return
}

func (p *Parser) parseEnumStmt() (stmt *ast.EnumStmt) {
	stmt = new(ast.EnumStmt)
	p.l.NextToken()
	p.l.NextTokenKind(TOKEN_SEP_LCURLY)
	var enum int64
	var ok bool
	for {
		if p.l.ConsumeIf(TOKEN_SEP_RCURLY) {
			break
		}
		if !p.l.Expect(TOKEN_IDENTIFIER) {
			panic(p.l.Line())
		}
		stmt.Names = append(stmt.Names, p.l.NextToken().Content)
		if p.l.ConsumeIf(TOKEN_OP_ASSIGN) {
			if enum, ok = p.l.NextToken().Value.(int64); !ok {
				panic(p.l.Line())
			}
		}
		stmt.Values = append(stmt.Values, enum)
		enum++
		if p.l.ConsumeIf(TOKEN_SEP_RCURLY) {
			break
		}
		if p.l.Expect(TOKEN_SEP_COMMA) || p.l.Expect(TOKEN_SEP_SEMI) {
			p.l.NextToken()
		}
	}
	p.EnumStmts = append(p.EnumStmts, stmt)
	return
}

func (p *Parser) parseClassStmt() (stmt *ast.ClassStmt) {
	stmt = new(ast.ClassStmt)
	p.l.NextToken()
	stmt.Name = p.l.NextTokenKind(TOKEN_IDENTIFIER).Content
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	p.l.NextTokenKind(TOKEN_SEP_LCURLY)

	for p.l.Expect(TOKEN_IDENTIFIER) {
		attr := p.l.NextToken().Content
		if p.l.Expect(TOKEN_SEP_LPAREN) {
			stmt.AttrName = append(stmt.AttrName, attr)
			stmt.AttrValue = append(stmt.AttrValue, &ast.FuncLiteralExp{FuncLiteral: p.parseFuncLiteral()})
		} else if p.l.ConsumeIf(TOKEN_OP_ASSIGN) {
			value := parseExp(p)
			stmt.AttrName = append(stmt.AttrName, attr)
			stmt.AttrValue = append(stmt.AttrValue, value)
		}
		if p.l.Expect(TOKEN_SEP_SEMI) || p.l.Expect(TOKEN_SEP_COMMA) {
			p.l.NextToken()
		}
	}

	p.l.NextTokenKind(TOKEN_SEP_RCURLY)
	p.ClassStmts = append(p.ClassStmts, stmt)
	return
}

// single stmt to block
func (p *Parser) parseBlockStmt() (stmt ast.Block) {
	switch p.l.LookAhead().Kind {
	case TOKEN_SEP_LCURLY:
		return p.parseBlock()
	default:
		stmt.Blocks = append(stmt.Blocks, p.parseStmt())
		return
	}
}

func (p *Parser) parseIfStmt() (stmt *ast.IfStmt) {
	stmt = new(ast.IfStmt)
	p.l.NextToken() // if
	stmt.Conditions = append(stmt.Conditions, p.parseExpBlock())
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	stmt.Blocks = append(stmt.Blocks, p.parseBlockStmt()) // block | stmt

	for p.l.ConsumeIf(TOKEN_KW_ELIF) { // elif
		stmt.Conditions = append(stmt.Conditions, p.parseExpBlock())
		p.l.ConsumeIf(TOKEN_SEP_SEMI)
		stmt.Blocks = append(stmt.Blocks, p.parseBlockStmt())
	}

	// else ==> elif(true)
	if p.l.ConsumeIf(TOKEN_KW_ELSE) { // else
		stmt.Conditions = append(stmt.Conditions, &ast.TrueExp{})
		stmt.Blocks = append(stmt.Blocks, p.parseBlockStmt())
	}

	return
}

// (exp)
func (p *Parser) parseExpBlock() (exp ast.Exp) {
	p.l.NextTokenKind(TOKEN_SEP_LPAREN) // (
	exp = parseExp(p)                   // condition
	p.l.NextTokenKind(TOKEN_SEP_RPAREN) // )
	return
}

func _isAsignOp(kind int) bool {
	return kind > TOKEN_ASIGN_START && kind < TOKEN_ASIGN_END
}

func (p *Parser) parseVarAssignStmt(prefix string) (stmt *ast.VarAssignStmt) {
	return p._parseAssignStmt(p._parseVar(prefix))
}

func (p *Parser) parseForStmt() (stmt *ast.ForStmt) {
	stmt = new(ast.ForStmt)
	p.l.NextToken()
	p.l.NextTokenKind(TOKEN_SEP_LPAREN) // (

	switch p.l.LookAhead().Kind {
	case TOKEN_KW_CONST, TOKEN_KW_LET:
		stmt.DeclStmt = p.parseVarDeclStmt()
	case TOKEN_IDENTIFIER:
		stmt.AsgnStmt = p.parseVarAssignStmt(p.l.NextToken().Content)
	case TOKEN_SEP_SEMI:
		break
	default:
		panic(p.l.Line())
	}

	p.l.NextTokenKind(TOKEN_SEP_SEMI) // ;
	if !p.l.Expect(TOKEN_SEP_SEMI) {
		stmt.Condition = parseExp(p)
	}
	p.l.NextTokenKind(TOKEN_SEP_SEMI) // ;
	if !p.l.Expect(TOKEN_SEP_RPAREN) {
		stmt.ForTail = p.parseForTail()
	}
	p.l.NextTokenKind(TOKEN_SEP_RPAREN) // )
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	stmt.Block = p.parseBlockStmt()
	return
}

// forTail ::= varAssign | varIncOrDec | incOrDecVar
func (p *Parser) parseForTail() *ast.VarAssignStmt {
	if p.l.Expect(TOKEN_OP_INC) || p.l.Expect(TOKEN_OP_DEC) { // incOrDecVar
		return p.parseIncOrDecVar()
	}
	v := p.parseVar()
	if p.l.Expect(TOKEN_OP_INC) || p.l.Expect(TOKEN_OP_DEC) { // varIncOrDec
		return _incOrDec2Assign(p.l.NextToken().Kind, v)
	}
	return p._parseAssignStmt(v) // varAssign
}

func (p *Parser) parseLoopStmt() (stmt *ast.LoopStmt) {
	p.l.NextToken()
	p.l.NextTokenKind(TOKEN_SEP_LPAREN)
	p.l.NextTokenKind(TOKEN_KW_LET)

	stmt = new(ast.LoopStmt)
	firstName := p.l.NextTokenKind(TOKEN_IDENTIFIER).Content
	if p.l.ConsumeIf(TOKEN_SEP_COMMA) {
		stmt.Key = firstName
		stmt.Val = p.l.NextTokenKind(TOKEN_IDENTIFIER).Content
	} else {
		stmt.Val = firstName
	}
	p.l.NextTokenKind(TOKEN_SEP_COLON) // :

	stmt.Iterator = parseExp(p)

	p.l.NextTokenKind(TOKEN_SEP_RPAREN) // )
	p.l.ConsumeIf(TOKEN_SEP_SEMI)

	stmt.Block = p.parseBlockStmt()

	return
}

func (p *Parser) parseWhileStmt() (stmt *ast.WhileStmt) {
	p.l.NextToken()
	stmt = new(ast.WhileStmt)
	stmt.Condition = p.parseExpBlock()
	p.l.ConsumeIf(TOKEN_SEP_SEMI)
	stmt.Block = p.parseBlockStmt()
	return
}

func (p *Parser) parseFallthroughStmt() (stmt *ast.FallthroughStmt) {
	p.l.NextToken()
	return new(ast.FallthroughStmt)
}

func (p *Parser) parseGotoStmt() (stmt *ast.GotoStmt) {
	p.l.NextToken()
	if !p.l.Expect(TOKEN_IDENTIFIER) {
		panic("")
	}
	stmt = &ast.GotoStmt{Label: p.l.NextToken().Content}
	return
}

func (p *Parser) parseNameList() (names []string) {
	names = append(names, p.parseName())
	for p.l.ConsumeIf(TOKEN_SEP_COMMA) {
		names = append(names, p.parseName())
	}
	return names
}

func (p *Parser) parseName() string {
	if !p.l.Expect(TOKEN_IDENTIFIER) {
		panic(p.l.Line())
	}
	return p.l.NextToken().Content
}

// stmt ::= func ID funcBody ';          			   	// function definition
//	      | funcLiteral '(' [expList] ')' callTail ';'  // anonymous function call
func (p *Parser) parseFunc() ast.Stmt {
	p.l.NextToken() // func
	switch p.l.LookAhead().Kind {
	case TOKEN_IDENTIFIER:
		return p.parseFuncDefStmt()
	case TOKEN_SEP_LPAREN:
		return p.parseAnonymousFuncCallStmt()
	default:
		panic(p.l.Line())
	}
}

func (p *Parser) parseFuncDefStmt() (stmt *ast.FuncDefStmt) {
	stmt = new(ast.FuncDefStmt)
	stmt.Name = p.l.NextToken().Content
	stmt.FuncLiteral = p.parseFuncLiteral()
	p.FuncDefs = append(p.FuncDefs, stmt)
	return
}

func (p *Parser) parseFuncLiteral() (literal ast.FuncLiteral) {
	var defaultValue bool
	p.l.NextTokenKind(TOKEN_SEP_LPAREN) // (

	literal.Parameters, defaultValue = p.parseParameters()

	if p.l.ConsumeIf(TOKEN_SEP_VARARG) { // ...
		if defaultValue {
			// TODO error
			panic(p.l.Line())
		}
		literal.VaArgs = p.l.NextTokenKind(TOKEN_IDENTIFIER).Content
	}
	p.l.NextTokenKind(TOKEN_SEP_RPAREN) // )

	literal.Block = p.parseBlock()
	return
}

func (p *Parser) parseParameters() (pars []ast.Parameter, defaultValue bool) {
	for p.l.Expect(TOKEN_IDENTIFIER) {
		par := ast.Parameter{}
		par.Name = p.l.NextToken().Content
		if p.l.Expect(TOKEN_OP_ASSIGN) {
			defaultValue = true
			p.l.NextToken()
			switch token := p.l.NextToken(); token.Kind {
			case TOKEN_KW_TRUE:
				par.Default = true
			case TOKEN_KW_FALSE:
				par.Default = false
			case TOKEN_STRING, TOKEN_NUMBER:
				par.Default = token.Value
			default:
				panic("invalid default value") // TODO
			}
		} else if defaultValue {
			// TODO error
			panic(p.l.Line())
		}
		pars = append(pars, par)
		if p.l.Expect(TOKEN_SEP_COMMA) {
			p.l.NextToken()
		}
	}
	return
}

func (p *Parser) parseAnonymousFuncCallStmt() (stmt *ast.AnonymousFuncCallStmt) {
	stmt = new(ast.AnonymousFuncCallStmt)
	stmt.FuncLiteral = p.parseFuncLiteral()
	stmt.CallTails = []ast.CallTail{{Args: parseExpListBlock(p)}}
	stmt.CallTails = append(stmt.CallTails, p.parseCallTails()...)
	return
}

// . or [ or (
func expectTail(token *Token) bool {
	return token.Kind == TOKEN_SEP_DOT || token.Kind == TOKEN_SEP_LBRACK || token.Kind == TOKEN_SEP_LPAREN
}

func (p *Parser) parseCallTails() (tails []ast.CallTail) {
	for expectTail(p.l.LookAhead()) {
		var tail ast.CallTail
		tail.Attrs = p.parseAttrTail()
		tail.Args = parseExpListBlock(p)
		tails = append(tails, tail)
	}
	return
}

func (p *Parser) parseBreakStmt() (stmt *ast.BreakStmt) {
	p.l.NextToken()
	return new(ast.BreakStmt)
}

func (p *Parser) parseContinueStmt() (stmt *ast.ContinueStmt) {
	p.l.NextToken()
	return new(ast.ContinueStmt)
}

func (p *Parser) parseReturnStmt() (stmt *ast.ReturnStmt) {
	stmt = new(ast.ReturnStmt)
	p.l.NextToken()
	if p.l.Expect(TOKEN_SEP_RCURLY) || p.l.Expect(TOKEN_SEP_SEMI) {
		return
	}
	stmt.Args = parseExpList(p)
	return
}
