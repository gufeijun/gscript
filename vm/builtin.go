package vm

import "fmt"

type builtinFunc struct {
	handler func(argCnt int, vm *VM) int
	name    string
}

var builtinFuncs = []builtinFunc{
	{builtinPrint, "print"},
}

func builtinPrint(argCnt int, vm *VM) (retCnt int) {
	for ; argCnt > 0; argCnt-- {
		arg := vm.curProto.stack.top(argCnt)
		print(arg)
		fmt.Printf(" ")
	}
	fmt.Println()
	return 0
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
		fmt.Printf("{")
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
	case []interface{}:
		fmt.Printf("[")
		for i, v := range val {
			print(v)
			if i != len(val)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf("]")
	default:
		fmt.Printf("%v", val)
	}
}
