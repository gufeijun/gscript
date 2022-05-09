package codegen

import (
	"encoding/binary"
	"fmt"
	"gscript/complier/ast"
	"gscript/complier/parser"
	"gscript/proto"
	"unsafe"
)

func Gen(parser *parser.Parser) (text []proto.Instruction, consts []interface{}, funcs []proto.Func) {
	prog := parser.Parse()
	ctx := newContext(parser)

	// make all enum statements global
	genEnumStmt(parser.EnumStmts, ctx)

	genBlockStmts(prog.BlockStmts, ctx)
	ctx.writeIns(proto.INS_STOP)

	genFuncDefStmts(parser.FuncDefs, ctx)

	text = *(*[]proto.Instruction)((unsafe.Pointer(&ctx.buf)))
	consts = ctx.ct.Constants
	funcs = ctx.ft.funcTable
	return
}

func genFuncDefStmts(stmts []*ast.FuncDefStmt, ctx *Context) {
	for i, stmt := range stmts {
		ctx.ft.funcTable[i].Addr = ctx.textSize()
		ctx.renew()
		genFuncDefStmt(stmt, ctx)
	}
}

func genFuncDefStmt(stmt *ast.FuncDefStmt, ctx *Context) {
	if stmt.VaArgs != "" {
		ctx.insPushName(stmt.VaArgs)
	}
	for i := len(stmt.Parameters) - 1; i >= 0; i-- {
		ctx.insPushName(stmt.Parameters[i].Name)
	}
	genBlockStmts(stmt.Block.Blocks, ctx)
	if !ctx.returnAtEnd {
		genReturnStmt(&ast.ReturnStmt{}, ctx)
	}
}

func genEnumStmt(stmts []*ast.EnumStmt, ctx *Context) {
	for _, stmt := range stmts {
		for i := range stmt.Names {
			ctx.ct.saveEnum(stmt.Names[i], stmt.Values[i])
		}
	}
}

func genBlockStmts(stmts []ast.BlockStmt, ctx *Context) (varDecl bool) {
	var gotos []unhandledGoto
	for _, stmt := range stmts {
		ctx.returnAtEnd = false
		if block, ok := stmt.(ast.Block); ok {
			genBlock(block, ctx)
			continue
		}
		switch stmt := stmt.(type) {
		case *ast.VarDeclStmt:
			genVarDeclStmt(stmt, ctx)
			varDecl = true
		case *ast.VarAssignStmt:
			genVarAssignStmt(stmt, ctx)
		case *ast.IfStmt:
			genIfStmt(stmt, ctx)
		case *ast.WhileStmt:
			genWhileStmt(stmt, ctx)
		case *ast.ForStmt:
			genForStmt(stmt, ctx)
		case *ast.BreakStmt:
			genBreakStmt(stmt, ctx)
		case *ast.ContinueStmt:
			genContinueStmt(stmt, ctx)
		case *ast.SwitchStmt:
			genSwitchStmt(stmt, ctx)
		case *ast.FallthroughStmt:
			genFallthroughStmt(stmt, ctx)
		case *ast.NamedFuncCallStmt:
			genFuncCallStmt(stmt, ctx)
		case *ast.ReturnStmt:
			genReturnStmt(stmt, ctx)
		case *ast.LabelStmt:
			ctx.validLabels[stmt.Name] = label{stmt.Name, ctx.textSize(), *ctx.nt.nameIdx}
			// when exit block, make labels inside block invalid
			defer func() { delete(ctx.validLabels, stmt.Name) }()
		case *ast.GotoStmt:
			gotos = append(gotos, genGotoStmt(stmt, ctx))
		case *ast.EnumStmt, *ast.FuncDefStmt:
			continue
		default:
			panic(fmt.Sprintf("do not support stmt:%T", stmt))
		}
	}
	handleGoto(ctx, gotos)
	return
}

func genReturnStmt(stmt *ast.ReturnStmt, ctx *Context) {
	ctx.returnAtEnd = true
	for _, exp := range stmt.Args {
		genExp(exp, ctx, 1)
	}
	ctx.insReturn(uint32(len(stmt.Args)))
}

func genFuncCallStmt(stmt *ast.NamedFuncCallStmt, ctx *Context) {
	last := len(stmt.CallTails) - 1
	for i, callTail := range stmt.CallTails {
		var wantRetCnt byte
		if i != last {
			wantRetCnt = 1
		}

		genExps(callTail.Args, ctx, len(callTail.Args))

		// function
		if i == 0 {
			ctx.insLoadName(stmt.Prefix)
		}
		for _, attr := range callTail.Attrs {
			genExp(attr, ctx, 1)
			ctx.writeIns(proto.INS_BINARY_ATTR)
		}

		// call function
		ctx.insCall(wantRetCnt, byte(len(callTail.Args)))
	}
}

