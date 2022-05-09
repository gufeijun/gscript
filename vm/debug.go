package vm

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"gscript/proto"
	"os"
	"strings"
	"unsafe"
)

var cmds = map[string]func(vm *VM, args []string){
	"help":  debugHelp,
	"n":     debugNext,
	"next":  debugNext,
	"v":     debugShowVar,
	"var":   debugShowVar,
	"s":     debugShowStack,
	"stack": debugShowStack,
	"c":     debugShowCode,
	"code":  debugShowCode,
	"const": debugShowConstant,
	"r":     debugRun,
}

func debugShowConstant(vm *VM, args []string) {
	fmt.Printf("constants: ")
	for i, value := range vm.constTable {
		showValue(value)
		if i != len(vm.frame.symbolTable.values)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Println()
}

func debugRun(vm *VM, args []string) {
	for {
		instruction := vm.text[vm.pc]
		if instruction == proto.Instruction(proto.INS_STOP) {
			break
		}
		vm.pc++
		Execute(vm, instruction)
	}
}

func debugShowCode(vm *VM, args []string) {
	cnt := 15
	start := vm.pc
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &start)
	}
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &cnt)
	}
	showCode(vm, start, cnt)
}

func showCode(vm *VM, pc uint32, cnt int) {
	for i := 0; i < cnt; i++ {
		pc += showInstruction(vm, pc)
		if int(pc) >= len(vm.text) {
			break
		}
	}
}

func debugShowStack(vm *VM, args []string) {
	fmt.Printf("stack: ")
	for i := 0; i < len(vm.stack.Buf); i++ {
		val := vm.stack.Buf[i]
		showValue(val)
		fmt.Printf(" ")
	}
	fmt.Println()
}

func debugShowVar(vm *VM, args []string) {
	fmt.Printf("variables: ")
	for i, value := range vm.frame.symbolTable.values {
		showValue(value)
		if i != len(vm.frame.symbolTable.values)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Println()
}

func debugNext(vm *VM, args []string) {
	instruction := vm.text[vm.pc]
	vm.pc++
	Execute(vm, instruction)
}

func debugHelp(vm *VM, args []string) {
	// TODO
	fmt.Println()
}

func (vm *VM) Debug() {
	for {
		if vm.stopped {
			fmt.Println("done!")
			break
		}
		bufr := bufio.NewReader(os.Stdin)

		line, _, _ := bufr.ReadLine()
		args := strings.Split(string(line), " ")
		if len(args) == 0 {
			continue
		}
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}
		handler, ok := cmds[args[0]]
		if !ok {
			continue
		}
		handler(vm, args[1:])
		fmt.Println()
	}
}

func showValue(val interface{}) {
	if str, ok := val.(string); ok {
		fmt.Printf("\"%s\"", str)
	} else {
		fmt.Printf("%v", val)
	}
}

