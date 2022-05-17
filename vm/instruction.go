package vm

import (
	"fmt"
	"gscript/proto"
	"strconv"
)

var actions = []func(vm *VM){
	actionUnaryNOT,
	actionUnaryNEG,
	actionUnaryLNOT,
	actionBinaryADD,
	actionBinarySUB,
	actionBinaryMUL,
	actionBinaryDIV,
	actionBinaryMOD,
	actionBinaryAND,
	actionBinaryXOR,
	actionBinaryOR,
	actionBinaryIDIV,
	actionBinarySHR,
	actionBinarySHL,
	actionBinaryLE,
	actionBinaryGE,
	actionBinaryLT,
	actionBinaryGT,
	actionBinaryEQ,
	actionBinaryNE,
	actionBinaryLAND,
	actionBinaryLOR,
	actionBinaryATTR,
	actionLoadNil,
	actionLoadConst,
	actionLoadName,
	actionLoadFunc,
	actionLoadAnonymous,
	actionLoadUpValue,
	actionStoreName,
	actionStoreUpValue,
	actionStoreKV,
	actionPushNameNil,
	actionPushName,
	actionCopyName,
	actionResizeNameTable,
	actionPopTop,
	actionStop,
	actionSliceNew,
	actionSliceAppend,
	actionNewMap,
	actionNewEmptyMap,
	actionAttrAssign,
	actionAttrAssignAddEq,
	actionAttrAssignSubEq,
	actionAttrAssignMulEq,
	actionAttrAssignDivEq,
	actionAttrAssignModEq,
	actionAttrAssignAndEq,
	actionAttrAssignXorEq,
	actionAttrAssignOrEq,
	actionAttrAccess,
	actionJumpRel,
	actionJumpAbs,
	actionJumpIf,
	actionJumpCase,
	actionCall,
	actionReturn,
	actionRotTwo,
}

func actionUnaryNOT(vm *VM) {
	top := vm.stack.Top()
	if v, ok := top.(int64); ok {
		vm.stack.Replace(^int64(v))
		return
	}
	panic("") // TODO
}

func actionUnaryNEG(vm *VM) {
	top := vm.stack.Top()
	if v, ok := top.(int64); ok {
		vm.stack.Replace(-int64(v))
		return
	}
	if v, ok := top.(float64); ok {
		vm.stack.Replace(-float64(v))
		return
	}
	panic("") // TODO
}

func actionUnaryLNOT(vm *VM) {
	vm.stack.Replace(!getBool(vm.stack.Top()))
}

func getTop2(vm *VM) (v1, v2 interface{}) {
	v2 = vm.stack.Top()
	vm.stack.Pop()
	v1 = vm.stack.Top()
	return
}

func actionBinaryADD(vm *VM) {
	vm.stack.Replace(addAction(getTop2(vm)))
}

func actionBinarySUB(vm *VM) {
	vm.stack.Replace(subAction(getTop2(vm)))
}

func actionBinaryMUL(vm *VM) {
	vm.stack.Replace(mulAction(getTop2(vm)))
}

func actionBinaryDIV(vm *VM) {
	vm.stack.Replace(divAction(getTop2(vm)))
}

func actionBinaryIDIV(vm *VM) {
	vm.stack.Replace(idivAction(getTop2(vm)))
}

func actionBinaryMOD(vm *VM) {
	vm.stack.Replace(modAction(getTop2(vm)))
}

func actionBinarySHL(vm *VM) {
	vm.stack.Replace(shlAction(getTop2(vm)))
}

func actionBinarySHR(vm *VM) {
	vm.stack.Replace(shrAction(getTop2(vm)))
}

func actionBinaryAND(vm *VM) {
	vm.stack.Replace(andAction(getTop2(vm)))
}

func actionBinaryOR(vm *VM) {
	vm.stack.Replace(orAction(getTop2(vm)))
}

func actionBinaryXOR(vm *VM) {
	vm.stack.Replace(xorAction(getTop2(vm)))
}

func actionBinaryLE(vm *VM) {
	vm.stack.Replace(leAction(getTop2(vm)))
}

func actionBinaryGE(vm *VM) {
	vm.stack.Replace(geAction(getTop2(vm)))
}

func actionBinaryLT(vm *VM) {
	vm.stack.Replace(ltAction(getTop2(vm)))
}

func actionBinaryGT(vm *VM) {
	vm.stack.Replace(gtAction(getTop2(vm)))
}

func actionBinaryEQ(vm *VM) {
	vm.stack.Replace(eqAction(getTop2(vm)))
}

