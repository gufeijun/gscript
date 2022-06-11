package vm

import (
	"bytes"
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

func GetBuiltinFuncNameByNum(num uint32) string {
	return builtinFuncs[num].name
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
	{builtinRemove, "__remove"},
	{builtinFChmod, "__fchmod"},
	{builtinChmod, "__chmod"},
	{builtinFChown, "__fchown"},
	{builtinChown, "__chown"},
	{builtinFChdir, "__fchdir"},
	{builtinChdir, "__chdir"},
	{builtinFStat, "__fstat"},
	{builtinStat, "__stat"},
	{builtinRename, "__rename"},
	{builtinMkdir, "__mkdir"},
	{builtinExit, "__exit"},
	{builtinGetEnv, "__getenv"},
	{builtinSetEnv, "__setenv"},
	{builtinReadDir, "__readdir"},
	{builtinFReadDir, "__freaddir"},
	{builtinThrow, "throw"},
}

func builtinMkdir(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 2)
	mode, ok := pop(vm).(int64)
	vm.assert(ok)
	str, ok := pop(vm).(string)
	vm.assert(ok)
	if err := os.Mkdir(str, os.FileMode(uint32(mode))); err != nil {
		throw(err, vm)
	}
	return 0
}

func throw(err error, vm *VM) {
	vm.builtinFuncFailed = true
	push(vm, err.Error())
	_throw(vm)
}

func _throw(vm *VM) {
	for {
		frame := vm.curProto.frame
		tryInfos := frame.tryInfos
		if len(tryInfos) > 0 {
			catchAddr, varCnt := frame.popTryInfo()
			frame.symbolTable.resizeTo(int(varCnt))
			frame.pc = catchAddr
			break
		}
		if frame.prev == nil {
			var buf bytes.Buffer
			fprint(&buf, pop(vm))
			vm.exit("uncaught exception: %s", buf.String())
		}
		vm.curProto.frame = frame.prev
	}
}

func builtinThrow(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	_throw(vm)
	return 0
}

func builtinRemove(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	path, ok := pop(vm).(string)
	vm.assert(ok)
	if err := os.Remove(path); err != nil {
		throw(err, vm)
	}
	return 0
}

// arg1: File, arg2: n
// return: []statObject
func builtinFReadDir(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 2)
	n, ok := pop(vm).(int64)
	vm.assert(ok)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	entrys, err := file.File.ReadDir(int(n))
	if err != nil {
		throw(err, vm)
		return 0
	}
	arr, err := newDirEntryObjects(entrys)
	if err != nil {
		throw(err, vm)
		return 0
	}
	push(vm, arr)
	return 1
}

func newDirEntryObjects(entrys []fs.DirEntry) (*types.Array, error) {
	arr := make([]interface{}, 0, len(entrys))
	for _, entry := range entrys {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		arr = append(arr, newStat(info))
	}
	return types.NewArray(arr), nil
}

// arg1: pathname
// return: []statObject
func builtinReadDir(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	pathname, ok := pop(vm).(string)
	vm.assert(ok)
	entrys, err := os.ReadDir(pathname)
	if err != nil {
		throw(err, vm)
		return 0
	}
	arr, err := newDirEntryObjects(entrys)
	if err != nil {
		throw(err, vm)
		return 0
	}
	push(vm, arr)
	return 1
}

// arg1: key, arg2: value
func builtinSetEnv(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 2)
	value, ok := pop(vm).(string)
	vm.assert(ok)
	key, ok := pop(vm).(string)
	vm.assert(ok)
	if err := os.Setenv(key, value); err != nil {
		throw(err, vm)
		return 0
	}
	return 0
}

// arg1: key
// return: value
func builtinGetEnv(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	key, ok := pop(vm).(string)
	vm.assert(ok)
	push(vm, os.Getenv(key))
	return 1
}

// arg1: code
func builtinExit(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	code, ok := pop(vm).(int64)
	vm.assert(ok)
	os.Exit(int(code))
	return 0
}

