package codegen

import (
	"encoding/binary"
	"fmt"
	"gscript/compiler/ast"
	"gscript/proto"
	"os"
)

func genExps(exps []ast.Exp, ctx *Context, wantCnt int) {
	if len(exps) == 0 {
		for i := 0; i < wantCnt; i++ {
			ctx.insLoadNil()
		}
		return
	}
	last := len(exps) - 1
	for i := 0; i < last; i++ {
		genExp(exps[i], ctx, 1)
	}
	genExp(exps[last], ctx, wantCnt-len(exps)+1)
}

func genExp(exp ast.Exp, ctx *Context, retCnt int) {
	switch exp := exp.(type) {
	case *ast.NumberLiteralExp:
		genNumberLiteralExp(exp, ctx)
		retCnt--
	case *ast.TrueExp:
		genTrueExp(ctx)
		retCnt--
	case *ast.FalseExp:
		genFalseExp(ctx)
		retCnt--
	case *ast.StringLiteralExp:
		genStringExp(exp, ctx)
		retCnt--
	case *ast.NilExp:
		genNilExp(exp, ctx)
		retCnt--
	case *ast.MapLiteralExp:
		genMapLiteralExp(exp, ctx)
		retCnt--
	case *ast.ArrLiteralExp:
		genArrLiteralExp(exp, ctx)
		retCnt--
	case *ast.FuncLiteralExp:
		genFuncLiteralExp(exp, ctx)
		retCnt--
	case *ast.NewObjectExp:
		genNewObjectExp(exp, ctx)
		retCnt--
	case *ast.NameExp:
		genNameExp(exp, ctx)
		retCnt--
	case *ast.UnOpExp:
		genUnOpExp(exp, ctx)
		retCnt--
	case *ast.BinOpExp:
		genBinOpExp(exp, ctx)
		retCnt--
	case *ast.TernaryOpExp:
		genTernaryOpExp(exp, ctx)
		retCnt--
	case *ast.FuncCallExp:
		genFuncCallExp(exp, ctx, retCnt)
		retCnt = 0
	default:
		fmt.Printf("[%s] unknown expression %v\n", curParsingFile, exp)
		os.Exit(0)
	}
	for i := 0; i < retCnt; i++ {
		ctx.insLoadNil()
	}
}

func genNewObjectExp(exp *ast.NewObjectExp, ctx *Context) {
	ctx.writeIns(proto.INS_NEW_EMPTY_MAP)
	genExps(exp.Args, ctx, len(exp.Args))
	idx, ok := ctx.classes[exp.Name]
	if !ok {
		fmt.Printf("[%s:%d] undefined class '%s'\n", curParsingFile, exp.Line, exp.Name)
		os.Exit(0)
	}
	ctx.insLoadAnonymous(idx)
	ctx.insCall(0, byte(len(exp.Args)))
}

func genFuncLiteralExp(exp *ast.FuncLiteralExp, ctx *Context) {
	idx := uint32(len(ctx.ft.anonymousFuncs))
	ctx.ft.anonymousFuncs = append(ctx.ft.anonymousFuncs, proto.AnonymousFuncProto{
		Info: &proto.BasicInfo{
			VaArgs:     exp.VaArgs != "",
			Parameters: exp.Parameters,
		},
	})
	ctx.frame.nowParsingAnonymous = int(idx)
	ctx.insLoadAnonymous(idx)
	genFuncLiteral(&exp.FuncLiteral, ctx, idx, true)
}

func genFuncCallExp(exp *ast.FuncCallExp, ctx *Context, retCnt int) {
	genExps(exp.Args, ctx, len(exp.Args))
	genExp(exp.Func, ctx, 1)
	ctx.insCall(byte(retCnt), byte(len(exp.Args)))
}

func genTernaryOpExp(exp *ast.TernaryOpExp, ctx *Context) {
	// TODO
	genExp(exp.Exp1, ctx, 1)
	ctx.writeIns(proto.INS_JUMP_IF)
	ctx.writeUint(0) // to determine steps to jump
	old := uint32(len(ctx.frame.text))
	genExp(exp.Exp3, ctx, 1)
	ctx.writeIns(proto.INS_JUMP_REL)
	ctx.writeUint(0) // to determine steps to jump
	now := uint32(len(ctx.frame.text))
	genExp(exp.Exp2, ctx, 1)
	last := uint32(len(ctx.frame.text))

	// TODO
	binary.LittleEndian.PutUint32(ctx.frame.text[old-4:old], now-old)
	binary.LittleEndian.PutUint32(ctx.frame.text[now-4:now], last-now)
}