func showInstruction(vm *VM, pc uint32) uint32 {
	skip := 1
	ins := byte(vm.text[pc])
	if pc == vm.pc {
		fmt.Printf("->")
	}
	fmt.Printf("%d	\t", pc)
	switch ins {
	case proto.INS_UNARY_NOT:
		fmt.Printf("NOT")
	case proto.INS_UNARY_NEG:
		fmt.Printf("NEG")
	case proto.INS_UNARY_LNOT:
		fmt.Printf("LNOT")
	case proto.INS_BINARY_ADD:
		fmt.Printf("ADD")
	case proto.INS_BINARY_SUB:
		fmt.Printf("SUB")
	case proto.INS_BINARY_MUL:
		fmt.Printf("MUL")
	case proto.INS_BINARY_DIV:
		fmt.Printf("DIV")
	case proto.INS_BINARY_MOD:
		fmt.Printf("MOD")
	case proto.INS_BINARY_AND:
		fmt.Printf("AND")
	case proto.INS_BINARY_XOR:
		fmt.Printf("XOR")
	case proto.INS_BINARY_OR:
		fmt.Printf("OR")
	case proto.INS_BINARY_IDIV:
		fmt.Printf("IDIV")
	case proto.INS_BINARY_SHR:
		fmt.Printf("SHR")
	case proto.INS_BINARY_SHL:
		fmt.Printf("SHL")
	case proto.INS_BINARY_LE:
		fmt.Printf("LE")
	case proto.INS_BINARY_GE:
		fmt.Printf("GE")
	case proto.INS_BINARY_LT:
		fmt.Printf("LT")
	case proto.INS_BINARY_GT:
		fmt.Printf("GT")
	case proto.INS_BINARY_EQ:
		fmt.Printf("EQ")
	case proto.INS_BINARY_NE:
		fmt.Printf("NE")
	case proto.INS_BINARY_LAND:
		fmt.Printf("LAND")
	case proto.INS_BINARY_LOR:
		fmt.Printf("LOR")
	case proto.INS_BINARY_ATTR:
		fmt.Printf("ATTR")
	case proto.INS_LOAD_CONST:
		fmt.Printf("LOAD_CONST %v", vm.constTable[getOpNum(vm, pc)])
		skip += 4
	case proto.INS_LOAD_NAME:
		fmt.Printf("LOAD_NAME %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_LOAD_FUNC:
		fmt.Printf("LOAD_FUNC %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_LOAD_VALUE:
		fmt.Printf("LOAD_VALUE %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_STORE_NAME:
		fmt.Printf("STORE_NAME %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_PUSH_NAME_NIL:
		fmt.Printf("PUSH_NAME_NIL")
	case proto.INS_CALL:
		pc++
		wantRtnCnt := byte(vm.text[pc])
		pc++
		argCnt := byte(vm.text[pc])
		fmt.Printf("CALL %d %d", wantRtnCnt, argCnt)
		skip += 2
	case proto.INS_RETURN:
		fmt.Printf("RETURN with %d values", getOpNum(vm, pc))
		skip += 4
	case proto.INS_PUSH_NAME:
		fmt.Printf("PUSH_NAME")
	case proto.INS_RESIZE_NAMETABLE:
		fmt.Printf("RESIZE_NAMETABLE %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_POP_TOP:
		fmt.Printf("POP_TOP")
	case proto.INS_STOP:
		fmt.Printf("STOP")
	case proto.INS_SLICE_NEW:
		fmt.Printf("SLICE_NEW %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_SLICE_APPEND:
		fmt.Printf("SLICE_APPEND")
	case proto.INS_MAP_NEW:
		fmt.Printf("MAP_NEW %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_ATTR_ASSIGN:
		fmt.Printf("ATTR_ASSIGN")
	case proto.INS_ATTR_ASSIGN_ADDEQ:
		fmt.Printf("ATTR_ASSIGN_ADDEQ")
	case proto.INS_ATTR_ASSIGN_SUBEQ:
		fmt.Printf("ATTR_ASSIGN_SUBEQ")
	case proto.INS_ATTR_ASSIGN_MULEQ:
		fmt.Printf("ATTR_ASSIGN_MULEQ")
	case proto.INS_ATTR_ASSIGN_DIVEQ:
		fmt.Printf("ATTR_ASSIGN_DIVEQ")
	case proto.INS_ATTR_ASSIGN_MODEQ:
		fmt.Printf("ATTR_ASSIGN_MODEQ")
	case proto.INS_ATTR_ASSIGN_ANDEQ:
		fmt.Printf("ATTR_ASSIGN_ANDEQ")
	case proto.INS_ATTR_ASSIGN_XOREQ:
		fmt.Printf("ATTR_ASSIGN_XOREQ")
	case proto.INS_ATTR_ASSIGN_OREQ:
		fmt.Printf("ATTR_ASSIGN_OREQ")
	case proto.INS_ATTR_ACCESS:
		fmt.Printf("ATTR_ACCESS")
	case proto.INS_JUMP_REL:
		fmt.Printf("JUMP_REL %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_JUMP_ABS:
		fmt.Printf("JUMP_ABS %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_JUMP_IF:
		fmt.Printf("JUMP_IF %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_JUMP_CASE:
		fmt.Printf("JUMP_CASE %d", getOpNum(vm, pc))
		skip += 4
	case proto.INS_ROT_TWO:
		fmt.Printf("ROT_TWO")
	}
	fmt.Println()
	return uint32(skip)
}

func getOpNum(vm *VM, pc uint32) uint32 {
	pc++
	_data := vm.text[pc : pc+4]
	data := *(*[]byte)(unsafe.Pointer((uintptr(unsafe.Pointer(&_data)))))
	return binary.LittleEndian.Uint32(data)
}
