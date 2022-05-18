package codegen

import (
	"fmt"
	"gscript/complier/ast"
	"gscript/complier/parser"
	"gscript/proto"
	"unsafe"
)

type Import struct {
	ProtoNumber uint32
	Alias       string
}

func Gen(parser *parser.Parser, prog *ast.Program, imports []Import, mainProto bool) proto.Proto {
	ctx := newContext(parser)

	genImports(imports, ctx)
	// make all enum and class statements global
	genEnumStmt(parser.EnumStmts, ctx)
	genClassStmts(parser.ClassStmts, ctx)

	genBlockStmts(prog.BlockStmts, ctx)
	if !mainProto {
		genExport(prog.Export, ctx)
	} else {
		ctx.writeIns(proto.INS_STOP)
	}

	return proto.Proto{
		Text:           bytesToInstructions(ctx.frame.text),
		Consts:         ctx.ct.Constants,
		Funcs:          ctx.ft.funcTable,
		AnonymousFuncs: ctx.ft.anonymousFuncs,
	}
}

func genExport(export ast.Export, ctx *Context) {
	exp := export.Exp
	if export.Exp == nil {
		exp = &ast.NilExp{}
	}
	genExp(exp, ctx, 1)
	ctx.writeIns(proto.INS_EXPORT)
}

func genImports(imports []Import, ctx *Context) {
	for _, _import := range imports {
		ctx.insLoadProto(_import.ProtoNumber)
		ctx.insPushName(_import.Alias)
	}
}

func genClassStmts(stmts []*ast.ClassStmt, ctx *Context) {
	for i, stmt := range stmts {
		genClassStmt(stmt, ctx)
		ctx.ft.anonymousFuncs[i].Info.Text = bytesToInstructions(ctx.frame.text)
		ctx.frame = newStackFrame()
	}
}

func genClassStmt(stmt *ast.ClassStmt, ctx *Context) {
	__self := stmt.Constructor
	if __self != nil {
		collectArgs(&__self.FuncLiteral, ctx)
	}
	ctx.insCopyName("this")
	for i := range stmt.AttrName {
		ctx.insLoadConst(stmt.AttrName[i])
		genExp(stmt.AttrValue[i], ctx, 1)
		ctx.writeIns(proto.INS_STORE_KV)
	}

	// codes of __self
	if __self != nil {
		genBlockStmts(__self.Block.Blocks, ctx)
	}
	if !ctx.frame.returnAtEnd {
		genReturnStmt(&ast.ReturnStmt{}, ctx)
	}
}

func genBlockStmts(stmts []ast.BlockStmt, ctx *Context) (varDecl bool) {
	var gotos []unhandledGoto
	for _, stmt := range stmts {
		ctx.frame.returnAtEnd = false
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
		case *ast.AnonymousFuncCallStmt:
			genAnonymousFuncCallStmt(stmt, ctx)
		case *ast.FuncDefStmt:
			genFuncDefStmt(stmt, ctx)
		case *ast.LoopStmt:
			// TODO
			panic("do not support loop statement for now")
		case *ast.LabelStmt:
			ctx.frame.validLabels[stmt.Name] = label{stmt.Name, ctx.textSize(), *ctx.frame.nt.nameIdx}
			// when exit block, make labels inside block invalid
			defer func() { delete(ctx.frame.validLabels, stmt.Name) }()
		case *ast.GotoStmt:
			gotos = append(gotos, genGotoStmt(stmt, ctx))
		case *ast.EnumStmt, *ast.ClassStmt:
			continue
		default:
			panic(fmt.Sprintf("do not support stmt:%T", stmt))
		}
	}
	handleGoto(ctx, gotos)
	return
}

func genFuncDefStmt(stmt *ast.FuncDefStmt, ctx *Context) {
	funcIdx := ctx.ft.funcMap[stmt.Name]

	genFuncLiteral(&stmt.FuncLiteral, ctx, funcIdx, false)
}

func collectArgs(literal *ast.FuncLiteral, ctx *Context) {
	if literal.VaArgs != "" {
		ctx.insPushName(literal.VaArgs)
	}
	for i := len(literal.Parameters) - 1; i >= 0; i-- {
		ctx.insPushName(literal.Parameters[i].Name)
	}
}

func genFuncLiteral(literal *ast.FuncLiteral, ctx *Context, idx uint32, anonymous bool) {
	ctx.pushFrame(anonymous, int(idx))
	collectArgs(literal, ctx)
	genBlockStmts(literal.Block.Blocks, ctx)
	if !ctx.frame.returnAtEnd {
		genReturnStmt(&ast.ReturnStmt{}, ctx)
	}

	bindUpValue(ctx, idx, anonymous)
}

func bytesToInstructions(text []byte) []proto.Instruction {
	return *((*[]proto.Instruction)(unsafe.Pointer(&text)))
}