func actionBinaryNE(vm *VM) {
	vm.stack.Replace(neAction(getTop2(vm)))
}

func actionBinaryLAND(vm *VM) {
	vm.stack.Replace(landAction(getTop2(vm)))
}

func actionBinaryLOR(vm *VM) {
	vm.stack.Replace(lorAction(getTop2(vm)))
}

func actionBinaryATTR(vm *VM) {
	key := vm.stack.Top()
	vm.stack.Pop()
	obj := vm.stack.Top()
	if slice, ok := obj.([]interface{}); ok {
		var idx int64
		if idx, ok = key.(int64); !ok {
			panic("array index should be integer") // TODO
		}
		if idx > int64(len(slice)) {
			panic("index out of range") // TODO
		}
		vm.stack.Replace(slice[idx])
		return
	}
	if _map, ok := obj.(map[interface{}]interface{}); ok {
		if key == nil {
			panic("map key should not be nil")
		}
		vm.stack.Replace(_map[key])
		return
	}
	// TODO
	panic(fmt.Sprintf("do not support attr access for %T", obj))
}

func actionLoadNil(vm *VM) {
	vm.stack.Push(nil)
}

func actionLoadConst(vm *VM) {
	vm.stack.Push(vm.constTable[vm.getOpNum()])
}

func actionLoadName(vm *VM) {
	vm.stack.Push(vm.frame.symbolTable.getValue(vm.getOpNum()))
}

func actionLoadFunc(vm *VM) {
	f := &vm.funcTable[vm.getOpNum()]
	if f.UpValueTable == nil {
		table := make([]*GsValue, 0, len(f.UpValues))
		for _, nameIdx := range f.UpValues {
			v := vm.topFrame.symbolTable.values[nameIdx]
			table = append(table, v)
		}
		f.UpValueTable = table
	}
	vm.stack.Push(&Closure{
		Info:     f.Info,
		UpValues: f.UpValueTable.([]*GsValue),
	})
}

func actionLoadAnonymous(vm *VM) {
	f := &vm.anonymousTable[vm.getOpNum()]
	closure := &Closure{
		Info:     f.Info,
		UpValues: make([]*GsValue, 0, len(f.UpValues)),
	}
	for _, upValue := range f.UpValues {
		var v *GsValue
		if !upValue.DirectDependent {
			v = vm.frame.upValues[upValue.Index]
		} else {
			v = vm.frame.symbolTable.values[upValue.Index]
		}
		closure.UpValues = append(closure.UpValues, v)
	}
	vm.stack.Push(closure)
}

func actionLoadUpValue(vm *VM) {
	vm.stack.Push(vm.frame.upValues[vm.getOpNum()].value)
}

func actionStoreName(vm *VM) {
	vm.frame.symbolTable.setValue(vm.getOpNum(), vm.stack.Top())
	vm.stack.Pop()
}

func actionStoreUpValue(vm *VM) {
	vm.frame.upValues[vm.getOpNum()].value = vm.stack.Top()
	vm.stack.Pop()
}

func actionStoreKV(vm *VM) {
	obj, ok := vm.frame.symbolTable.top().(map[interface{}]interface{})
	if !ok {
		panic("STORE_KV: target is not an object") // TODO
	}
	val := vm.stack.Top()
	vm.stack.Pop()
	key := vm.stack.Top()
	vm.stack.Pop()
	obj[key] = val
}

func actionPushNameNil(vm *VM) {
	vm.frame.symbolTable.pushSymbol(nil)
}

func actionPushName(vm *VM) {
	vm.frame.symbolTable.pushSymbol(vm.stack.Top())
	vm.stack.Pop()
}

func actionCopyName(vm *VM) {
	vm.frame.symbolTable.pushSymbol(vm.stack.Top())
}

func actionResizeNameTable(vm *VM) {
	length := int(vm.getOpNum())
	if length >= len(vm.frame.symbolTable.values) {
		return
	}
	vm.frame.symbolTable.values = vm.frame.symbolTable.values[:length]
}

func actionPopTop(vm *VM) {
	vm.stack.Pop()
}

func actionStop(vm *VM) {
	vm.Stop()
}

func actionSliceNew(vm *VM) {
	cnt := vm.getOpNum()
	arr := make([]interface{}, cnt)
	for i := int(cnt) - 1; i >= 0; i-- {
		val := vm.stack.Top()
		vm.stack.Pop()
		arr[i] = val
	}
	vm.stack.Push(arr)
}

func actionSliceAppend(vm *VM) {
	val := vm.stack.Top()
	vm.stack.Pop()
	top := vm.stack.Top()
	if slice, ok := top.([]interface{}); ok {
		slice = append(slice, val)
		vm.stack.Replace(slice)
		return
	}
	panic("append operate for illegal type") // TODO
}

