package vm

import (
	"encoding/binary"
	"fmt"
	"gscript/vm/types"
	"io"
	"io/fs"
	"math"
	"os"
)

type builtinFunc struct {
	handler func(argCnt int, vm *VM) int
	name    string
}

// if modify the following, should modify builtinFuncs in context.go too
var builtinFuncs = []builtinFunc{
	{builtinPrint, "print"},
	{builtinLen, "len"},
	{builtinAppend, "append"},
	{builtinSub, "sub"},
	{builtinType, "type"},
	{builtinDelete, "delete"},
	{builtinClone, "clone"},
	{builtinBufferNew, "__buffer_new"},
	{builtinBufferReadNumber, "__buffer_readNumber"},
	{builtinBufferWriteNumber, "__buffer_writeNumber"},
	{builtinBufferToString, "__buffer_toString"},
	{builtinBufferSlice, "__buffer_slice"},
	{builtinBufferConcat, "__buffer_concat"},
	{builtinBufferCopy, "__buffer_copy"},
	{builtinBufferFrom, "__buffer_from"},
	{builtinOpen, "__open"},
	{builtinRead, "__read"},
	{builtinWrite, "__write"},
	{builtinClose, "__close"},
	{builtinSeek, "__seek"},
	{builtinFChmod, "__fchmod"},
	{builtinChmod, "__chmod"},
	{builtinFChown, "__fchown"},
	{builtinChown, "__chown"},
	{builtinFChdir, "__fchdir"},
	{builtinChdir, "__chdir"},
	{builtinFStat, "__fstat"},
	{builtinStat, "__stat"},
	{builtinRename, "__rename"},
	{builtinExit, "__exit"},
	{builtinGetEnv, "__getenv"},
	{builtinSetEnv, "__setenv"},
	{builtinReadDir, "__readdir"},
	{builtinFReadDir, "__freaddir"},
}

// arg1: File, arg2: n
// return: []statObject
func builtinFReadDir(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 2, "")
	n, ok := pop(vm).(int64)
	assertS(ok, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	entrys, err := file.File.ReadDir(int(n))
	if err != nil {
		panic(err)
	}
	push(vm, newDirEntryObjects(entrys))
	return 1
}

func newDirEntryObjects(entrys []fs.DirEntry) *types.Array {
	arr := make([]interface{}, 0, len(entrys))
	for _, entry := range entrys {
		info, err := entry.Info()
		if err != nil {
			panic("") // TODO
		}
		arr = append(arr, newStat(info))
	}
	return types.NewArray(arr)
}

// arg1: pathname
// return: []statObject
func builtinReadDir(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	pathname, ok := pop(vm).(string)
	assertS(ok, "")
	entrys, err := os.ReadDir(pathname)
	if err != nil {
		panic("") // TODO
	}
	push(vm, newDirEntryObjects(entrys))
	return 1
}

// arg1: key, arg2: value
func builtinSetEnv(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 2, "")
	value, ok := pop(vm).(string)
	assertS(ok, "")
	key, ok := pop(vm).(string)
	assertS(ok, "")
	if err := os.Setenv(key, value); err != nil {
		panic(err) // TODO
	}
	return 0
}

// arg1: key
// return: value
func builtinGetEnv(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	key, ok := pop(vm).(string)
	assertS(ok, "")
	push(vm, os.Getenv(key))
	return 1
}

// arg1: code
func builtinExit(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	code, ok := pop(vm).(int64)
	assertS(ok, "")
	os.Exit(int(code))
	return 0
}

// arg1: oldpath, arg2: newpath
func builtinRename(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 2, "")
	newpath, ok := pop(vm).(string)
	assertS(ok, "")
	oldpath, ok := pop(vm).(string)
	assertS(ok, "")
	if err := os.Rename(oldpath, newpath); err != nil {
		panic(err)
	}
	return 0
}

func newStat(info fs.FileInfo) *types.Object {
	obj := types.NewObject()
	obj.Data["is_dir"] = info.IsDir()
	obj.Data["mode"] = int64(info.Mode())
	obj.Data["name"] = info.Name()
	obj.Data["size"] = info.Size()
	obj.Data["mod_time"] = info.ModTime().Unix()
	return obj
}

// arg1: filepath
// return: statObject
func builtinStat(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	path, ok := pop(vm).(string)
	assertS(ok, "")
	stat, err := os.Stat(path)
	if err != nil {
		panic(err) // TODO
	}
	push(vm, newStat(stat))
	return 1
}

// arg1: File
// return: statObject
func builtinFStat(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	stat, err := file.File.Stat()
	if err != nil {
		panic(err) // TODO
	}
	push(vm, newStat(stat))
	return 1
}

