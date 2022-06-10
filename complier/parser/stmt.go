package parser

import (
	"gscript/complier/ast"
	"gscript/complier/token"
)

func (p *Parser) parseStmt(atTop bool) ast.Stmt {
	var stmt ast.Stmt
	ahead := p.l.LookAhead()
	switch ahead.Kind {
	case token.TOKEN_KW_LET:
		stmt = p.parseVarDeclStmt()
	case token.TOKEN_IDENTIFIER:
		stmt = p.parseVarOpOrLabel()
	case token.TOKEN_KW_FUNC:
		stmt = p.parseFunc()
		if !atTop {
			if _stmt, ok := stmt.(*ast.FuncDefStmt); ok {
				stmt = toVarDeclStmt(_stmt)
			}
		}
	case token.TOKEN_KW_BREAK:
		stmt = p.parseBreakStmt()
	case token.TOKEN_KW_CONTINUE:
		stmt = p.parseContinueStmt()
	case token.TOKEN_KW_RETURN:
		stmt = p.parseReturnStmt()
	case token.TOKEN_KW_GOTO:
		stmt = p.parseGotoStmt()
	case token.TOKEN_KW_FALLTHROUGH:
		stmt = p.parseFallthroughStmt()
	case token.TOKEN_KW_LOOP:
		stmt = p.parseLoopStmt()
	case token.TOKEN_KW_WHILE:
		stmt = p.parseWhileStmt()
	case token.TOKEN_KW_FOR:
		stmt = p.parseForStmt()
	case token.TOKEN_KW_IF:
		stmt = p.parseIfStmt()
	case token.TOKEN_KW_CLASS:
		stmt = p.parseClassStmt()
	case token.TOKEN_KW_ENUM:
		stmt = p.parseEnumStmt()
	case token.TOKEN_KW_SWITCH:
		stmt = p.parseSwitchStmt()
	case token.TOKEN_OP_INC, token.TOKEN_OP_DEC:
		stmt = p.parseIncOrDecVar()
	case token.TOKEN_KW_TRY:
		stmt = p.parseTryCatchStmt()
	default:
		p.exit("unexpected token '%s' to make a statement", ahead.Content)
	}
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	return stmt
}

func (p *Parser) parseTryCatchStmt() (stmt *ast.TryCatchStmt) {
	stmt = new(ast.TryCatchStmt)
	p.l.NextToken()
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	stmt.TryBlocks = p.parseBlock().Blocks
	p.NextTokenKind(token.TOKEN_KW_CATCH)
	if p.ConsumeIf(token.TOKEN_SEP_LPAREN) {
		if p.Expect(token.TOKEN_IDENTIFIER) {
			stmt.CatchLine = p.l.Line()
			stmt.CatchValue = p.l.NextToken().Content
		}
		p.NextTokenKind(token.TOKEN_SEP_RPAREN)
	}
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	stmt.CatchBlocks = p.parseBlock().Blocks
	return stmt
}

func toVarDeclStmt(stmt *ast.FuncDefStmt) ast.Stmt {
	varDecl := new(ast.VarDeclStmt)
	varDecl.Lefts = []string{stmt.Name}
	varDecl.Rights = []ast.Exp{&ast.FuncLiteralExp{FuncLiteral: stmt.FuncLiteral}}
	return varDecl
}

// varDeclStmt ::= let id {,id} = exp {,exp} ;
func (p *Parser) parseVarDeclStmt() (stmt *ast.VarDeclStmt) {
	stmt = new(ast.VarDeclStmt)
	p.l.NextToken()
	stmt.Line = p.l.Line()
	stmt.Lefts = p.parseNameList() // id {,id}
	if !p.Expect(token.TOKEN_OP_ASSIGN) {
		stmt.Rights = make([]ast.Exp, len(stmt.Lefts))
		for i := range stmt.Rights {
			stmt.Rights[i] = &ast.NilExp{}
		}
		return
	}
	p.l.NextToken()
	stmt.Rights = parseExpList(p) // exp {,exp}
	if len(stmt.Lefts) < len(stmt.Rights) {
		stmt.Rights = stmt.Rights[:len(stmt.Lefts)]
	}
	return
}

