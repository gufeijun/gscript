package vm

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"gscript/proto"
	"gscript/vm/types"
	"os"
	"strconv"
	"strings"
)

var cmds = map[string]func(vm *VM, args []string){
	"help":    debugHelp,
	"n":       debugNext,
	"next":    debugNext,
	"v":       debugShowVar,
	"var":     debugShowVar,
	"s":       debugShowStack,
	"stack":   debugShowStack,
	"c":       debugShowCode,
	"code":    debugShowCode,
	"const":   debugShowConstant,
	"r":       debugRun,
	"f":       debugShowFunc,
	"ff":      debugShowAnonymousFunc,
	"upvalue": debugShowUpValue,
}

func debugShowUpValue(vm *VM, args []string) {
	fmt.Printf("upValues: ")
	for i, v := range vm.curProto.frame.upValues {
		fmt.Printf("%v", v.Value)
		if i != len(vm.curProto.frame.upValues) {
			fmt.Printf(", ")
		}
	}
	fmt.Println()
}

func showFunc(upValues []*types.GsValue) {
	fmt.Printf("upvalues: [")
	for i, upValue := range upValues {
		showValue(upValue.Value)
		if i != len(upValues)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("]\n")
}

func debugShowAnonymousFunc(vm *VM, args []string) {
	if len(args) == 0 {
		return
	}
	num, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}
	f := vm.curProto.anonymousTable[num]
	cnt := -1
	if len(args) > 1 {
		fmt.Scanf("%d", &cnt)
	}
	showCode(vm, f.Info.Text, 0, cnt)
}

func debugShowFunc(vm *VM, args []string) {
	if len(args) != 0 {
		num, err := strconv.Atoi(args[0])
		if err != nil || num >= len(vm.curProto.funcTable) {
			return
		}
		f := vm.curProto.funcTable[num]
		upValues, _ := f.UpValueTable.([]*types.GsValue)
		showFunc(upValues)
		cnt := -1
		if len(args) > 1 {
			fmt.Scanf("%d", &cnt)
		}
		showCode(vm, f.Info.Text, 0, cnt)
		return
	}
	for i, f := range vm.curProto.funcTable {
		fmt.Printf("%dth: ", i)
		upValues, _ := f.UpValueTable.([]*types.GsValue)
		showFunc(upValues)
	}
	fmt.Println()
}

func debugShowConstant(vm *VM, args []string) {
	fmt.Printf("constants: ")
	if len(args) == 0 {
		return
	}
	var num int
	fmt.Sscanf(args[0], "%d", &num)
	for i, value := range vm.protos[num].Consts {
		showValue(value)
		if i != len(vm.curProto.frame.symbolTable.values)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Println()
}

func debugRun(vm *VM, args []string) {
	for {
		if vm.stopped {
			return
		}
		instruction := vm.curProto.frame.text[vm.curProto.frame.pc]
		if instruction == proto.INS_STOP {
			break
		}
		vm.curProto.frame.pc++
		Execute(vm, instruction)
	}
}

func debugShowCode(vm *VM, args []string) {
	cnt := 15
	start := vm.curProto.frame.pc
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &start)
	}
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &cnt)
	}
	showCode(vm, vm.curProto.frame.text, start, cnt)
}

func showCode(vm *VM, text []byte, pc uint32, cnt int) {
	for i := 0; ; i++ {
		if int(pc) >= len(text) || i == cnt {
			break
		}
		pc += showInstruction(vm, text, pc)
	}
}

func debugShowStack(vm *VM, args []string) {
	fmt.Printf("stack: ")
	for i := 0; i < len(vm.curProto.stack.Buf); i++ {
		val := vm.curProto.stack.Buf[i]
		showValue(val)
		fmt.Printf(" ")
	}
	fmt.Println()
}

