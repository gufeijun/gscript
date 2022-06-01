package vm

import (
	"encoding/binary"
	"gscript/proto"
	"unsafe"
)

type VM struct {
	stopped           bool
	protos            []proto.Proto
	curProto          *protoFrame
	builtinFuncFailed bool
}

func NewVM(protos []proto.Proto) *VM {
	return &VM{
		protos:   protos,
		curProto: newProtoFrame(protos[0]),
	}
}

func (vm *VM) Run() {
	for {
		if vm.stopped {
			break
		}
		instruction := vm.curProto.frame.text[vm.curProto.frame.pc]
		vm.curProto.frame.pc++
		Execute(vm, instruction)
	}
}

func (vm *VM) Stop() {
	vm.stopped = true
}

func (vm *VM) getOpNum() uint32 {
	arr := *(*[4]byte)(unsafe.Pointer(&vm.curProto.frame.text[vm.curProto.frame.pc]))
	v := binary.LittleEndian.Uint32(arr[:])
	vm.curProto.frame.pc += 4
	return v
}