// arg1: File
func builtinFChdir(argCnt int, vm *VM) (retCnt int) {
	// var file *os.File
	assertS(argCnt == 1, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	err := file.File.Chdir()
	if err != nil {
		panic(err) // TODO
	}
	return 0
}

// arg1: File, arg2: uid, arg3: gid
func builtinFChown(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "")
	gid, ok := pop(vm).(int64)
	assertS(ok, "")
	uid, ok := pop(vm).(int64)
	assertS(ok, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	err := file.File.Chown(int(uid), int(gid))
	if err != nil {
		panic(err) // TODO
	}
	return 0
}

// arg1: File, arg2: mode
func builtinFChmod(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 2, "")
	mode, ok := pop(vm).(int64)
	assertS(ok, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	err := file.File.Chmod(os.FileMode(uint32(mode)))
	if err != nil {
		panic(err) // TODO
	}
	return 0
}

// arg1: path
func builtinChdir(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	path, ok := pop(vm).(string)
	assertS(ok, "")
	err := os.Chdir(path)
	if err != nil {
		panic(err) // TODO
	}
	return 0
}

// arg1: filepath, arg2: uid, arg3: gid
func builtinChown(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "")
	gid, ok := pop(vm).(int64)
	assertS(ok, "")
	uid, ok := pop(vm).(int64)
	assertS(ok, "")
	path, ok := pop(vm).(string)
	assertS(ok, "")
	err := os.Chown(path, int(uid), int(gid))
	if err != nil {
		panic(err) // TODO
	}
	return 0
}

// arg1: filepath, arg2: mode
func builtinChmod(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 2, "")
	mode, ok := pop(vm).(int64)
	assertS(ok, "")
	file, ok := pop(vm).(string)
	assertS(ok, "")
	err := os.Chmod(file, os.FileMode(uint32(mode)))
	if err != nil {
		panic(err) // TODO
	}
	return 0
}

// arg1: File, arg2: offset, arg3: whence("cur","end","start")
// return: n
func builtinSeek(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "")
	whenceS, ok := pop(vm).(string)
	assertS(ok, "")
	var whence int
	if whenceS == "cur" {
		whence = io.SeekCurrent
	} else if whenceS == "end" {
		whence = io.SeekEnd
	} else if whenceS == "start" {
		whence = io.SeekStart
	} else {
		panic("") // TODO
	}
	offset, ok := pop(vm).(int64)
	assertS(ok, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")

	n, err := file.File.Seek(offset, whence)
	if err != nil {
		panic(err) // TODO
	}
	push(vm, int64(n))
	return 1
}

// arg1: File
func builtinClose(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	file.File.Close()
	return 0
}

// arg1: File, arg2: Buffer, arg3: size
// arg1: File, arg2: stringTODv
// return: n
func builtinWrite(argCnt int, vm *VM) (retCnt int) {
	var data []byte
	if argCnt == 3 {
		size, ok := pop(vm).(int64)
		assertS(ok, "")
		buff, ok := pop(vm).(*types.Buffer)
		assertS(ok, "")
		data = buff.Data[:size]
	} else if argCnt == 2 {
		str, ok := pop(vm).(string)
		assertS(ok, "")
		data = []byte(str)
	} else {
		panic("") // TODO
	}
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	n, err := file.File.Write(data)
	if err != nil {
		panic("read failed") // TODO
	}
	push(vm, int64(n))
	return 1
}

// arg1: File, arg2: Buffer, arg3: size
// return: n
func builtinRead(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "")
	size, ok := pop(vm).(int64)
	assertS(ok, "")
	buff, ok := pop(vm).(*types.Buffer)
	assertS(ok, "")
	file, ok := pop(vm).(*types.File)
	assertS(ok, "")
	n, err := file.File.Read(buff.Data[:size])
	if err != nil {
		panic("read failed") // TODO
	}
	push(vm, int64(n))
	return 1
}

func getFileFlag(flagS string) int {
	var flag int
	var read, write bool
	for _, ch := range []byte(flagS) {
		if ch == 'r' {
			read = true
			continue
		}
		if ch == 'w' {
			write = true
			continue
		}
		if ch == 'c' {
			flag |= os.O_CREATE
			continue
		}
		if ch == 'a' {
			flag |= os.O_APPEND
			continue
		}
		if ch == 't' {
			flag |= os.O_TRUNC
			continue
		}
		if ch == 'e' {
			flag |= os.O_EXCL
		}
	}
	if read && write {
		flag |= os.O_RDWR
	} else if read && !write {
		flag |= os.O_RDONLY
	} else if write && !read {
		flag |= os.O_WRONLY
	}
	return flag
}

