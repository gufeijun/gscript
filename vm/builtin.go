package vm

import (
	"encoding/binary"
	"fmt"
	"math"
)

type builtinFunc struct {
	handler func(argCnt int, vm *VM) int
	name    string
}

var builtinFuncs = []builtinFunc{
	{builtinPrint, "print"},
	{builtinLen, "len"},
	{builtinAppend, "append"},
	{builtinSub, "sub"},
	{builtinType, "type"},
	{builtinDelete, "delete"},
	{builtinClone, "clone"},
	{builtinForeach, "foreach"},
	{builtinBufferNew, "__buffer_new"},
	{builtinBufferReadNumber, "__buffer_readNumber"},
	{builtinBufferWriteNumber, "__buffer_writeNumber"},
	{builtinBufferToString, "__buffer_toString"},
}

// arg1: Buffer, arg2: offset, arg3: length
func builtinBufferToString(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "")
	length, ok := vm.curProto.stack.pop().(int64)
	assertS(ok, "")
	offset, ok := vm.curProto.stack.pop().(int64)
	assertS(ok, "")
	buffer, ok := vm.curProto.stack.pop().([]byte)
	assertS(ok, "")
	vm.curProto.stack.Push(string(buffer[offset : offset+length]))
	return 1
}

// arg1: Buffer, arg2: offset, arg3: size, arg4: littleEndian arg5: number
func builtinBufferWriteNumber(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 5, "") // TODO
	number := vm.curProto.stack.pop()
	littleEndian, ok := vm.curProto.stack.pop().(bool)
	assertS(ok, "")
	size, ok := vm.curProto.stack.pop().(int64)
	assertS(ok, "")
	offset, ok := vm.curProto.stack.pop().(int64)
	assertS(ok, "") // TODO
	buffer, ok := vm.curProto.stack.pop().([]byte)
	assertS(ok, "") // TODO

	switch size {
	case 1:
		v := byte(number.(int64))
		buffer[offset] = v
	case 2:
		v := uint16(number.(int64))
		if littleEndian {
			binary.LittleEndian.PutUint16(buffer[offset:], v)
		} else {
			binary.BigEndian.PutUint16(buffer[offset:], v)
		}
	case 4:
		var v uint32
		if vf, ok := number.(float64); ok {
			v = math.Float32bits(float32(vf))
		} else {
			v = uint32(number.(int64))
		}
		if littleEndian {
			binary.LittleEndian.PutUint32(buffer[offset:], v)
		} else {
			binary.BigEndian.PutUint32(buffer[offset:], v)
		}
	case 8:
		var v uint64
		if vf, ok := number.(float64); ok {
			v = math.Float64bits(vf)
		} else {
			v = uint64(number.(int64))
		}
		if littleEndian {
			binary.LittleEndian.PutUint64(buffer[offset:], v)
		} else {
			binary.BigEndian.PutUint64(buffer[offset:], v)
		}
	}
	return 0
}

// arg1: Buffer, arg2: offset, arg3: size, arg4: signed arg5: littleEndian arg6: isFloat
func builtinBufferReadNumber(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 6, "") // TODO
	isFloat, ok := vm.curProto.stack.pop().(bool)
	assertS(ok, "")
	littleEndian, ok := vm.curProto.stack.pop().(bool)
	assertS(ok, "")
	signed, ok := vm.curProto.stack.pop().(bool)
	assertS(ok, "") // TODO
	size, ok := vm.curProto.stack.pop().(int64)
	assertS(ok, "")
	offset, ok := vm.curProto.stack.pop().(int64)
	assertS(ok, "") // TODO
	buffer, ok := vm.curProto.stack.pop().([]byte)
	assertS(ok, "") // TODO
	var result interface{}
	switch size {
	case 1:
		var v uint8 = buffer[offset]
		if signed {
			result = int64(int8(v))
		} else {
			result = int64(v)
		}
	case 2:
		var v uint16
		if littleEndian {
			v = binary.LittleEndian.Uint16(buffer[offset:])
		} else {
			v = binary.BigEndian.Uint16(buffer[offset:])
		}
		if signed {
			result = int64(int16(v))
		} else {
			result = int64(v)
		}
	case 4:
		var v uint32
		if littleEndian {
			v = binary.LittleEndian.Uint32(buffer[offset:])
		} else {
			v = binary.BigEndian.Uint32(buffer[offset:])
		}
		if isFloat {
			result = float64(math.Float32frombits(v))
			break
		}
		if signed {
			result = int64(int32(v))
		} else {
			result = int64(v)
		}
	case 8:
		var v uint64
		if littleEndian {
			v = binary.LittleEndian.Uint64(buffer[offset:])
		} else {
			v = binary.BigEndian.Uint64(buffer[offset:])
		}
		if isFloat {
			result = float64(math.Float64frombits(v))
			break
		}
		result = int64(v) // TODO uint64 to int64 may overflow
	default:
		panic("") // TODO
	}
	vm.curProto.stack.Push(result)
	return 1
}

func builtinBufferNew(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 1 {
		panic("") // TODO
	}
	capacity, ok := vm.curProto.stack.pop().(int64)
	if !ok {
		panic("") // TODO
	}
	vm.curProto.stack.Push(make([]byte, capacity))
	return 1
}

func pushTwo(vm *VM, k, v interface{}) {
	vm.curProto.stack.Push(k)
	vm.curProto.stack.Push(v)
}