// stmt ::= | varAssign	;						# case1: var {, var} assignOp exp {, exp}
//          | varIncOrDec ;						# case2: var++ or var--
//          | var '(' {explist} ')' callTail ;	# case3: function call
//          | ID ':'							# case4: label
func (p *Parser) parseVarOpOrLabel() ast.Stmt {
	name := p.l.NextToken().Content
	if p.ConsumeIf(token.TOKEN_SEP_COLON) { // case4
		label := &ast.LabelStmt{Name: name}
		return label
	}
	v := p._parseVar(name)
	switch p.l.LookAhead().Kind {
	case token.TOKEN_OP_INC, token.TOKEN_OP_DEC: // case2
		stmt := _incOrDec2Assign(p.l.NextToken().Kind, v)
		return stmt
	case token.TOKEN_SEP_LPAREN: // case3
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
	for p.ConsumeIf(token.TOKEN_SEP_COMMA) {
		stmt.Lefts = append(stmt.Lefts, p.parseVar())
	}
	token := p.l.NextToken()
	if !_isAsignOp(token.Kind) {
		p.exit("expect assign operation for assign statement, but got '%s'", token.Content)
	}
	stmt.AssignOp = token.Kind
	stmt.Rights = parseExpList(p)
	if len(stmt.Lefts) != len(stmt.Rights) {
		p.exit("expression count(%d) on the right of operator '%s' is not equal to variable count(%d) on the left",
			len(stmt.Rights), token.Content, len(stmt.Lefts))
	}
	return &stmt
}

// i++ i-- ==> i+=1 i-=1
func _incOrDec2Assign(kind int, v ast.Var) (stmt *ast.VarAssignStmt) {
	op := ast.ASIGN_OP_SUBEQ
	if kind == token.TOKEN_OP_INC {
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
	return p._parseVar(p.NextTokenKind(token.TOKEN_IDENTIFIER).Content)
}

func (p *Parser) parseAttrTail() (exps []ast.Exp) {
	for {
		switch p.l.LookAhead().Kind {
		case token.TOKEN_SEP_DOT:
			p.l.NextToken()
			// a.b ==> a["b"]
			token := p.NextTokenKind(token.TOKEN_IDENTIFIER)
			exps = append(exps, &ast.StringLiteralExp{Value: token.Content})
		case token.TOKEN_SEP_LBRACK: //[
			p.l.NextToken()
			exps = append(exps, parseExp(p))
			p.NextTokenKind(token.TOKEN_SEP_RBRACK)
		default:
			return
		}
	}
}

func (p *Parser) parseSwitchStmt() (stmt *ast.SwitchStmt) {
	stmt = new(ast.SwitchStmt)
	p.l.NextToken()                // switch
	stmt.Value = p.parseExpBlock() // (condition)
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	p.NextTokenKind(token.TOKEN_SEP_LCURLY) // {

	for p.ConsumeIf(token.TOKEN_KW_CASE) {
		stmt.Cases = append(stmt.Cases, parseExpList(p))
		p.NextTokenKind(token.TOKEN_SEP_COLON) // :
		stmt.Blocks = append(stmt.Blocks, p.parseBlockStmts(false))
	}

	if p.ConsumeIf(token.TOKEN_KW_DEFAULT) { // default
		p.NextTokenKind(token.TOKEN_SEP_COLON)
		stmt.Default = p.parseBlockStmts(false)
	}
	p.NextTokenKind(token.TOKEN_SEP_RCURLY) // }
	return
}

func (p *Parser) parseEnumStmt() (stmt *ast.EnumStmt) {
	stmt = new(ast.EnumStmt)
	p.l.NextToken()
	p.NextTokenKind(token.TOKEN_SEP_LCURLY)
	var enum int64
	var ok bool
	for {
		if p.ConsumeIf(token.TOKEN_SEP_RCURLY) {
			break
		}
		stmt.Lines = append(stmt.Lines, p.l.Line())
		stmt.Names = append(stmt.Names, p.NextTokenKind(token.TOKEN_IDENTIFIER).Content)
		if p.ConsumeIf(token.TOKEN_OP_ASSIGN) {
			t := p.l.NextToken()
			if enum, ok = t.Value.(int64); !ok {
				p.exit("expect enum value after '=', but got '%s'", t.Content)
			}
		}
		stmt.Values = append(stmt.Values, enum)
		enum++
		if p.ConsumeIf(token.TOKEN_SEP_RCURLY) {
			break
		}
		if p.Expect(token.TOKEN_SEP_COMMA) || p.Expect(token.TOKEN_SEP_SEMI) {
			p.l.NextToken()
		}
	}
	p.EnumStmts = append(p.EnumStmts, stmt)
	return
}

func (p *Parser) parseClassStmt() (stmt *ast.ClassStmt) {
	stmt = new(ast.ClassStmt)
	p.l.NextToken()
	stmt.Name = p.NextTokenKind(token.TOKEN_IDENTIFIER).Content
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	p.NextTokenKind(token.TOKEN_SEP_LCURLY)

	for p.Expect(token.TOKEN_IDENTIFIER) {
		attr := p.l.NextToken().Content
		var exp ast.Exp
		if p.Expect(token.TOKEN_SEP_LPAREN) {
			exp = &ast.FuncLiteralExp{FuncLiteral: p.parseFuncLiteral()}
		} else if p.ConsumeIf(token.TOKEN_OP_ASSIGN) {
			exp = parseExp(p)
		}
		if attr == "__self" {
			_exp, ok := exp.(*ast.FuncLiteralExp)
			if !ok {
				p.exit("__self of class '%s' should be a method", stmt.Name)
			}
			stmt.Constructor = _exp
		} else if exp != nil {
			stmt.AttrName = append(stmt.AttrName, attr)
			stmt.AttrValue = append(stmt.AttrValue, exp)
		}
		if p.Expect(token.TOKEN_SEP_SEMI) || p.Expect(token.TOKEN_SEP_COMMA) {
			p.l.NextToken()
		}
	}

	p.NextTokenKind(token.TOKEN_SEP_RCURLY)
	p.ClassStmts = append(p.ClassStmts, stmt)
	return
}

// single stmt to block
func (p *Parser) parseBlockStmt() (stmt ast.Block) {
	switch p.l.LookAhead().Kind {
	case token.TOKEN_SEP_LCURLY:
		return p.parseBlock()
	default:
		stmt.Blocks = append(stmt.Blocks, p.parseStmt(false))
		return
	}
}

func (p *Parser) parseIfStmt() (stmt *ast.IfStmt) {
	stmt = new(ast.IfStmt)
	p.l.NextToken() // if
	stmt.Conditions = append(stmt.Conditions, p.parseExpBlock())
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	stmt.Blocks = append(stmt.Blocks, p.parseBlockStmt()) // block | stmt

	for p.ConsumeIf(token.TOKEN_KW_ELIF) { // elif
		stmt.Conditions = append(stmt.Conditions, p.parseExpBlock())
		p.ConsumeIf(token.TOKEN_SEP_SEMI)
		stmt.Blocks = append(stmt.Blocks, p.parseBlockStmt())
	}

	// else ==> elif(true)
	if p.ConsumeIf(token.TOKEN_KW_ELSE) { // else
		stmt.Conditions = append(stmt.Conditions, &ast.TrueExp{})
		stmt.Blocks = append(stmt.Blocks, p.parseBlockStmt())
	}

	return
}

// (exp)
func (p *Parser) parseExpBlock() (exp ast.Exp) {
	p.NextTokenKind(token.TOKEN_SEP_LPAREN) // (
	exp = parseExp(p)                       // condition
	p.NextTokenKind(token.TOKEN_SEP_RPAREN) // )
	return
}

func _isAsignOp(kind int) bool {
	return kind > token.TOKEN_ASIGN_START && kind < token.TOKEN_ASIGN_END
}

func (p *Parser) parseVarAssignStmt(prefix string) (stmt *ast.VarAssignStmt) {
	return p._parseAssignStmt(p._parseVar(prefix))
}

func (p *Parser) parseForStmt() (stmt *ast.ForStmt) {
	stmt = new(ast.ForStmt)
	p.l.NextToken()
	p.NextTokenKind(token.TOKEN_SEP_LPAREN) // (

	switch ahead := p.l.LookAhead(); ahead.Kind {
	case token.TOKEN_KW_LET:
		stmt.DeclStmt = p.parseVarDeclStmt()
	case token.TOKEN_IDENTIFIER:
		stmt.AsgnStmt = p.parseVarAssignStmt(p.l.NextToken().Content)
	case token.TOKEN_SEP_SEMI:
		break
	default:
		p.exit("unexpected token '%s' for ForStatement", ahead.Content)
	}

	p.NextTokenKind(token.TOKEN_SEP_SEMI) // ;
	if !p.Expect(token.TOKEN_SEP_SEMI) {
		stmt.Condition = parseExp(p)
	}
	p.NextTokenKind(token.TOKEN_SEP_SEMI) // ;
	if !p.Expect(token.TOKEN_SEP_RPAREN) {
		stmt.ForTail = p.parseForTail()
	}
	p.NextTokenKind(token.TOKEN_SEP_RPAREN) // )
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	stmt.Block = p.parseBlockStmt()
	return
}

// forTail ::= varAssign | varIncOrDec | incOrDecVar
func (p *Parser) parseForTail() *ast.VarAssignStmt {
	if p.Expect(token.TOKEN_OP_INC) || p.Expect(token.TOKEN_OP_DEC) { // incOrDecVar
		return p.parseIncOrDecVar()
	}
	v := p.parseVar()
	if p.Expect(token.TOKEN_OP_INC) || p.Expect(token.TOKEN_OP_DEC) { // varIncOrDec
		return _incOrDec2Assign(p.l.NextToken().Kind, v)
	}
	return p._parseAssignStmt(v) // varAssign
}

func (p *Parser) parseLoopStmt() (stmt *ast.LoopStmt) {
	p.l.NextToken()
	p.NextTokenKind(token.TOKEN_SEP_LPAREN)
	p.NextTokenKind(token.TOKEN_KW_LET)

	stmt = new(ast.LoopStmt)
	firstName := p.NextTokenKind(token.TOKEN_IDENTIFIER).Content
	if p.ConsumeIf(token.TOKEN_SEP_COMMA) {
		stmt.Key = firstName
		stmt.Val = p.NextTokenKind(token.TOKEN_IDENTIFIER).Content
	} else {
		stmt.Val = firstName
	}
	p.NextTokenKind(token.TOKEN_SEP_COLON) // :

	stmt.Iterator = parseExp(p)

	p.NextTokenKind(token.TOKEN_SEP_RPAREN) // )
	p.ConsumeIf(token.TOKEN_SEP_SEMI)

	stmt.Block = p.parseBlockStmt()

	return
}

func (p *Parser) parseWhileStmt() (stmt *ast.WhileStmt) {
	p.l.NextToken()
	stmt = new(ast.WhileStmt)
	stmt.Condition = p.parseExpBlock()
	p.ConsumeIf(token.TOKEN_SEP_SEMI)
	stmt.Block = p.parseBlockStmt()
	return
}

func (p *Parser) parseFallthroughStmt() (stmt *ast.FallthroughStmt) {
	p.l.NextToken()
	stmt = new(ast.FallthroughStmt)
	stmt.Line = p.l.Line()
	return stmt
}

func (p *Parser) parseGotoStmt() *ast.GotoStmt {
	p.l.NextToken()
	label := p.NextTokenKind(token.TOKEN_IDENTIFIER).Content
	stmt := &ast.GotoStmt{
		Label: label,
		Line:  p.l.Line(),
	}
	return stmt
}

func (p *Parser) parseNameList() (names []string) {
	names = append(names, p.parseName())
	for p.ConsumeIf(token.TOKEN_SEP_COMMA) {
		names = append(names, p.parseName())
	}
	return names
}

func (p *Parser) parseName() string {
	return p.NextTokenKind(token.TOKEN_IDENTIFIER).Content
}

// stmt ::= func ID funcBody ';          			   	// function definition
//	      | funcLiteral '(' [expList] ')' callTail ';'  // anonymous function call
func (p *Parser) parseFunc() ast.Stmt {
	p.l.NextToken() // func
	ahead := p.l.LookAhead()
	switch ahead.Kind {
	case token.TOKEN_IDENTIFIER:
		return p.parseFuncDefStmt()
	case token.TOKEN_SEP_LPAREN:
		return p.parseAnonymousFuncCallStmt()
	default:
		p.exit("unexpexted token '%s' after keyword 'func' for function definition", ahead.Content)
	}
	return nil
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
	p.NextTokenKind(token.TOKEN_SEP_LPAREN) // (

	literal.Line = p.l.Line()
	literal.Parameters, defaultValue = p.parseParameters()

	if p.ConsumeIf(token.TOKEN_SEP_VARARG) { // ...
		vaArgs := p.NextTokenKind(token.TOKEN_IDENTIFIER).Content
		if defaultValue {
			p.exit("vaargs '%s' can not appear after parameter with default value", vaArgs)
		}
		literal.VaArgs = vaArgs
	}
	p.NextTokenKind(token.TOKEN_SEP_RPAREN) // )

	literal.Block = p.parseBlock()
	return
}

func (p *Parser) parseParameters() (pars []ast.Parameter, defaultValue bool) {
	for p.Expect(token.TOKEN_IDENTIFIER) {
		par := ast.Parameter{}
		par.Name = p.l.NextToken().Content
		if p.Expect(token.TOKEN_OP_ASSIGN) {
			defaultValue = true
			p.l.NextToken()
			switch t := p.l.NextToken(); t.Kind {
			case token.TOKEN_KW_TRUE:
				par.Default = true
			case token.TOKEN_KW_FALSE:
				par.Default = false
			case token.TOKEN_STRING, token.TOKEN_NUMBER:
				par.Default = t.Value
			case token.TOKEN_OP_SUB:
				par.Default = -p.NextTokenKind(token.TOKEN_NUMBER).Value.(int64)
			default:
				p.exit("invalid default value '%s' for parameter '%s'", t.Content, par.Name)
			}
		} else if defaultValue {
			p.exit("parameter '%s' without default value can not appear after those with default values", par.Name)
		}
		pars = append(pars, par)
		if p.Expect(token.TOKEN_SEP_COMMA) {
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
func expectTail(t *token.Token) bool {
	return t.Kind == token.TOKEN_SEP_DOT || t.Kind == token.TOKEN_SEP_LBRACK || t.Kind == token.TOKEN_SEP_LPAREN
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
	stmt = new(ast.BreakStmt)
	stmt.Line = p.l.Line()
	return stmt
}

func (p *Parser) parseContinueStmt() (stmt *ast.ContinueStmt) {
	p.l.NextToken()
	stmt = new(ast.ContinueStmt)
	stmt.Line = p.l.Line()
	return stmt
}

func (p *Parser) parseReturnStmt() (stmt *ast.ReturnStmt) {
	stmt = new(ast.ReturnStmt)
	p.l.NextToken()
	if p.Expect(token.TOKEN_SEP_RCURLY) || p.Expect(token.TOKEN_SEP_SEMI) {
		return
	}
	stmt.Args = parseExpList(p)
	return
}
