package codegen

import (
	"encoding/binary"
	"fmt"
	"gscript/complier/ast"
	"gscript/proto"
)

func genExps(exps []ast.Exp, ctx *Context, wantCnt int) {
	if len(exps) == 0 {
		for i := 0; i < wantCnt; i++ {
			ctx.insLoadConst(nil)
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
		// TODO
	case *ast.NewObjectExp:
		// TODO
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
		panic(fmt.Sprintf("do not support exp: %v", exp))
	}
	for i := 0; i < retCnt; i++ {
		ctx.insLoadConst(nil)
	}
}

func genFuncCallExp(exp *ast.FuncCallExp, ctx *Context, retCnt int) {
	genExps(exp.Args, ctx, len(exp.Args))
	genExp(exp.Func, ctx, 1)
	ctx.insCall(byte(retCnt), byte(len(exp.Args)))
}

func genTernaryOpExp(exp *ast.TernaryOpExp, ctx *Context) {
	genExp(exp.Exp1, ctx, 1)
	ctx.writeIns(proto.INS_JUMP_IF)
	ctx.writeUint(0) // to determine steps to jump
	old := uint32(len(ctx.buf))
	genExp(exp.Exp3, ctx, 1)
	ctx.writeIns(proto.INS_JUMP_REL)
	ctx.writeUint(0) // to determine steps to jump
	now := uint32(len(ctx.buf))
	genExp(exp.Exp2, ctx, 1)
	last := uint32(len(ctx.buf))

	binary.LittleEndian.PutUint32(ctx.buf[old-4:old], now-old)
	binary.LittleEndian.PutUint32(ctx.buf[now-4:now], last-now)
}

func genBinOpExp(exp *ast.BinOpExp, ctx *Context) {
	genExp(exp.Exp1, ctx, 1)
	genExp(exp.Exp2, ctx, 1)
	ctx.writeIns(byte(exp.BinOp-ast.BINOP_START) + proto.INS_BINARY_START)
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
	ctx.insLoadName(exp.Name)
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
	ctx.writeIns(proto.INS_MAP_NEW)
	ctx.writeUint(uint32(len(exp.Keys)))
}

func genNilExp(exp *ast.NilExp, ctx *Context) {
	ctx.insLoadConst(nil)
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
	be, ok := exp.(*ast.BinOpExp)
	if !ok {
		panic("") // TODO
	}
	for {
		v.Attrs = append(v.Attrs, be.Exp2)
		if e, ok := be.Exp1.(*ast.NameExp); ok {
			v.Prefix = e.Name
			break
		}
		if e, ok := be.Exp1.(*ast.BinOpExp); ok {
			be = e
			continue
		}
		panic("") // TODO
	}
	for i := 0; i < len(v.Attrs)/2-1; i++ {
		mirror := len(v.Attrs) - i - 1
		v.Attrs[i], v.Attrs[mirror] = v.Attrs[mirror], v.Attrs[i]
	}
	stmt.Lefts = []ast.Var{v}
	genVarAssignStmt(stmt, ctx)
}
