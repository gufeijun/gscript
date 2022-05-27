package vm

import (
	"fmt"
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

// TODO Buffer?
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
	default:
		fmt.Printf("%v", val)
	}
}