// arg1: filepath, arg2: flag(r,w,c,a,t,e), arg3: mode
// return: File
func builtinOpen(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "")
	mode, ok := pop(vm).(int64)
	assertS(ok, "")
	flagS, ok := pop(vm).(string)
	assertS(ok, "")
	filepath, ok := pop(vm).(string)
	assertS(ok, "")
	flag := getFileFlag(flagS)

	file, err := os.OpenFile(filepath, flag, os.FileMode(uint32(mode)))
	if err != nil {
		panic(err) // TODO
	}
	push(vm, types.NewFile(file))
	return 1
}

// arg1: string
// return: Buffer
func builtinBufferFrom(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 1, "")
	str, ok := pop(vm).(string)
	assertS(ok, "")
	push(vm, types.NewBufferFromString(str))
	return 1
}

// arg1: Buffer1, arg2: Buffer2, arg3: length, arg4: idx1, arg5: idx2
func builtinBufferCopy(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 5, "")
	idx2, ok := pop(vm).(int64)
	assertS(ok, "")
	idx1, ok := pop(vm).(int64)
	assertS(ok, "")
	length, ok := pop(vm).(int64)
	assertS(ok, "")
	buf2, ok := pop(vm).(*types.Buffer)
	assertS(ok, "")
	buf1, ok := pop(vm).(*types.Buffer)
	assertS(ok, "")
	copy(buf1.Data[idx1:], buf2.Data[idx2:idx2+length])
	return 0
}

// args: Buffer,Buffer[,Buffer]
// return: Buffer
func builtinBufferConcat(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt >= 2, "")
	bufs := make([][]byte, argCnt)
	var length int
	for i := argCnt - 1; i >= 0; i-- {
		buf, ok := pop(vm).(*types.Buffer)
		assertS(ok, "")
		length += len(buf.Data)
		bufs[i] = buf.Data
	}
	result := make([]byte, 0, length)
	for _, buf := range bufs {
		result = append(result, buf...)
	}
	push(vm, types.NewBuffer(result))
	return 1
}

// arg1: Buffer, arg2: offset, arg3: length
// return: Buffer
func builtinBufferSlice(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "argCnt should be 3")
	length, ok := pop(vm).(int64)
	assertS(ok, "")
	offset, ok := pop(vm).(int64)
	assertS(ok, "")
	buffer, ok := pop(vm).(*types.Buffer)
	assertS(ok, "")
	push(vm, types.NewBuffer(buffer.Data[offset:offset+length]))
	return 1
}

// arg1: Buffer, arg2: offset, arg3: length
// return: string
func builtinBufferToString(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 3, "")
	length, ok := pop(vm).(int64)
	assertS(ok, "")
	offset, ok := pop(vm).(int64)
	assertS(ok, "")
	buffer, ok := pop(vm).(*types.Buffer)
	assertS(ok, "")
	push(vm, string(buffer.Data[offset:offset+length]))
	return 1
}

// arg1: Buffer, arg2: offset, arg3: size, arg4: littleEndian arg5: number
func builtinBufferWriteNumber(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 5, "") // TODO
	number := pop(vm)
	littleEndian, ok := pop(vm).(bool)
	assertS(ok, "")
	size, ok := pop(vm).(int64)
	assertS(ok, "")
	offset, ok := pop(vm).(int64)
	assertS(ok, "") // TODO
	buf, ok := pop(vm).(*types.Buffer)
	assertS(ok, "") // TODO

	switch size {
	case 1:
		v := byte(number.(int64))
		buf.Data[offset] = v
	case 2:
		v := uint16(number.(int64))
		if littleEndian {
			binary.LittleEndian.PutUint16(buf.Data[offset:], v)
		} else {
			binary.BigEndian.PutUint16(buf.Data[offset:], v)
		}
	case 4:
		var v uint32
		if vf, ok := number.(float64); ok {
			v = math.Float32bits(float32(vf))
		} else {
			v = uint32(number.(int64))
		}
		if littleEndian {
			binary.LittleEndian.PutUint32(buf.Data[offset:], v)
		} else {
			binary.BigEndian.PutUint32(buf.Data[offset:], v)
		}
	case 8:
		var v uint64
		if vf, ok := number.(float64); ok {
			v = math.Float64bits(vf)
		} else {
			v = uint64(number.(int64))
		}
		if littleEndian {
			binary.LittleEndian.PutUint64(buf.Data[offset:], v)
		} else {
			binary.BigEndian.PutUint64(buf.Data[offset:], v)
		}
	}
	return 0
}

