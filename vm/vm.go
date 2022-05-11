package vm

import (
	"encoding/binary"
	"gscript/proto"
	"unsafe"
)

type VM struct {
	stopped        bool
	pc             uint32
	topFrame       *stackFrame
	frame          *stackFrame
	stack          *evalStack
	text           []proto.Instruction
	constTable     []interface{}
	funcTable      []proto.FuncProto
	anonymousTable []proto.AnonymousFuncProto
}

func NewVM(_proto *proto.Proto) *VM {
	topFrame := newFuncFrame()
	return &VM{
		topFrame:       topFrame,
		frame:          topFrame,
		stack:          newEvalStack(),
		text:           _proto.Text,
		constTable:     _proto.Consts,
		funcTable:      _proto.Funcs,
		anonymousTable: _proto.AnonymousFuncs,
	}
}

func (vm *VM) Run() {
	for {
		if vm.stopped {
			break
		}
		instruction := vm.text[vm.pc]
		vm.pc++
		Execute(vm, instruction)
	}
}

func (vm *VM) Stop() {
	vm.stopped = true
}

func (vm *VM) getOpNum() uint32 {
	arr := *(*[4]byte)(unsafe.Pointer(&vm.text[vm.pc]))
	v := binary.LittleEndian.Uint32(arr[:])
	vm.pc += 4
	return v
}