func builtinForeach(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 2 {
		panic("") // TODO
	}
	_func := vm.curProto.stack.pop()
	src := vm.curProto.stack.pop()
	switch src := src.(type) {
	case *[]interface{}:
		for idx, val := range *src {
			pushTwo(vm, int64(idx), val)
			call(_func, vm, 2, 0)
		}
	case string:
		for idx, ch := range src {
			pushTwo(vm, int64(idx), int64(ch))
			call(_func, vm, 2, 0)
		}
	case map[interface{}]interface{}:
		for k, v := range src {
			pushTwo(vm, k, v)
			call(_func, vm, 2, 0)
		}
	case []byte:
		for idx, val := range src {
			pushTwo(vm, int64(idx), int64(val))
			call(_func, vm, 2, 0)
		}
	default:
		panic("") // TODO
	}
	return 0
}

func builtinClone(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 1 {
		panic("") // TODO
	}
	switch src := vm.curProto.stack.pop().(type) {
	case *[]interface{}:
		arr := make([]interface{}, len(*src))
		copy(arr, *src)
		vm.curProto.stack.Push(&arr)
	case map[interface{}]interface{}:
		m := make(map[interface{}]interface{}, len(src))
		for k, v := range src {
			m[k] = v
		}
		vm.curProto.stack.Push(m)
	case []byte:
		arr := make([]byte, len(src))
		copy(arr, src)
		vm.curProto.stack.Push(arr)
	default:
		vm.curProto.stack.Push(src)
	}
	return 1
}

func builtinDelete(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 2 {
		panic("") // TODO
	}
	key := vm.curProto.stack.pop()
	m := vm.curProto.stack.pop().(map[interface{}]interface{}) // TODO
	delete(m, key)
	return 0
}

func builtinType(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 1 {
		panic("") // TODO
	}
	var t string
	switch vm.curProto.stack.pop().(type) {
	case string:
		t = "String"
	case *Closure:
		t = "Closure"
	case *builtinFunc:
		t = "Builtin"
	case map[interface{}]interface{}:
		t = "Object"
	case *[]interface{}:
		t = "Array"
	case []byte:
		t = "Buffer"
	case int64, float64:
		t = "Number"
	case bool:
		t = "Boolean"
	case nil:
		t = "Nil"
	default:
		panic("") // TODO
	}
	vm.curProto.stack.Push(t)
	return 1
}

func builtinSub(argCnt int, vm *VM) (retCnt int) {
	if argCnt < 2 || argCnt > 3 {
		panic("") // TODO
	}
	target := vm.curProto.stack.top(argCnt)
	start := vm.curProto.stack.top(argCnt - 1).(int64) // TODO
	var end int64
	if argCnt == 3 {
		end = vm.curProto.stack.Top().(int64)
	}
	vm.curProto.stack.popN(argCnt)
	switch target := target.(type) {
	case string:
		if argCnt == 2 {
			end = int64(len(target))
		}
		vm.curProto.stack.Push(target[start:end])
	case *[]interface{}:
		if argCnt == 2 {
			end = int64(len(*target))
		}
		subSlice := (*target)[start:end]
		vm.curProto.stack.Push(&subSlice)
	default:
		panic("") // TODO
	}
	return 1
}

func builtinAppend(argCnt int, vm *VM) (retCnt int) {
	if argCnt < 2 {
		panic("") // TODO
	}
	target := vm.curProto.stack.top(argCnt)
	arr, ok := target.(*[]interface{})
	if !ok {
		panic("") // TODO
	}
	n := argCnt
	for argCnt--; argCnt > 0; argCnt-- {
		arg := vm.curProto.stack.top(argCnt)
		*arr = append(*arr, arg)
	}
	vm.curProto.stack.popN(n)
	return 0
}

func builtinPrint(argCnt int, vm *VM) (retCnt int) {
	n := argCnt
	for ; argCnt > 0; argCnt-- {
		arg := vm.curProto.stack.top(argCnt)
		print(arg)
		fmt.Printf(" ")
	}
	vm.curProto.stack.popN(n)
	fmt.Println()
	return 0
}

func builtinLen(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 1 {
		panic("") // TODO
	}
	var length int64
	switch val := vm.curProto.stack.pop().(type) {
	case *[]interface{}:
		length = int64(len(*val))
	case map[interface{}]interface{}:
		length = int64(len(val))
	case string:
		length = int64(len(val))
	default:
		panic("") // TODO
	}
	vm.curProto.stack.Push(length)
	return 1
}

func print(val interface{}) {
	switch val := val.(type) {
	case *Closure:
		fmt.Printf("<closure>")
	case *builtinFunc:
		fmt.Printf("<builtin:\"%s\">", val.name)
	case string:
		fmt.Printf("\"%s\"", val)
	case map[interface{}]interface{}:
		fmt.Printf("Object{")
		i := 0
		for k, v := range val {
			print(k)
			fmt.Printf(": ")
			print(v)
			if i != len(val)-1 {
				fmt.Printf(", ")
				i++
			}
		}
		fmt.Printf("}")
	case *[]interface{}:
		fmt.Printf("Array[")
		for i, v := range *val {
			print(v)
			if i != len(*val)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf("]")
	case []byte:
		fmt.Printf("<Buffer>")
	default:
		fmt.Printf("%v", val)
	}
}

func assertS(condition bool, str string) {
	if !condition {
		panic(str)
	}
}
