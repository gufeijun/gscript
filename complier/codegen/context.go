package codegen

import (
	"encoding/binary"
	"fmt"
	"gscript/complier/ast"
	"gscript/complier/parser"
	"gscript/proto"
)

type Context struct {
	protoNum uint32
	parser   *parser.Parser
	ct       *ConstTable
	ft       *FuncTable
	classes  map[string]uint32 // class name -> FuncTable index

	frame *StackFrame
}

func newContext(parser *parser.Parser) *Context {
	ft := newFuncTable(parser.FuncDefs)
	return &Context{
		parser:  parser,
		ct:      newConstTable(),
		ft:      ft,
		classes: newClassTable(parser.ClassStmts, ft),
		frame:   newStackFrame(),
	}
}

func newClassTable(stmts []*ast.ClassStmt, ft *FuncTable) map[string]uint32 {
	classes := map[string]uint32{}
	ft.anonymousFuncs = make([]proto.AnonymousFuncProto, len(stmts))
	for i, stmt := range stmts {
		classes[stmt.Name] = uint32(i)
		info := &proto.BasicInfo{}
		__self := stmt.Constructor
		if __self != nil {
			info.Parameters = __self.Parameters
			info.VaArgs = __self.VaArgs != ""
		}
		ft.anonymousFuncs[i].Info = info
	}
	return classes
}

func (ctx *Context) pushFrame(anonymous bool, idx int) {
	frame := newStackFrame()
	frame.prev = ctx.frame
	ctx.frame = frame
}

func (ctx *Context) popFrame() *StackFrame {
	frame := ctx.frame
	ctx.frame = ctx.frame.prev
	return frame
}

func (ctx *Context) writeIns(ins byte) {
	ctx.frame.text = append(ctx.frame.text, ins)
}

func (ctx *Context) writeUint(idx uint32) {
	var arr [4]byte
	binary.LittleEndian.PutUint32(arr[:], idx)
	ctx.frame.text = append(ctx.frame.text, arr[:]...)
}

func (ctx *Context) writeByte(v byte) {
	ctx.frame.text = append(ctx.frame.text, v)
}

func (ctx *Context) insCopyName(name string) {
	ctx.writeIns(proto.INS_COPY_STACK_TOP)
	ctx.frame.nt.Set(name)
}

func (ctx *Context) insPushName(name string) {
	ctx.writeIns(proto.INS_PUSH_NAME)
	ctx.frame.nt.Set(name)
}

func (ctx *Context) insCall(wantRtnCnt byte, argCnt byte) {
	ctx.writeIns(proto.INS_CALL)
	ctx.writeByte(wantRtnCnt)
	ctx.writeByte(argCnt)
}

func (ctx *Context) insTry(catchAddr uint32) int {
	ctx.writeIns(proto.INS_TRY)
	pos := len(ctx.frame.text)
	ctx.writeUint(catchAddr)
	return pos
}

func (ctx *Context) insReturn(argCnt uint32) {
	ctx.writeIns(proto.INS_RETURN)
	ctx.writeUint(argCnt)
}

func (ctx *Context) insResizeNameTable(size uint32) int {
	ctx.writeIns(proto.INS_RESIZE_NAMETABLE)
	pos := len(ctx.frame.text)
	ctx.writeUint(size)
	return pos
}

func (ctx *Context) setAddr(pos int, addr uint32) {
	binary.LittleEndian.PutUint32(ctx.frame.text[pos:pos+4], addr)
}

func (ctx *Context) setSteps(pos int, addr uint32) {
	binary.LittleEndian.PutUint32(ctx.frame.text[pos:pos+4], addr-uint32(pos+4))
}

func (ctx *Context) insPopTop() {
	ctx.writeIns(proto.INS_POP_TOP)
}

func (ctx *Context) insJumpCase(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_CASE)
	pos := len(ctx.frame.text)
	ctx.writeUint(addr - uint32(pos+4))
	return pos
}

func (ctx *Context) insJumpIfLAnd(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_LAND)
	pos := len(ctx.frame.text)
	ctx.writeUint(addr - uint32(pos+4))
	return pos
}

func (ctx *Context) insJumpIfLOr(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_LOR)
	pos := len(ctx.frame.text)
	ctx.writeUint(addr - uint32(pos+4))
	return pos
}

func (ctx *Context) insJumpIf(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_IF)
	pos := len(ctx.frame.text)
	ctx.writeUint(addr - uint32(pos+4))
	return pos
}

func (ctx *Context) insJumpRel(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_REL)
	pos := len(ctx.frame.text)
	ctx.writeUint(addr - uint32(pos+4))
	return pos
}

func (ctx *Context) insJumpAbs(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_ABS)
	pos := len(ctx.frame.text)
	ctx.writeUint(addr)
	return pos
}