func handleGoto(ctx *Context, gotos []unhandledGoto) {
	writeUint := func(start int, val uint32) {
		binary.LittleEndian.PutUint32(ctx.buf[start:start+4], val)
	}
	for _, _goto := range gotos {
		label, ok := ctx.validLabels[_goto.label]
		if !ok {
			panic(fmt.Sprintf("invalid goto label: %s", _goto.label))
		}
		writeUint(_goto.jumpPos, label.addr)
		writeUint(_goto.resizePos, label.nameTableSize)
	}
}

func genGotoStmt(stmt *ast.GotoStmt, ctx *Context) unhandledGoto {
	resizePos := ctx.insResizeNameTable(0)
	jumpPos := ctx.insJumpAbs(0)
	return unhandledGoto{
		label:     stmt.Label,
		resizePos: resizePos,
		jumpPos:   jumpPos,
	}
}

func genFallthroughStmt(stmt *ast.FallthroughStmt, ctx *Context) {
	b := ctx.bs.latestSwitch()
	ctx.insResizeNameTable(b.nameCnt)
	pos := ctx.insJumpAbs(0)
	b._fallthrough = &pos
}

func genContinueStmt(stmt *ast.ContinueStmt, ctx *Context) {
	b := ctx.bs.latestFor()
	b.continues = append(b.continues, ctx.insJumpAbs(0))
}

func genBreakStmt(stmt *ast.BreakStmt, ctx *Context) {
	b := ctx.bs.top()
	var breaks *[]int
	if fb, ok := b.(*forBlock); ok {
		breaks = &fb.breaks
	} else {
		sb := b.(*switchBlock)
		breaks = &sb.breaks
		ctx.insResizeNameTable(sb.nameCnt)
	}
	*breaks = append(*breaks, ctx.insJumpAbs(0))
}

/*
------------------------------
switch (e0){
case e1,e2:
	case_block0
case e3:
	case_block1
default:
	case_block2
}
other_code
------------------------------
the code above will be translated to following code:

start:
	genExp(e0)
	genExp(e1)
	jump_case p0
	genExp(e2)
	jump_case p0
	genExp(e3)
	jump_case p1
	jump p2
	#jump end	replace "jump p1" with this when default doesn't appear
p0:
	case_block0
	jump end
p1:
	case_block1
	jump end
p2:
	case_block2
end:
	pop_top
	other_code
*/
func genSwitchStmt(stmt *ast.SwitchStmt, ctx *Context) {
	ctx.bs.pushSwitch(*ctx.nt.nameIdx)
	sb := ctx.bs.top().(*switchBlock)

	var pos_ptrs []int
	var end_ptrs []int
	genExp(stmt.Value, ctx, 1)
	for i, exps := range stmt.Cases {
		for _, exp := range exps {
			genExp(exp, ctx, 1)
			// if case block is nil, just jump to end
			if stmt.Blocks[i] == nil {
				end_ptrs = append(end_ptrs, ctx.insJumpCase(0))
			} else {
				pos_ptrs = append(pos_ptrs, ctx.insJumpCase(0))
			}
		}
	}
	if stmt.Default != nil {
		pos_ptrs = append(pos_ptrs, ctx.insJumpAbs(0))
	} else {
		end_ptrs = append(end_ptrs, ctx.insJumpAbs(0))
	}
	var i int
	for j, stmts := range stmt.Blocks {
		if sb._fallthrough != nil {
			pos := *sb._fallthrough
			sb._fallthrough = nil
			if stmts == nil {
				end_ptrs = append(end_ptrs, pos)
				continue
			}
			ctx.setAddr(pos, ctx.textSize())
		}
		if stmts == nil {
			continue
		}
		for k := 0; k < len(stmt.Cases[j]); k++ {
			ctx.setAddr(pos_ptrs[i], ctx.textSize())
			i++
		}
		genStmtsWithBlock(stmts, ctx)
		if stmt.Default != nil || j != len(stmt.Blocks)-1 {
			end_ptrs = append(end_ptrs, ctx.insJumpAbs(0))
		}
	}
	if stmt.Default != nil {
		if sb._fallthrough != nil {
			pos := *sb._fallthrough
			sb._fallthrough = nil
			ctx.setAddr(pos, ctx.textSize())
		}
		ctx.setAddr(pos_ptrs[i], ctx.textSize())
		genStmtsWithBlock(stmt.Default, ctx)
	}
	if sb._fallthrough != nil {
		panic("invalid fallthrough") // TODO
	}
	end := ctx.textSize()
	for _, end_ptr := range end_ptrs {
		ctx.setAddr(end_ptr, end)
	}
	ctx.insPopTop()
	for i := range sb.breaks {
		ctx.setAddr(sb.breaks[i], end)
	}
	ctx.bs.pop()
}

/*
------------------------------
for(e1;condition;e2){
	block_code
}
other_code
------------------------------
the code above will be translated to following code:

start:
	e1
p0:
	genExp(condition)
	JUMP_IF p1
	JUMP p2
p1:
	block_code
p3:
	e2
	JUMP p0
p2:
	other_code

break => jump to p2
continue => jump to p3
*/