// arg1: oldpath, arg2: newpath
func builtinRename(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 2)
	newpath, ok := pop(vm).(string)
	vm.assert(ok)
	oldpath, ok := pop(vm).(string)
	vm.assert(ok)
	if err := os.Rename(oldpath, newpath); err != nil {
		throw(err, vm)
	}
	return 0
}

func newStat(info fs.FileInfo) *types.Object {
	obj := types.NewObject()
	obj.Set("is_dir", info.IsDir())
	obj.Set("mode", int64(info.Mode()))
	obj.Set("name", info.Name())
	obj.Set("size", info.Size())
	obj.Set("mod_time", info.ModTime().Unix())
	return obj
}

// arg1: filepath
// return: statObject
func builtinStat(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	path, ok := pop(vm).(string)
	vm.assert(ok)
	stat, err := os.Stat(path)
	if err != nil {
		throw(err, vm)
		return 0
	}
	push(vm, newStat(stat))
	return 1
}

// arg1: File
// return: statObject
func builtinFStat(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	stat, err := file.File.Stat()
	if err != nil {
		throw(err, vm)
		return 0
	}
	push(vm, newStat(stat))
	return 1
}

// arg1: File
func builtinFChdir(argCnt int, vm *VM) (retCnt int) {
	// var file *os.File
	vm.assert(argCnt == 1)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	err := file.File.Chdir()
	if err != nil {
		throw(err, vm)
		return 0
	}
	return 0
}

// arg1: File, arg2: uid, arg3: gid
func builtinFChown(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 3)
	gid, ok := pop(vm).(int64)
	vm.assert(ok)
	uid, ok := pop(vm).(int64)
	vm.assert(ok)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	err := file.File.Chown(int(uid), int(gid))
	if err != nil {
		throw(err, vm)
		return 0
	}
	return 0
}

// arg1: File, arg2: mode
func builtinFChmod(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 2)
	mode, ok := pop(vm).(int64)
	vm.assert(ok)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	err := file.File.Chmod(os.FileMode(uint32(mode)))
	if err != nil {
		throw(err, vm)
		return 0
	}
	return 0
}

// arg1: path
func builtinChdir(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	path, ok := pop(vm).(string)
	vm.assert(ok)
	err := os.Chdir(path)
	if err != nil {
		throw(err, vm)
		return 0
	}
	return 0
}

// arg1: filepath, arg2: uid, arg3: gid
func builtinChown(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 3)
	gid, ok := pop(vm).(int64)
	vm.assert(ok)
	uid, ok := pop(vm).(int64)
	vm.assert(ok)
	path, ok := pop(vm).(string)
	vm.assert(ok)
	err := os.Chown(path, int(uid), int(gid))
	if err != nil {
		throw(err, vm)
		return 0
	}
	return 0
}

// arg1: filepath, arg2: mode
func builtinChmod(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 2)
	mode, ok := pop(vm).(int64)
	vm.assert(ok)
	file, ok := pop(vm).(string)
	vm.assert(ok)
	err := os.Chmod(file, os.FileMode(uint32(mode)))
	if err != nil {
		throw(err, vm)
		return 0
	}
	return 0
}

// arg1: File, arg2: offset, arg3: whence("cur","end","start")
// return: n
func builtinSeek(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 3)
	whenceS, ok := pop(vm).(string)
	vm.assert(ok)
	var whence int
	if whenceS == "cur" {
		whence = io.SeekCurrent
	} else if whenceS == "end" {
		whence = io.SeekEnd
	} else if whenceS == "start" {
		whence = io.SeekStart
	} else {
		vm.assert(false)
	}
	offset, ok := pop(vm).(int64)
	vm.assert(ok)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)

	n, err := file.File.Seek(offset, whence)
	if err != nil {
		throw(err, vm)
		return 0
	}
	push(vm, int64(n))
	return 1
}

// arg1: File
func builtinClose(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	if err := file.File.Close(); err != nil {
		throw(err, vm)
	}
	return 0
}