func actionNewMap(vm *VM) {
	m := make(map[interface{}]interface{})
	cnt := vm.getOpNum()
	for i := 0; i < int(cnt); i++ {
		val := vm.stack.Top()
		vm.stack.Pop()
		key := vm.stack.Top()
		vm.stack.Pop()
		if key == nil {
			panic("map key should not be nil") // TODO
		}
		m[key] = val
	}
	vm.stack.Push(m)
}

func actionNewEmptyMap(vm *VM) {
	vm.stack.Push(map[interface{}]interface{}{})
}

func actionAttrAssign(vm *VM) {
	attrAssign(vm, func(ori, val interface{}) interface{} {
		return val
	})
}

func actionAttrAssignAddEq(vm *VM) {
	attrAssign(vm, addAction)
}

func actionAttrAssignSubEq(vm *VM) {
	attrAssign(vm, subAction)
}

func actionAttrAssignMulEq(vm *VM) {
	attrAssign(vm, mulAction)
}

func actionAttrAssignDivEq(vm *VM) {
	attrAssign(vm, divAction)
}

func actionAttrAssignModEq(vm *VM) {
	attrAssign(vm, modAction)
}

func actionAttrAssignAndEq(vm *VM) {
	attrAssign(vm, andAction)
}

func actionAttrAssignXorEq(vm *VM) {
	attrAssign(vm, xorAction)
}

func actionAttrAssignOrEq(vm *VM) {
	attrAssign(vm, orAction)
}

func attrAssign(vm *VM, cb func(ori, val interface{}) interface{}) {
	obj := vm.stack.Top()
	vm.stack.Pop()
	key := vm.stack.Top()
	vm.stack.Pop()
	val := vm.stack.Top()
	vm.stack.Pop()
	if slice, ok := obj.([]interface{}); ok {
		var idx int64
		if idx, ok = key.(int64); !ok {
			panic("array index should be integer") // TODO
		}
		if idx > int64(len(slice)) {
			panic("index out of range")
		}
		slice[idx] = cb(slice[idx], val)
		return
	}
	if _map, ok := obj.(map[interface{}]interface{}); ok {
		if key == nil {
			panic("map key should not be nil")
		}
		_map[key] = cb(_map[key], val)
		return
	}
	panic(fmt.Sprintf("do not support attr assign for %T", obj))

}

func actionAttrAccess(vm *VM) {
	obj := vm.stack.Top()
	vm.stack.Pop()
	key := vm.stack.Top()
	if slice, ok := obj.([]interface{}); ok {
		var idx int64
		if idx, ok = key.(int64); !ok {
			panic("array index should be integer") // TODO
		}
		if idx > int64(len(slice)) {
			panic("index out of range")
		}
		vm.stack.Replace(slice[idx])
		return
	}
	if _map, ok := obj.(map[interface{}]interface{}); ok {
		if key == nil {
			panic("map key should not be nil")
		}
		vm.stack.Replace(_map[key])
		return
	}
	// TODO
	panic(fmt.Sprintf("do not support attr access for %T", obj))
}

func actionJumpRel(vm *VM) {
	steps := vm.getOpNum()
	vm.frame.pc += steps
}

func actionJumpAbs(vm *VM) {
	vm.frame.pc = vm.getOpNum()
}

func actionJumpIf(vm *VM) {
	top := vm.stack.Top()
	vm.stack.Pop()
	steps := vm.getOpNum()
	if getBool(top) {
		vm.frame.pc += steps
	}
}

func actionJumpCase(vm *VM) {
	caseCond := vm.stack.Top()
	vm.stack.Pop()
	switchCond := vm.stack.Top()
	steps := vm.getOpNum()
	if eqAction(caseCond, switchCond).(bool) {
		vm.frame.pc += steps
	}
}

func actionCall(vm *VM) {
	_func := vm.stack.Top().(*Closure)
	vm.stack.Pop()

	wantRtnCnt := int(vm.frame.text[vm.frame.pc])
	vm.frame.pc++
	argCnt := uint32(vm.frame.text[vm.frame.pc])
	vm.frame.pc++

	// generate a new function call frame
	frame := &stackFrame{
		pc:          0,
		prev:        vm.frame,
		symbolTable: newSymbolTable(),
		wantRetCnt:  wantRtnCnt,
		text:        _func.Info.Text,
		upValues:    _func.UpValues,
	}
	vm.frame = frame

	parCnt := uint32(len(_func.Info.Parameters))

	// if arguments is fewer than parameters, push several nil values to make up
	for argCnt < parCnt {
		vm.stack.Push(_func.Info.Parameters[argCnt].Default)
		argCnt++
	}
	if _func.Info.VaArgs {
		// collect VaArgs
		i := argCnt - parCnt
		arr := make([]interface{}, i)
		for i > 0 {
			i--
			arr[i] = vm.stack.Top()
			vm.stack.Pop()
			argCnt--
		}
		vm.stack.Push(arr)
	} else {
		// pop out extra arguments
		for argCnt > parCnt {
			vm.stack.Pop()
			argCnt--
		}
	}
}