func (ctx *Context) insLoadConst(c interface{}) {
	idx := ctx.ct.Get(c)
	if stdLibGenMode {
		ctx.writeIns(proto.INS_LOAD_STD_CONST)
	} else {
		ctx.writeIns(proto.INS_LOAD_CONST)
	}
	ctx.writeUint(ctx.protoNum)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadNil() {
	ctx.writeIns(proto.INS_LOAD_NIL)
}

func (ctx *Context) insLoadFunc(idx uint32) {
	if stdLibGenMode {
		ctx.writeIns(proto.INS_LOAD_STD_FUNC)
	} else {
		ctx.writeIns(proto.INS_LOAD_FUNC)
	}
	ctx.writeUint(ctx.protoNum)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadAnonymous(idx uint32) {
	if stdLibGenMode {
		ctx.writeIns(proto.INS_LOAD_STD_ANONYMOUS)
	} else {
		ctx.writeIns(proto.INS_LOAD_ANONYMOUS)
	}
	ctx.writeUint(ctx.protoNum)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadUpValue(idx uint32) {
	ctx.writeIns(proto.INS_LOAD_UPVALUE)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadProto(idx uint32) {
	ctx.writeIns(proto.INS_LOAD_PROTO)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadStdlib(idx uint32) {
	ctx.writeIns(proto.INS_LOAD_STDLIB)
	ctx.writeUint(idx)
}

func (ctx *Context) insStoreUpValue(idx uint32) {
	ctx.writeIns(proto.INS_STORE_UPVALUE)
	ctx.writeUint(idx)
}

func searchFrame(frame *StackFrame, name string) (uint32, bool) {
	for nt := frame.nt; nt != nil; nt = nt.prev {
		if idx, ok := nt.nameTable[name]; ok {
			return idx, true
		}
	}
	return 0, false
}

func tryLoadUpValue(ctx *Context, name string) (upValueIdx uint32, ok bool) {
	upValueIdx, _, ok = ctx.frame.vt.get(name)
	if ok {
		return upValueIdx, true
	}
	var level uint32 = 0
	for frame := ctx.frame.prev; frame != nil; frame = frame.prev {
		idx, ok := searchFrame(frame, name)
		if ok {
			upValueIdx = ctx.frame.vt.set(name, level, idx)
			return upValueIdx, true
		}
		level++
	}
	return 0, false
}

func (ctx *Context) insLoadName(name string) {
	// ctx.frame.nt.nameTable

	// name is a defined variable?
	idx, ok := searchFrame(ctx.frame, name)
	if ok {
		ctx.writeIns(proto.INS_LOAD_NAME)
		ctx.writeUint(idx)
		return
	}

	// name is an upValue?
	if idx, ok := tryLoadUpValue(ctx, name); ok {
		ctx.insLoadUpValue(idx)
		return
	}

	// name is a function?
	idx, ok = ctx.ft.funcMap[name]
	if ok {
		ctx.insLoadFunc(idx)
		return
	}

	// name is a enum constant?
	idx, ok = ctx.ct.getEnum(name)
	if ok {
		ctx.writeIns(proto.INS_LOAD_CONST)
		ctx.writeUint(ctx.protoNum)
		ctx.writeUint(idx)
		return
	}

	// name is a builtin function?
	idx, ok = builtinFuncs[name]
	if ok {
		ctx.writeIns(proto.INS_LOAD_BUILTIN)
		ctx.writeUint(idx)
		return
	}

	panic(fmt.Sprintf("undefined name %s", name)) // TODO
}

func (ctx *Context) insStoreName(name string) {
	idx, ok := ctx.frame.nt.get(name)
	if ok {
		ctx.writeIns(proto.INS_STORE_NAME)
		ctx.writeUint(idx)
		return
	}

	// if name is an upvalue?
	if idx, ok := tryLoadUpValue(ctx, name); ok {
		ctx.insStoreUpValue(idx)
		return
	}

	panic(fmt.Sprintf("name %s do not exist", name))
}

func (ctx *Context) enterBlock() {
	nt := newNameTable(ctx.frame.nt.nameIdx)
	nt.prev = ctx.frame.nt
	ctx.frame.nt = nt
}

func (ctx *Context) leaveBlock(size uint32, varDecl bool) {
	nt := ctx.frame.nt
	if nt.prev == nil {
		panic("") // TODO
	}
	ctx.frame.nt = nt.prev
	*ctx.frame.nt.nameIdx = size
	if varDecl {
		ctx.insResizeNameTable(size)
	}
}

func (ctx *Context) textSize() uint32 {
	return uint32(len(ctx.frame.text))
}

type unhandledGoto struct {
	label     string
	resizePos int
	jumpPos   int
}

type label struct {
	name          string
	addr          uint32
	nameTableSize uint32
}

var builtinFuncs = map[string]uint32{
	"print":                0,
	"len":                  1,
	"append":               2,
	"sub":                  3,
	"type":                 4,
	"delete":               5,
	"clone":                6,
	"__buffer_new":         7,
	"__buffer_readNumber":  8,
	"__buffer_writeNumber": 9,
	"__buffer_toString":    10,
	"__buffer_slice":       11,
	"__buffer_concat":      12,
	"__buffer_copy":        13,
	"__buffer_from":        14,
	"__open":               15,
	"__read":               16,
	"__write":              17,
	"__close":              18,
	"__seek":               19,
	"__remove":             20,
	"__fchmod":             21,
	"__chmod":              22,
	"__fchown":             23,
	"__chown":              24,
	"__fchdir":             25,
	"__chdir":              26,
	"__fstat":              27,
	"__stat":               28,
	"__rename":             29,
	"__mkdir":              30,
	"__exit":               31,
	"__getenv":             32,
	"__setenv":             33,
	"__readdir":            34,
	"__freaddir":           35,
	"throw":                36,
}