func genBinOpExp(exp *ast.BinOpExp, ctx *Context) {
	genExp(exp.Exp1, ctx, 1)
	switch exp.BinOp {
	case ast.BINOP_LAND:
		pos := ctx.insJumpIfLAnd(0)
		genExp(exp.Exp2, ctx, 1)
		ctx.setSteps(pos, ctx.textSize())
	case ast.BINOP_LOR:
		pos := ctx.insJumpIfLOr(0)
		genExp(exp.Exp2, ctx, 1)
		ctx.setSteps(pos, ctx.textSize())
	default:
		genExp(exp.Exp2, ctx, 1)
		ctx.writeIns(byte(exp.BinOp-ast.BINOP_START) + proto.INS_BINARY_START)
	}
}

func genUnOpExp(exp *ast.UnOpExp, ctx *Context) {
	switch exp.Op {
	case ast.UNOP_NOT:
		genExp(exp.Exp, ctx, 1)
		ctx.writeIns(proto.INS_UNARY_NOT)
	case ast.UNOP_LNOT:
		genExp(exp.Exp, ctx, 1)
		ctx.writeIns(proto.INS_UNARY_LNOT)
	case ast.UNOP_NEG:
		genExp(exp.Exp, ctx, 1)
		ctx.writeIns(proto.INS_UNARY_NEG)
	case ast.UNOP_DEC: // --i
		toAssignStmt(exp.Exp, ast.ASIGN_OP_SUBEQ, ctx)
		genExp(exp.Exp, ctx, 1)
	case ast.UNOP_INC: // ++i
		toAssignStmt(exp.Exp, ast.ASIGN_OP_ADDEQ, ctx)
		genExp(exp.Exp, ctx, 1)
	case ast.UNOP_DEC_: // i--
		genExp(exp.Exp, ctx, 1)
		toAssignStmt(exp.Exp, ast.ASIGN_OP_SUBEQ, ctx)
	case ast.UNOP_INC_: // i++
		genExp(exp.Exp, ctx, 1)
		toAssignStmt(exp.Exp, ast.ASIGN_OP_ADDEQ, ctx)
	}
}

func genNameExp(exp *ast.NameExp, ctx *Context) {
	ctx.insLoadName(exp.Name, exp.Line)
}

func genArrLiteralExp(exp *ast.ArrLiteralExp, ctx *Context) {
	for _, val := range exp.Vals {
		genExp(val, ctx, 1)
	}
	ctx.writeIns(proto.INS_SLICE_NEW)
	ctx.writeUint(uint32(len(exp.Vals)))
}

func genMapLiteralExp(exp *ast.MapLiteralExp, ctx *Context) {
	for i, key := range exp.Keys {
		ctx.insLoadConst(key)
		genExp(exp.Vals[i], ctx, 1)
	}
	ctx.writeIns(proto.INS_NEW_MAP)
	ctx.writeUint(uint32(len(exp.Keys)))
}

func genNilExp(exp *ast.NilExp, ctx *Context) {
	ctx.insLoadNil()
}

func genNumberLiteralExp(exp *ast.NumberLiteralExp, ctx *Context) {
	ctx.insLoadConst(exp.Value)
}

func genTrueExp(ctx *Context) {
	ctx.insLoadConst(true)
}

func genFalseExp(ctx *Context) {
	ctx.insLoadConst(false)
}

func genStringExp(exp *ast.StringLiteralExp, ctx *Context) {
	ctx.insLoadConst(exp.Value)
}

func toAssignStmt(exp ast.Exp, op int, ctx *Context) {
	stmt := &ast.VarAssignStmt{AssignOp: op, Rights: []ast.Exp{&ast.NumberLiteralExp{Value: int64(1)}}}
	if e, ok := exp.(*ast.NameExp); ok {
		stmt.Lefts = []ast.Var{{Prefix: e.Name}}
		genVarAssignStmt(stmt, ctx)
		return
	}
	var v ast.Var
	be := exp.(*ast.BinOpExp)
	for {
		v.Attrs = append(v.Attrs, be.Exp2)
		if e, ok := be.Exp1.(*ast.NameExp); ok {
			v.Prefix = e.Name
			break
		}
		be = be.Exp1.(*ast.BinOpExp)
	}
	for i := 0; i < len(v.Attrs)/2-1; i++ {
		mirror := len(v.Attrs) - i - 1
		v.Attrs[i], v.Attrs[mirror] = v.Attrs[mirror], v.Attrs[i]
	}
	stmt.Lefts = []ast.Var{v}
	genVarAssignStmt(stmt, ctx)
}
