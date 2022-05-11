package codegen

import (
	"encoding/binary"
	"fmt"
	"gscript/complier/parser"
	"gscript/proto"
)

type Context struct {
	parser *parser.Parser
	ct     *ConstTable
	ft     *FuncTable
	frame  *StackFrame
}

func newContext(parser *parser.Parser) *Context {
	return &Context{
		parser: parser,
		ct:     newConstTable(),
		ft:     newFuncTable(parser.FuncDefs),
		frame:  newStackFrame(),
	}
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

func (ctx *Context) insPushName(name string) {
	ctx.writeIns(proto.INS_PUSH_NAME)
	ctx.frame.nt.Set(name)
}

func (ctx *Context) insCall(wantRtnCnt byte, argCnt byte) {
	ctx.writeIns(proto.INS_CALL)
	ctx.writeByte(wantRtnCnt)
	ctx.writeByte(argCnt)
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
	ctx.writeIns(proto.INS_LOAD_CONST)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadFunc(idx uint32) {
	ctx.writeIns(proto.INS_LOAD_FUNC)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadAnonymous(idx uint32) {
	ctx.writeIns(proto.INS_LOAD_ANONYMOUS)
	ctx.writeUint(idx)
}

func (ctx *Context) insLoadUpValue(idx uint32) {
	ctx.writeIns(proto.INS_LOAD_UPVALUE)
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
	idx, ok := searchFrame(ctx.frame, name)
	if ok {
		ctx.writeIns(proto.INS_LOAD_NAME)
		ctx.writeUint(idx)
		return
	}

	// // if name is an upValue?
	if idx, ok := tryLoadUpValue(ctx, name); ok {
		ctx.insLoadUpValue(idx)
		return
	}

	// if name is a function?
	idx, ok = ctx.ft.funcMap[name]
	if ok {
		ctx.insLoadFunc(idx)
		return
	}

	// if name is a enum constant?
	idx, ok = ctx.ct.getEnum(name)
	if ok {
		ctx.writeIns(proto.INS_LOAD_CONST)
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