func bindUpValue(ctx *Context, funcIdx uint32, anonymous bool) {
	oldFrame := ctx.popFrame()
	upValues := oldFrame.vt.upValues

	if anonymous {
		ctx.ft.anonymousFuncs[ctx.frame.nowParsingAnonymous].Info.Text = bytesToInstructions(oldFrame.text)
		ft := ctx.ft.anonymousFuncs
		for _, upValue := range upValues {
			vptr := getUpValueIdx(ctx.frame, ctx, &upValue)
			ft[funcIdx].UpValues = append(ft[funcIdx].UpValues, vptr)
		}
	} else {
		ctx.ft.funcTable[funcIdx].Info.Text = bytesToInstructions(oldFrame.text)
		ft := ctx.ft.funcTable
		for _, upValue := range upValues {
			ft[funcIdx].UpValues = append(ft[funcIdx].UpValues, upValue.nameIdx)
		}
	}
}

func getUpValueIdx(frame *StackFrame, ctx *Context, upValue *UpValue) proto.UpValuePtr {
	if upValue.level == 0 {
		return proto.UpValuePtr{true, upValue.nameIdx}
	}
	upValue.level--
	valueIdx, _, ok := frame.vt.get(upValue.name)
	if ok {
		return proto.UpValuePtr{false, valueIdx}
	}
	valueIdx = frame.vt.set(upValue.name, upValue.level, upValue.nameIdx)
	return proto.UpValuePtr{false, valueIdx}
}

func genAnonymousFuncCallStmt(stmt *ast.AnonymousFuncCallStmt, ctx *Context) {
	genFuncCall(&ast.FuncLiteralExp{FuncLiteral: stmt.FuncLiteral}, stmt.CallTails, ctx)
}

func genReturnStmt(stmt *ast.ReturnStmt, ctx *Context) {
	ctx.frame.returnAtEnd = true
	for _, exp := range stmt.Args {
		genExp(exp, ctx, 1)
	}
	ctx.insReturn(uint32(len(stmt.Args)))
}

func genFuncCallStmt(stmt *ast.NamedFuncCallStmt, ctx *Context) {
	genFuncCall(&ast.NameExp{Name: stmt.Prefix}, stmt.CallTails, ctx)
}

