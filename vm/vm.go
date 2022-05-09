package vm

import (
	"encoding/binary"
	"gscript/proto"
	"unsafe"
)

type VM struct {
	text       []proto.Instruction
	constTable []interface{}
	frame      *funcFrame
	funcTable  []proto.Func
	stack      *evalStack
	pc         uint32

	stopped bool
}

func NewVM(Text []proto.Instruction, ct []interface{}, ft []proto.Func) *VM {
	return &VM{
		text:       Text,
		constTable: ct,
		stack:      newEvalStack(),
		frame:      newFuncFrame(),
		funcTable:  ft,
		pc:         0,
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