// arg1: Buffer, arg2: offset, arg3: size, arg4: signed arg5: littleEndian arg6: isFloat
// return: Number
func builtinBufferReadNumber(argCnt int, vm *VM) (retCnt int) {
	assertS(argCnt == 6, "") // TODO
	isFloat, ok := pop(vm).(bool)
	assertS(ok, "")
	littleEndian, ok := pop(vm).(bool)
	assertS(ok, "")
	signed, ok := pop(vm).(bool)
	assertS(ok, "") // TODO
	size, ok := pop(vm).(int64)
	assertS(ok, "")
	offset, ok := pop(vm).(int64)
	assertS(ok, "") // TODO
	buf, ok := pop(vm).(*types.Buffer)
	assertS(ok, "") // TODO
	var result interface{}
	switch size {
	case 1:
		var v uint8 = buf.Data[offset]
		if signed {
			result = int64(int8(v))
		} else {
			result = int64(v)
		}
	case 2:
		var v uint16
		if littleEndian {
			v = binary.LittleEndian.Uint16(buf.Data[offset:])
		} else {
			v = binary.BigEndian.Uint16(buf.Data[offset:])
		}
		if signed {
			result = int64(int16(v))
		} else {
			result = int64(v)
		}
	case 4:
		var v uint32
		if littleEndian {
			v = binary.LittleEndian.Uint32(buf.Data[offset:])
		} else {
			v = binary.BigEndian.Uint32(buf.Data[offset:])
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
			v = binary.LittleEndian.Uint64(buf.Data[offset:])
		} else {
			v = binary.BigEndian.Uint64(buf.Data[offset:])
		}
		if isFloat {
			result = float64(math.Float64frombits(v))
			break
		}
		result = int64(v) // TODO uint64 to int64 may overflow
	default:
		panic("") // TODO
	}
	push(vm, result)
	return 1
}

// arg1: capacity
// return: Buffer
func builtinBufferNew(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 1 {
		panic("") // TODO
	}
	capacity, ok := pop(vm).(int64)
	if !ok {
		panic("") // TODO
	}
	push(vm, types.NewBufferN(int(capacity)))
	return 1
}

func builtinClone(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 1 {
		panic("") // TODO
	}
	switch src := pop(vm).(type) {
	case *types.Array:
		arr := make([]interface{}, len(src.Data))
		copy(arr, src.Data)
		push(vm, types.NewArray(arr))
	case *types.Object:
		obj := types.NewObjectN(len(src.Data))
		for k, v := range src.Data {
			obj.Data[k] = v
		}
		push(vm, obj)
	case *types.Buffer:
		data := make([]byte, len(src.Data))
		copy(data, src.Data)
		push(vm, types.NewBuffer(data))
	default:
		push(vm, src)
	}
	return 1
}

// arg1: Object, arg2: key
func builtinDelete(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 2 {
		panic("") // TODO
	}
	key := pop(vm)
	obj := pop(vm).(*types.Object) // TODO
	delete(obj.Data, key)
	return 0
}

func builtinType(argCnt int, vm *VM) (retCnt int) {
	if argCnt != 1 {
		panic("") // TODO
	}
	var t string
	switch pop(vm).(type) {
	case string:
		t = "String"
	case *types.Closure:
		t = "Closure"
	case *builtinFunc:
		t = "Builtin"
	case *types.Object:
		t = "Object"
	case *types.Array:
		t = "Array"
	case *types.Buffer:
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
	push(vm, t)
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
		push(vm, target[start:end])
	case *types.Array:
		if argCnt == 2 {
			end = int64(len(target.Data))
		}
		subSlice := target.Data[start:end]
		push(vm, types.NewArray(subSlice))
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
	arr, ok := target.(*types.Array)
	if !ok {
		panic("") // TODO
	}
	n := argCnt
	for argCnt--; argCnt > 0; argCnt-- {
		arg := vm.curProto.stack.top(argCnt)
		arr.Data = append(arr.Data, arg)
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
	switch val := pop(vm).(type) {
	case *types.Array:
		length = int64(len(val.Data))
	case *types.Object:
		length = int64(len(val.Data))
	case string:
		length = int64(len(val))
	case *types.Buffer:
		length = int64(len(val.Data))
	default:
		panic("") // TODO
	}
	push(vm, length)
	return 1
}

func print(val interface{}) {
	switch val := val.(type) {
	case *types.Closure:
		fmt.Printf("<closure>")
	case *builtinFunc:
		fmt.Printf("<builtin:\"%s\">", val.name)
	case string:
		fmt.Printf("%s", val)
	case *types.Object:
		fmt.Printf("Object{")
		i := 0
		for k, v := range val.Data {
			print(k)
			fmt.Printf(": ")
			print(v)
			if i != len(val.Data)-1 {
				fmt.Printf(", ")
				i++
			}
		}
		fmt.Printf("}")
	case *types.Array:
		fmt.Printf("Array[")
		for i, v := range val.Data {
			print(v)
			if i != len(val.Data)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf("]")
	case *types.Buffer:
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

func push(vm *VM, val interface{}) {
	vm.curProto.stack.Push(val)
}

func pop(vm *VM) interface{} {
	return vm.curProto.stack.pop()
}
