package vm

import (
	"encoding/binary"
	"fmt"
	"gscript/proto"
	"os"
	"unsafe"
)

type VM struct {
	stopped           bool
	protos            []proto.Proto
	stdlibs           []proto.Proto
	curProto          *protoFrame
	builtinFuncFailed bool
	curCallingBuiltin string
}

func NewVM(protos []proto.Proto, stdlibs []proto.Proto) *VM {
	return &VM{
		protos:   protos,
		stdlibs:  stdlibs,
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

func (vm *VM) exit(format string, args ...interface{}) {
	fmt.Printf("[%s] ", vm.curProto.filepath)
	exit(format, args...)
}

func exit(format string, args ...interface{}) {
	fmt.Printf("runtime error: ")
	fmt.Printf(format, args...)
	os.Exit(0)
}

func (vm *VM) assert(cond bool) {
	if cond {
		return
	}
	fmt.Printf("[%s] runtime error: ", vm.curProto.filepath)
	fmt.Printf("call builtin function '%s' failed, ", vm.curCallingBuiltin)
	fmt.Println("please check count and type of arguments")
	os.Exit(0)
}