func actionReturn(vm *VM) {
	realRtnCnt := int(vm.getOpNum())
	wantRtnCnt := vm.frame.wantRetCnt

	for wantRtnCnt < realRtnCnt {
		vm.stack.Pop()
		wantRtnCnt++
	}
	for wantRtnCnt > realRtnCnt {
		vm.stack.Push(nil)
		wantRtnCnt--
	}

	if vm.frame.prev == nil {
		vm.Stop()
	}
	vm.frame.symbolTable = nil
	vm.frame = vm.frame.prev
}

func actionRotTwo(vm *VM) {
	top := len(vm.stack.Buf) - 1
	vm.stack.Buf[top], vm.stack.Buf[top-1] = vm.stack.Buf[top-1], vm.stack.Buf[top]
}

func Execute(vm *VM, ins proto.Instruction) {
	actions[ins](vm)
}

func getBool(val interface{}) bool {
	if v, ok := val.(bool); ok {
		return v
	}
	if v, ok := val.(int64); ok {
		return v != 0
	}
	if v, ok := val.(float64); ok {
		return v != 0
	}
	return val != nil
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func boolToFloat(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func addAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, addInt, addFloat, addBool, addString)
}

func subAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, subInt, subFloat, subBool, nil)
}

func mulAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, mulInt, mulFloat, mulBool, nil)
}

func divAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, divInt, divFloat, nil, nil)
}

func idivAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, idivInt, idivFloat, nil, nil)
}

func modAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, modInt, modFloat, nil, nil)
}

func andAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, andInt, nil, nil, nil)
}

func orAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, orInt, nil, nil, nil)
}

func xorAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, xorInt, nil, nil, nil)
}

func shrAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, shrInt, nil, nil, nil)
}

func shlAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, shlInt, nil, nil, nil)
}

func leAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, leInt, leFloat, leBool, leString)
}

func geAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, geInt, geFloat, geBool, geString)
}

func ltAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, ltInt, ltFloat, ltBool, ltString)
}

func gtAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, gtInt, gtFloat, gtBool, gtString)
}

func eqAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, eqInt, eqFloat, eqBool, eqString)
}

func neAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, neInt, neFloat, neBool, neString)
}

func landAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, landInt, landFloat, landBool, landString)
}

func lorAction(v1, v2 interface{}) interface{} {
	return binaryAction(v1, v2, lorInt, lorFloat, lorBool, lorString)
}

func binaryAction(v1, v2 interface{}, intOP func(a, b int64) interface{},
	floatOP func(a, b float64) interface{}, boolOP func(a, b bool) interface{},
	stringOP func(a, b string) interface{}) interface{} {
	if v1 == nil || v2 == nil {
		return boolOP(getBool(v1), getBool(v2))
	}
	var result interface{}

	switch v1 := v1.(type) {
	case int64:
		if v, ok := v2.(int64); ok {
			result = intOP(v1, v)
			break
		}
		if v, ok := v2.(float64); ok {
			result = floatOP(float64(v1), v)
			break
		}
		if v, ok := v2.(bool); ok {
			result = intOP(v1, boolToInt(v))
			break
		}
		if v, ok := v2.(string); ok {
			result = stringOP(strconv.Itoa(int(v1)), v)
			break
		}
		panic("") // TODO
	case float64:
		if v, ok := v2.(int64); ok {
			result = floatOP(v1, float64(v))
			break
		}
		if v, ok := v2.(float64); ok {
			result = floatOP(v1, v)
			break
		}
		if v, ok := v2.(bool); ok {
			result = floatOP(v1, boolToFloat(v))
			break
		}
		if v, ok := v2.(string); ok {
			result = stringOP(fmt.Sprintf("%f", v1), v)
			break
		}
		panic("") // TODO
	case bool:
		if v, ok := v2.(int64); ok {
			result = intOP(boolToInt(v1), v)
			break
		}
		if v, ok := v2.(float64); ok {
			result = floatOP(boolToFloat(v1), v)
			break
		}
		if v, ok := v2.(bool); ok {
			result = boolOP(v1, v)
			break
		}
		if _, ok := v2.(string); ok {
			result = boolOP(v1, true)
			break
		}
	case string:
		if v, ok := v2.(int64); ok {
			result = stringOP(v1, strconv.Itoa(int(v)))
			break
		}
		if v, ok := v2.(float64); ok {
			result = stringOP(v1, fmt.Sprintf("%f", v))
			break
		}
		if v, ok := v2.(bool); ok {
			result = boolOP(true, v)
			break
		}
		if v, ok := v2.(string); ok {
			result = stringOP(v1, v)
			break
		}
	default:
		panic("")
	}

	if v, ok := result.(float64); ok {
		vi := int64(v)
		if float64(vi) == v {
			result = vi
		}
	}
	return result
}

