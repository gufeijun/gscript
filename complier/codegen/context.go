package codegen

import (
	"encoding/binary"
	"fmt"
	"gscript/complier/parser"
	"gscript/proto"
)

type Context struct {
	parser      *parser.Parser
	buf         []byte
	ct          *ConstTable
	nt          *NameTable
	ft          *FuncTable
	bs          *blockStack
	validLabels map[string]label

	returnAtEnd bool
}

func newContext(parser *parser.Parser) *Context {
	return &Context{
		parser:      parser,
		ct:          newConstTable(),
		nt:          NewNameTable(),
		ft:          newFuncTable(parser.FuncDefs),
		bs:          newBlockStack(),
		validLabels: make(map[string]label),
	}
}

func (ctx *Context) renew() {
	ctx.returnAtEnd = false
	ctx.nt = NewNameTable()
	ctx.bs = newBlockStack()
}

func (ctx *Context) writeIns(ins byte) {
	ctx.buf = append(ctx.buf, ins)
}

func (ctx *Context) writeUint(idx uint32) {
	var arr [4]byte
	binary.LittleEndian.PutUint32(arr[:], idx)
	ctx.buf = append(ctx.buf, arr[:]...)
}

func (ctx *Context) writeByte(v byte) {
	ctx.buf = append(ctx.buf, v)
}

func (ctx *Context) insPushName(name string) {
	ctx.writeIns(proto.INS_PUSH_NAME)
	ctx.nt.Set(name)
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
	pos := len(ctx.buf)
	ctx.writeUint(size)
	return pos
}

func (ctx *Context) setAddr(pos int, addr uint32) {
	binary.LittleEndian.PutUint32(ctx.buf[pos:pos+4], addr)
}

func (ctx *Context) insPopTop() {
	ctx.writeIns(proto.INS_POP_TOP)
}

func (ctx *Context) insJumpCase(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_CASE)
	pos := len(ctx.buf)
	ctx.writeUint(addr)
	return pos
}

func (ctx *Context) insJumpIf(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_IF)
	pos := len(ctx.buf)
	ctx.writeUint(addr)
	return pos
}

func (ctx *Context) insJumpRel(steps uint32) int {
	ctx.writeIns(proto.INS_JUMP_REL)
	pos := len(ctx.buf)
	ctx.writeUint(steps)
	return pos
}

func (ctx *Context) insJumpAbs(addr uint32) int {
	ctx.writeIns(proto.INS_JUMP_ABS)
	pos := len(ctx.buf)
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

func (ctx *Context) insLoadName(name string) {
	for nt := ctx.nt; nt != nil; nt = nt.prev {
		if idx, ok := nt.nameTable[name]; ok {
			ctx.writeIns(proto.INS_LOAD_NAME)
			ctx.writeUint(idx)
			return
		}
	}

	// if name is a function?
	idx, ok := ctx.ft.funcMap[name]
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
	ctx.writeIns(proto.INS_STORE_NAME)

	idx, ok := ctx.nt.get(name)
	if !ok {
		panic(fmt.Sprintf("name %s do not exist", name))
	}
	ctx.writeUint(idx)
}

func (ctx *Context) enterBlock() {
	nt := newNameTable(ctx.nt.nameIdx)
	nt.prev = ctx.nt
	ctx.nt = nt
}

func (ctx *Context) leaveBlock(size uint32, varDecl bool) {
	nt := ctx.nt
	if nt.prev == nil {
		panic("") // TODO
	}
	ctx.nt = nt.prev
	*ctx.nt.nameIdx = size
	if varDecl {
		ctx.insResizeNameTable(size)
	}
}

func (ctx *Context) textSize() uint32 {
	return uint32(len(ctx.buf))
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