// arg1: File, arg2: Buffer or string, arg3: size
// return: n
func builtinWrite(argCnt int, vm *VM) (retCnt int) {
	var data []byte
	size, ok := pop(vm).(int64)
	vm.assert(ok)
	src := pop(vm)
	if buf, ok := src.(*types.Buffer); ok {
		data = buf.Data[:size]
	} else if str, ok := src.(string); ok {
		data = []byte(str)[:size]
	} else {
		vm.assert(false)
	}
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	n, err := file.File.Write(data)
	if err != nil {
		throw(err, vm)
		return 0
	}
	push(vm, int64(n))
	return 1
}

// arg1: File, arg2: Buffer, arg3: size
// return: n
func builtinRead(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 3)
	size, ok := pop(vm).(int64)
	vm.assert(ok)
	buff, ok := pop(vm).(*types.Buffer)
	vm.assert(ok)
	file, ok := pop(vm).(*types.File)
	vm.assert(ok)
	n, err := file.File.Read(buff.Data[:size])
	if err != nil {
		throw(err, vm)
		return 0
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
	vm.assert(argCnt == 3)
	mode, ok := pop(vm).(int64)
	vm.assert(ok)
	flagS, ok := pop(vm).(string)
	vm.assert(ok)
	filepath, ok := pop(vm).(string)
	vm.assert(ok)
	flag := getFileFlag(flagS)

	file, err := os.OpenFile(filepath, flag, os.FileMode(uint32(mode)))
	if err != nil {
		throw(err, vm)
		return 0
	}
	push(vm, types.NewFile(file))
	return 1
}

// arg1: string
// return: Buffer
func builtinBufferFrom(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	str, ok := pop(vm).(string)
	vm.assert(ok)
	push(vm, types.NewBufferFromString(str))
	return 1
}

// arg1: Buffer1, arg2: Buffer2, arg3: length, arg4: idx1, arg5: idx2
func builtinBufferCopy(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 5)
	idx2, ok := pop(vm).(int64)
	vm.assert(ok)
	idx1, ok := pop(vm).(int64)
	vm.assert(ok)
	length, ok := pop(vm).(int64)
	vm.assert(ok)
	buf2, ok := pop(vm).(*types.Buffer)
	vm.assert(ok)
	buf1, ok := pop(vm).(*types.Buffer)
	vm.assert(ok)
	copy(buf1.Data[idx1:], buf2.Data[idx2:idx2+length])
	return 0
}

// args: Buffer,Buffer[,Buffer]
// return: Buffer
func builtinBufferConcat(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt >= 2)
	bufs := make([][]byte, argCnt)
	var length int
	for i := argCnt - 1; i >= 0; i-- {
		buf, ok := pop(vm).(*types.Buffer)
		vm.assert(ok)
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
	vm.assert(argCnt == 3)
	length, ok := pop(vm).(int64)
	vm.assert(ok)
	offset, ok := pop(vm).(int64)
	vm.assert(ok)
	buffer, ok := pop(vm).(*types.Buffer)
	vm.assert(ok)
	push(vm, types.NewBuffer(buffer.Data[offset:offset+length]))
	return 1
}

// arg1: Buffer, arg2: offset, arg3: length
// return: string
func builtinBufferToString(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 3)
	length, ok := pop(vm).(int64)
	vm.assert(ok)
	offset, ok := pop(vm).(int64)
	vm.assert(ok)
	buffer, ok := pop(vm).(*types.Buffer)
	vm.assert(ok)
	push(vm, string(buffer.Data[offset:offset+length]))
	return 1
}

// arg1: Buffer, arg2: offset, arg3: size, arg4: littleEndian arg5: number
func builtinBufferWriteNumber(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 5)
	number := pop(vm)
	littleEndian, ok := pop(vm).(bool)
	vm.assert(ok)
	size, ok := pop(vm).(int64)
	vm.assert(ok)
	offset, ok := pop(vm).(int64)
	vm.assert(ok)
	buf, ok := pop(vm).(*types.Buffer)
	vm.assert(ok)

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
	vm.assert(argCnt == 6)
	isFloat, ok := pop(vm).(bool)
	vm.assert(ok)
	littleEndian, ok := pop(vm).(bool)
	vm.assert(ok)
	signed, ok := pop(vm).(bool)
	vm.assert(ok)
	size, ok := pop(vm).(int64)
	vm.assert(ok)
	offset, ok := pop(vm).(int64)
	vm.assert(ok)
	buf, ok := pop(vm).(*types.Buffer)
	vm.assert(ok)
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
		vm.assert(false)
	}
	push(vm, result)
	return 1
}