func genFuncCall(exp ast.Exp, callTails []ast.CallTail, ctx *Context) {
	last := len(callTails) - 1
	for i, callTail := range callTails {
		var wantRetCnt byte
		if i != last {
			wantRetCnt = 1
		}

		genExps(callTail.Args, ctx, len(callTail.Args))
		// function
		if i == 0 {
			genExp(exp, ctx, 1)
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
	for _, _goto := range gotos {
		label, ok := ctx.frame.validLabels[_goto.label]
		if !ok {
			panic(fmt.Sprintf("invalid goto label: %s", _goto.label))
		}
		ctx.setSteps(_goto.jumpPos, label.addr)
		ctx.setSteps(_goto.resizePos, label.nameTableSize)
	}
}

func genGotoStmt(stmt *ast.GotoStmt, ctx *Context) unhandledGoto {
	resizePos := ctx.insResizeNameTable(0)
	jumpPos := ctx.insJumpRel(0)
	return unhandledGoto{
		label:     stmt.Label,
		resizePos: resizePos,
		jumpPos:   jumpPos,
	}
}

func genFallthroughStmt(stmt *ast.FallthroughStmt, ctx *Context) {
	b := ctx.frame.bs.latestSwitch()
	ctx.insResizeNameTable(b.nameCnt)
	pos := ctx.insJumpRel(0)
	b._fallthrough = &pos
}

func genContinueStmt(stmt *ast.ContinueStmt, ctx *Context) {
	b := ctx.frame.bs.latestFor()
	b.continues = append(b.continues, ctx.insJumpRel(0))
}

func genBreakStmt(stmt *ast.BreakStmt, ctx *Context) {
	b := ctx.frame.bs.top()
	var breaks *[]int
	if fb, ok := b.(*forBlock); ok {
		breaks = &fb.breaks
	} else {
		sb := b.(*switchBlock)
		breaks = &sb.breaks
		ctx.insResizeNameTable(sb.nameCnt)
	}
	*breaks = append(*breaks, ctx.insJumpRel(0))
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
	ctx.frame.bs.pushSwitch(*ctx.frame.nt.nameIdx)
	sb := ctx.frame.bs.top().(*switchBlock)

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
		pos_ptrs = append(pos_ptrs, ctx.insJumpRel(0))
	} else {
		end_ptrs = append(end_ptrs, ctx.insJumpRel(0))
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
			ctx.setSteps(pos, ctx.textSize())
		}
		if stmts == nil {
			continue
		}
		for k := 0; k < len(stmt.Cases[j]); k++ {
			ctx.setSteps(pos_ptrs[i], ctx.textSize())
			i++
		}
		genStmtsWithBlock(stmts, ctx)
		if stmt.Default != nil || j != len(stmt.Blocks)-1 {
			end_ptrs = append(end_ptrs, ctx.insJumpRel(0))
		}
	}
	if stmt.Default != nil {
		if sb._fallthrough != nil {
			pos := *sb._fallthrough
			sb._fallthrough = nil
			ctx.setSteps(pos, ctx.textSize())
		}
		ctx.setSteps(pos_ptrs[i], ctx.textSize())
		genStmtsWithBlock(stmt.Default, ctx)
	}
	if sb._fallthrough != nil {
		panic("invalid fallthrough") // TODO
	}
	end := ctx.textSize()
	for _, end_ptr := range end_ptrs {
		ctx.setSteps(end_ptr, end)
	}
	ctx.insPopTop()
	for i := range sb.breaks {
		ctx.setSteps(sb.breaks[i], end)
	}
	ctx.frame.bs.pop()
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
	startSize := *ctx.frame.nt.nameIdx
	var varDecl bool
	// e1
	if stmt.AsgnStmt != nil {
		genVarAssignStmt(stmt.AsgnStmt, ctx)
	} else if stmt.DeclStmt != nil {
		genVarDeclStmt(stmt.DeclStmt, ctx)
		varDecl = true
	}
	curSize := *ctx.frame.nt.nameIdx
	ctx.frame.bs.pushFor(curSize)

	p0 := ctx.textSize()
	genExp(stmt.Condition, ctx, 1)
	p1ptr := ctx.insJumpIf(0)
	p2ptr := ctx.insJumpRel(0)
	ctx.setSteps(p1ptr, ctx.textSize())
	varDecl = genBlockStmts(stmt.Block.Blocks, ctx) || varDecl
	p3 := ctx.textSize()
	ctx.insResizeNameTable(curSize)
	if stmt.ForTail != nil {
		genVarAssignStmt(stmt.ForTail, ctx) // e2
	}
	ctx.insJumpRel(p0)
	p2 := ctx.textSize()
	ctx.setSteps(p2ptr, p2)

	ctx.leaveBlock(startSize, varDecl)

	fs := ctx.frame.bs.top().(*forBlock)
	for i := range fs.breaks {
		ctx.setSteps(fs.breaks[i], p2)
	}
	for i := range fs.continues {
		ctx.setSteps(fs.continues[i], p3)
	}
	ctx.frame.bs.pop()
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
	var jr_addr_ptrs []int // jump_rel
	for _, condition := range stmt.Conditions {
		genExp(condition, ctx, 1)
		jf_addr_ptrs = append(jf_addr_ptrs, ctx.insJumpIf(0))
	}
	jr_addr_ptrs = append(jr_addr_ptrs, ctx.insJumpRel(0))
	last := len(stmt.Blocks) - 1
	for i := 0; i < len(stmt.Blocks); i++ {
		ctx.setSteps(jf_addr_ptrs[i], ctx.textSize())
		genBlock(stmt.Blocks[i], ctx)
		if i != last {
			jr_addr_ptrs = append(jr_addr_ptrs, ctx.insJumpRel(0))
		}
	}
	end := ctx.textSize()
	for _, pos := range jr_addr_ptrs {
		ctx.setSteps(pos, end)
	}
}

func genVarDeclStmt(stmt *ast.VarDeclStmt, ctx *Context) {
	genExps(stmt.Rights, ctx, len(stmt.Lefts))
	for i := len(stmt.Lefts) - 1; i >= 0; i-- {
		ctx.insPushName(stmt.Lefts[i])
	}
}

func needRotate(op int) bool {
	return op == ast.ASIGN_OP_SUBEQ || op == ast.ASIGN_OP_ADDEQ ||
		op == ast.ASIGN_OP_DIVEQ || op == ast.ASIGN_OP_MODEQ
}

func genVarAssignStmt(stmt *ast.VarAssignStmt, ctx *Context) {
	genExps(stmt.Rights, ctx, len(stmt.Lefts))

	for i := len(stmt.Lefts) - 1; i >= 0; i-- {
		target := stmt.Lefts[i]
		length := len(target.Attrs)
		if length == 0 {
			if stmt.AssignOp != ast.ASIGN_OP_ASSIGN {
				genExp(&ast.NameExp{Name: target.Prefix}, ctx, 1)
				if needRotate(stmt.AssignOp) {
					ctx.writeIns(proto.INS_ROT_TWO)
				}
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
	size := *ctx.frame.nt.nameIdx
	varDecl := genBlockStmts(stmts, ctx)
	ctx.leaveBlock(size, varDecl)
}

func genEnumStmt(stmts []*ast.EnumStmt, ctx *Context) {
	for _, stmt := range stmts {
		for i := range stmt.Names {
			ctx.ct.saveEnum(stmt.Names[i], stmt.Values[i])
		}
	}
}