func addInt(a, b int64) interface{} {
	return a + b
}

func addFloat(a, b float64) interface{} {
	return a + b
}

func addBool(a, b bool) interface{} {
	return boolToInt(a) + boolToInt(b)
}

func addString(a, b string) interface{} {
	return a + b
}

func subInt(a, b int64) interface{} {
	return a - b
}

func subFloat(a, b float64) interface{} {
	return a - b
}

func subBool(a, b bool) interface{} {
	return boolToInt(a) - boolToInt(b)
}

func mulInt(a, b int64) interface{} {
	return a * b
}

func mulFloat(a, b float64) interface{} {
	return a * b
}

func mulBool(a, b bool) interface{} {
	if a && b {
		return int64(1)
	}
	return 0
}

func divInt(a, b int64) interface{} {
	return float64(a) / float64(b)
}

func divFloat(a, b float64) interface{} {
	return a / b
}

func idivInt(a, b int64) interface{} {
	return a / b
}

func idivFloat(a, b float64) interface{} {
	return int64(a / b)
}

func modInt(a, b int64) interface{} {
	return a % b
}

func modFloat(a, b float64) interface{} {
	var i int64
	for v := b; v < a; v += b {
		i++
	}
	return i
}

func shlInt(a, b int64) interface{} {
	return a << b
}

func shrInt(a, b int64) interface{} {
	return a >> b
}

func andInt(a, b int64) interface{} {
	return a & b
}

func orInt(a, b int64) interface{} {
	return a != 0 || b != 0
}

func xorInt(a, b int64) interface{} {
	return a ^ b
}

func leInt(a, b int64) interface{} {
	return a <= b
}

func leFloat(a, b float64) interface{} {
	return a <= b
}

func leBool(a, b bool) interface{} {
	return !(a && !b)
}

func leString(a, b string) interface{} {
	return a <= b
}

func geInt(a, b int64) interface{} {
	return a >= b
}

func geFloat(a, b float64) interface{} {
	return a >= b
}

func geBool(a, b bool) interface{} {
	return !(!a && b)
}

func geString(a, b string) interface{} {
	return a >= b
}

func ltInt(a, b int64) interface{} {
	return a < b
}

func ltFloat(a, b float64) interface{} {
	return a < b
}

func ltBool(a, b bool) interface{} {
	return !a && b
}

func ltString(a, b string) interface{} {
	return a < b
}

func gtInt(a, b int64) interface{} {
	return a > b
}

func gtFloat(a, b float64) interface{} {
	return a > b
}

func gtBool(a, b bool) interface{} {
	return a && !b
}

func gtString(a, b string) interface{} {
	return a > b
}

func eqInt(a, b int64) interface{} {
	return a == b
}

func eqFloat(a, b float64) interface{} {
	return a == b
}

func eqBool(a, b bool) interface{} {
	return a == b
}

func eqString(a, b string) interface{} {
	return a == b
}

func neInt(a, b int64) interface{} {
	return a != b
}

func neFloat(a, b float64) interface{} {
	return a != b
}

func neBool(a, b bool) interface{} {
	return a != b
}

func neString(a, b string) interface{} {
	return a != b
}

func landInt(a, b int64) interface{} {
	return a != 0 && b != 0
}

func landFloat(a, b float64) interface{} {
	return a != 0 && b != 0
}

func landBool(a, b bool) interface{} {
	return a && b
}

func landString(a, b string) interface{} {
	return true
}

func lorInt(a, b int64) interface{} {
	return a != 0 || b != 0
}

func lorFloat(a, b float64) interface{} {
	return a != 0 || b != 0
}

func lorBool(a, b bool) interface{} {
	return a || b
}

func lorString(a, b string) interface{} {
	return true
}