func debugShowVar(vm *VM, args []string) {
	fmt.Printf("variables: ")
	for i, val := range vm.curProto.frame.symbolTable.values {
		showValue(val.Value)
		if i != len(vm.curProto.frame.symbolTable.values)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Println()
}

func debugNext(vm *VM, args []string) {
	instruction := vm.curProto.frame.text[vm.curProto.frame.pc]
	vm.curProto.frame.pc++
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
	switch val := val.(type) {
	case string:
		fmt.Printf("\"%s\"", val)
	case *types.Closure:
		fmt.Printf("closure{")
		fmt.Printf("upvalues: [")
		for i, upValue := range val.UpValues {
			fmt.Printf("%v", *upValue)
			if i != len(val.UpValues)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf("]}")
	case *builtinFunc:
		fmt.Printf("builtin(\"%s\")", val.name)
	case *types.Array:
		fmt.Printf("%v", val.Data)
	case *types.Object:
		fmt.Printf("%v", val.Data)
	case *types.Buffer:
		fmt.Printf("Buffer")
	case *types.File:
		fmt.Printf("File")
	default:
		fmt.Printf("%v", val)
	}

}

func showInstruction(vm *VM, text []byte, pc uint32) uint32 {
	skip := 1
	ins := byte(text[pc])
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
	case proto.INS_LOAD_NIL:
		fmt.Printf("LOAD_NIL")
	case proto.INS_LOAD_CONST:
		pc++
		protoNum := getOpNum(text, pc)
		pc += 4
		constNum := getOpNum(text, pc)
		fmt.Printf("LOAD_CONST %v", vm.protos[protoNum].Consts[constNum])
		skip += 8
	case proto.INS_LOAD_STD_CONST:
		pc++
		protoNum := getOpNum(text, pc)
		pc += 4
		constNum := getOpNum(text, pc)
		fmt.Printf("LOAD_STD_CONST %v", vm.stdlibs[protoNum].Consts[constNum])
		skip += 8
	case proto.INS_LOAD_NAME:
		pc++
		fmt.Printf("LOAD_NAME %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_LOAD_FUNC:
		pc++
		protoNum := getOpNum(text, pc)
		pc += 4
		fmt.Printf("LOAD_FUNC %d %d", protoNum, getOpNum(text, pc))
		skip += 8
	case proto.INS_LOAD_STD_FUNC:
		pc++
		protoNum := getOpNum(text, pc)
		pc += 4
		fmt.Printf("LOAD_STD_FUNC %d %d", protoNum, getOpNum(text, pc))
		skip += 8
	case proto.INS_LOAD_ANONYMOUS:
		pc++
		protoNum := getOpNum(text, pc)
		pc += 4
		fmt.Printf("LOAD_ANONYMOUS %d %d", protoNum, getOpNum(text, pc))
		skip += 8
	case proto.INS_LOAD_STD_ANONYMOUS:
		pc++
		protoNum := getOpNum(text, pc)
		pc += 4
		fmt.Printf("LOAD_STD_ANONYMOUS %d %d", protoNum, getOpNum(text, pc))
		skip += 8
	case proto.INS_LOAD_UPVALUE:
		pc++
		fmt.Printf("LOAD_UPVALUE %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_LOAD_PROTO:
		pc++
		fmt.Printf("LOAD_PROTO %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_LOAD_STDLIB:
		pc++
		fmt.Printf("LOAD_STDLIB %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_STORE_NAME:
		pc++
		fmt.Printf("STORE_NAME %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_STORE_UPVALUE:
		pc++
		fmt.Printf("STORE_UPVALUE %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_LOAD_BUILTIN:
		pc++
		fmt.Printf("LOAD_BUILTIN \"%s\"", builtinFuncs[getOpNum(text, pc)].name)
		skip += 4
	case proto.INS_STORE_KV:
		fmt.Printf("STORE_KV")
	case proto.INS_PUSH_NAME_NIL:
		fmt.Printf("PUSH_NAME_NIL")
	case proto.INS_CALL:
		pc++
		wantRtnCnt := byte(vm.curProto.frame.text[pc])
		pc++
		argCnt := byte(vm.curProto.frame.text[pc])
		fmt.Printf("CALL %d %d", wantRtnCnt, argCnt)
		skip += 2
	case proto.INS_RETURN:
		pc++
		fmt.Printf("RETURN with %d values", getOpNum(text, pc))
		skip += 4
	case proto.INS_PUSH_NAME:
		fmt.Printf("PUSH_NAME")
	case proto.INS_COPY_NAME:
		fmt.Printf("COPY_NAME")
	case proto.INS_RESIZE_NAMETABLE:
		pc++
		fmt.Printf("RESIZE_NAMETABLE %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_POP_TOP:
		fmt.Printf("POP_TOP")
	case proto.INS_STOP:
		fmt.Printf("STOP")
	case proto.INS_SLICE_NEW:
		pc++
		fmt.Printf("SLICE_NEW %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_NEW_EMPTY_MAP:
		fmt.Printf("NEW_EMPTY_MAP")
	case proto.INS_NEW_MAP:
		pc++
		fmt.Printf("MAP_NEW %d", getOpNum(text, pc))
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
		pc++
		fmt.Printf("JUMP_REL %d", int32(getOpNum(text, pc)))
		skip += 4
	case proto.INS_JUMP_ABS:
		pc++
		fmt.Printf("JUMP_ABS %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_JUMP_IF:
		pc++
		fmt.Printf("JUMP_IF %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_JUMP_LAND:
		pc++
		fmt.Printf("JUMP_LAND %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_JUMP_LOR:
		pc++
		fmt.Printf("JUMP_LOR %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_JUMP_CASE:
		pc++
		fmt.Printf("JUMP_CASE %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_ROT_TWO:
		fmt.Printf("ROT_TWO")
	case proto.INS_EXPORT:
		fmt.Printf("EXPORT")
	case proto.INS_TRY:
		pc++
		fmt.Printf("TRY %d", getOpNum(text, pc))
		skip += 4
	case proto.INS_END_TRY:
		fmt.Printf("END_TRY")
	}
	fmt.Println()
	return uint32(skip)
}

func getOpNum(text []byte, pc uint32) uint32 {
	data := text[pc : pc+4]
	return binary.LittleEndian.Uint32(data)
}