func genForStmt(stmt *ast.ForStmt, ctx *Context) {
	if stmt.Condition == nil {
		stmt.Condition = &ast.TrueExp{}
	}
	ctx.enterBlock()
	startSize := *ctx.nt.nameIdx
	var varDecl bool
	// e1
	if stmt.AsgnStmt != nil {
		genVarAssignStmt(stmt.AsgnStmt, ctx)
	} else if stmt.DeclStmt != nil {
		genVarDeclStmt(stmt.DeclStmt, ctx)
		varDecl = true
	}
	curSize := *ctx.nt.nameIdx
	ctx.bs.pushFor(curSize)

	p0 := ctx.textSize()
	genExp(stmt.Condition, ctx, 1)
	p1ptr := ctx.insJumpIf(0)
	p2ptr := ctx.insJumpAbs(0)
	ctx.setAddr(p1ptr, ctx.textSize())
	varDecl = genBlockStmts(stmt.Block.Blocks, ctx) || varDecl
	p3 := ctx.textSize()
	ctx.insResizeNameTable(curSize)
	if stmt.ForTail != nil {
		genVarAssignStmt(stmt.ForTail, ctx) // e2
	}
	ctx.insJumpAbs(p0)
	p2 := ctx.textSize()
	ctx.setAddr(p2ptr, p2)

	ctx.leaveBlock(startSize, varDecl)

	fs := ctx.bs.top().(*forBlock)
	for i := range fs.breaks {
		ctx.setAddr(fs.breaks[i], p2)
	}
	for i := range fs.continues {
		ctx.setAddr(fs.continues[i], p3)
	}
	ctx.bs.pop()
}

func genWhileStmt(stmt *ast.WhileStmt, ctx *Context) {
	genForStmt(&ast.ForStmt{
		Condition: stmt.Condition,
		Block:     stmt.Block,
	}, ctx)
}

/*
------------------------------
if(e0) block0
elif(e1) block1
else block2
other_code
------------------------------
the code above will be translated to following code:

start:
	e0
	JUMP_IF p0
	e1
	JUMP_IF p1
	trueExp
	JUMP_IF p2
	JUMP end
p0:
	block0
	JUMP end
p1:
	block1
	JUMP end
p2:
	block2
end:
	other_code

*/
func genIfStmt(stmt *ast.IfStmt, ctx *Context) {
	var jf_addr_ptrs []int // jump_if
	var ja_addr_ptrs []int // jump_abs
	for _, condition := range stmt.Conditions {
		genExp(condition, ctx, 1)
		jf_addr_ptrs = append(jf_addr_ptrs, ctx.insJumpIf(0))
	}
	ja_addr_ptrs = append(ja_addr_ptrs, ctx.insJumpAbs(0))
	last := len(stmt.Blocks) - 1
	for i := 0; i < len(stmt.Blocks); i++ {
		ctx.setAddr(jf_addr_ptrs[i], ctx.textSize())
		genBlock(stmt.Blocks[i], ctx)
		if i != last {
			ja_addr_ptrs = append(ja_addr_ptrs, ctx.insJumpAbs(0))
		}
	}
	end := ctx.textSize()
	for _, pos := range ja_addr_ptrs {
		ctx.setAddr(pos, end)
	}
}

func genVarDeclStmt(stmt *ast.VarDeclStmt, ctx *Context) {
	genExps(stmt.Rights, ctx, len(stmt.Lefts))
	for i := len(stmt.Lefts) - 1; i >= 0; i-- {
		ctx.insPushName(stmt.Lefts[i])
	}
}

func genVarAssignStmt(stmt *ast.VarAssignStmt, ctx *Context) {
	genExps(stmt.Rights, ctx, len(stmt.Lefts))

	for i := len(stmt.Lefts) - 1; i >= 0; i-- {
		target := stmt.Lefts[i]
		length := len(target.Attrs)
		if length == 0 {
			if stmt.AssignOp != ast.ASIGN_OP_ASSIGN {
				genExp(&ast.NameExp{Name: target.Prefix}, ctx, 1)
				ctx.writeIns(proto.INS_ROT_TWO)
				ctx.writeIns(byte(stmt.AssignOp-ast.ASIGN_OP_ASSIGN) + proto.INS_BINARY_START)
			}
			ctx.insStoreName(target.Prefix)
			continue
		}
		genExp(target.Attrs[len(target.Attrs)-1], ctx, 1)
		ctx.insLoadName(target.Prefix)
		for i := 0; i < length-1; i++ {
			genExp(target.Attrs[i], ctx, 1)
			ctx.writeIns(proto.INS_BINARY_ATTR)
		}
		ctx.writeIns(byte(stmt.AssignOp-ast.ASIGN_OP_START) + proto.INS_ATTR_ASSIGN_START)
	}
}

func genBlock(block ast.Block, ctx *Context) {
	genStmtsWithBlock(block.Blocks, ctx)
}

func genStmtsWithBlock(stmts []ast.BlockStmt, ctx *Context) {
	ctx.enterBlock()
	size := *ctx.nt.nameIdx
	varDecl := genBlockStmts(stmts, ctx)
	ctx.leaveBlock(size, varDecl)
}