// arg1: capacity
// return: Buffer
func builtinBufferNew(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	capacity, ok := pop(vm).(int64)
	vm.assert(ok)
	push(vm, types.NewBufferN(int(capacity)))
	return 1
}

func builtinClone(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	switch src := pop(vm).(type) {
	case *types.Array:
		arr := make([]interface{}, len(src.Data))
		copy(arr, src.Data)
		push(vm, types.NewArray(arr))
	case *types.Object:
		obj := types.NewObjectN(src.KVCount())
		push(vm, obj.Clone())
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
	vm.assert(argCnt == 2)
	key := pop(vm)
	obj := pop(vm).(*types.Object)
	obj.Delete(key)
	return 0
}

func getType(v interface{}) string {
	switch v.(type) {
	case string:
		return "String"
	case *types.Closure:
		return "Closure"
	case *builtinFunc:
		return "Builtin"
	case *types.Object:
		return "Object"
	case *types.Array:
		return "Array"
	case *types.Buffer:
		return "Buffer"
	case int64, float64:
		return "Number"
	case bool:
		return "Boolean"
	case nil:
		return "Nil"
	}
	return ""
}

func builtinType(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	t := getType(pop(vm))
	push(vm, t)
	return 1
}

func builtinSub(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 2 || argCnt == 3)
	target := vm.curProto.stack.top(argCnt)
	start := vm.curProto.stack.top(argCnt - 1).(int64)
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
		vm.assert(false)
	}
	return 1
}

func builtinAppend(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt < 2)
	target := vm.curProto.stack.top(argCnt)
	arr, ok := target.(*types.Array)
	vm.assert(ok)
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
		fprint(os.Stdout, arg)
		fmt.Printf(" ")
	}
	vm.curProto.stack.popN(n)
	fmt.Println()
	return 0
}

func builtinLen(argCnt int, vm *VM) (retCnt int) {
	vm.assert(argCnt == 1)
	var length int64
	switch val := pop(vm).(type) {
	case *types.Array:
		length = int64(len(val.Data))
	case *types.Object:
		length = int64(val.KVCount())
	case string:
		length = int64(len(val))
	case *types.Buffer:
		length = int64(len(val.Data))
	default:
		vm.assert(false)
	}
	push(vm, length)
	return 1
}

func fprint(w io.Writer, val interface{}) {
	switch val := val.(type) {
	case *types.Closure:
		fmt.Fprintf(w, "<closure>")
	case *builtinFunc:
		fmt.Fprintf(w, "<builtin:\"%s\">", val.name)
	case string:
		fmt.Fprintf(w, "%s", val)
	case *types.Object:
		fmt.Fprintf(w, "Object{")
		i := 0
		cnt := val.KVCount()
		val.ForEach(func(k, v interface{}) {
			fprint(w, k)
			fmt.Fprintf(w, ": ")
			fprint(w, v)
			if i != cnt-1 {
				fmt.Fprintf(w, ", ")
				i++
			}
		})
		fmt.Fprintf(w, "}")
	case *types.Array:
		fmt.Fprintf(w, "Array[")
		for i, v := range val.Data {
			fprint(w, v)
			if i != len(val.Data)-1 {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprintf(w, "]")
	case *types.Buffer:
		fmt.Fprintf(w, "<Buffer>")
	case *types.File:
		fmt.Fprintf(w, "<File>")
	default:
		fmt.Fprintf(w, "%v", val)
	}
}

func push(vm *VM, val interface{}) {
	vm.curProto.stack.Push(val)
}

func pop(vm *VM) interface{} {
	return vm.curProto.stack.pop()
}
