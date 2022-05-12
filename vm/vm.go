package vm

import (
	"encoding/binary"
	"gscript/proto"
	"unsafe"
)

type VM struct {
	stopped        bool
	topFrame       *stackFrame
	frame          *stackFrame
	stack          *evalStack
	constTable     []interface{}
	funcTable      []proto.FuncProto
	anonymousTable []proto.AnonymousFuncProto
}

func NewVM(_proto *proto.Proto) *VM {
	topFrame := newFuncFrame()
	topFrame.text = _proto.Text
	return &VM{
		topFrame:       topFrame,
		frame:          topFrame,
		stack:          newEvalStack(),
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
		instruction := vm.frame.text[vm.frame.pc]
		vm.frame.pc++
		Execute(vm, instruction)
	}
}

func (vm *VM) Stop() {
	vm.stopped = true
}

func (vm *VM) getOpNum() uint32 {
	arr := *(*[4]byte)(unsafe.Pointer(&vm.frame.text[vm.frame.pc]))
	v := binary.LittleEndian.Uint32(arr[:])
	vm.frame.pc += 4
	return v
}
