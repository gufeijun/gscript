package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"gscript/complier/ast"
	"io"
	"math"
	"os"
)

const (
	magicNumber       = 0x00686a6c
	VersionMajor byte = 0
	VersionMinor byte = 1
)

const (
	typeString  = 0
	typeInteger = 1
	typeFloat   = 2
	typeTrue    = 3
	typeFalse   = 4
	typeNil     = 5
)

type Proto struct {
	FilePath       string
	Consts         []interface{}
	Funcs          []FuncProto
	AnonymousFuncs []AnonymousFuncProto
	Text           []byte
}

type Header struct {
	Magic        uint32
	VersionMajor byte
	VersionMinor byte
}

func writeHeader(w io.Writer) error {
	var buff bytes.Buffer
	writeUint32(&buff, magicNumber)
	writeVersion(&buff)
	_, err := io.Copy(w, &buff)
	return err
}

func readHeader(r *bufio.Reader) (h Header, err error) {
	magic, err := readUint32(r)
	if err != nil {
		return
	}
	if magic != magicNumber {
		return h, fmt.Errorf("invalid magic number")
	}
	h.Magic = magic
	h.VersionMajor, err = r.ReadByte()
	if err != nil {
		return
	}
	h.VersionMinor, err = r.ReadByte()
	return
}

func writeVersion(buff *bytes.Buffer) {
	buff.WriteByte(VersionMajor)
	buff.WriteByte(VersionMinor)
}

func WriteProtos(w io.Writer, protos []Proto) error {
	if err := writeHeader(w); err != nil {
		return err
	}
	var buff bytes.Buffer
	writeUint32(&buff, uint32(len(protos)))
	for i := range protos {
		writeProto(&buff, &protos[i])
	}
	_, err := io.Copy(w, &buff)
	return err
}

func WriteProto(w io.Writer, proto *Proto) error {
	var buff bytes.Buffer
	writeProto(&buff, proto)
	_, err := io.Copy(w, &buff)
	return err
}

func writeProto(w *bytes.Buffer, proto *Proto) {
	writeString(w, proto.FilePath)
	// for i, cons := range proto.Consts {
	// 	fmt.Println(i, cons)
	// }
	writeConsts(w, proto.Consts)
	writeFuncs(w, proto.Funcs)
	writeAnonymousFuncs(w, proto.AnonymousFuncs)
	writeText(w, proto.Text)
}

// [funcsCnt:4B] [funcs]
// func:	[UpValueCnt:4B] [UpValues] [BasicInfo]
// UpValue: [DirectDependent:1B] [Index:4B]
func writeAnonymousFuncs(w *bytes.Buffer, funcs []AnonymousFuncProto) {
	writeUint32(w, uint32(len(funcs)))
	for _, _func := range funcs {
		writeUint32(w, uint32(len(_func.UpValues)))
		for _, upvalue := range _func.UpValues {
			if upvalue.DirectDependent {
				w.WriteByte(1)
			} else {
				w.WriteByte(0)
			}
			writeUint32(w, upvalue.Index)
		}
		writeBasicInfo(w, _func.Info)
	}
}

func readAnonymousFuncs(r *bufio.Reader) (funcs []AnonymousFuncProto, err error) {
	funcCnt, err := readUint32(r)
	if err != nil {
		return
	}
	funcs = make([]AnonymousFuncProto, 0, funcCnt)
	for i := 0; i < int(funcCnt); i++ {
		var _func AnonymousFuncProto
		upvalueCnt, err := readUint32(r)
		if err != nil {
			return nil, err
		}
		for j := 0; j < int(upvalueCnt); j++ {
			var upvalue UpValuePtr
			b, err := r.ReadByte()
			if err != nil {
				return funcs, err
			}
			upvalue.DirectDependent = b == 1
			upvalue.Index, err = readUint32(r)
			if err != nil {
				return funcs, err
			}
			_func.UpValues = append(_func.UpValues, upvalue)
		}
		_func.Info, err = readBasicInfo(r)
		if err != nil {
			return funcs, err
		}
		funcs = append(funcs, _func)
	}
	return
}

// [funcsCnt:4B] [funcs]
// func: 		[funcName] [upValueCnt:4B] [UpValues] [BasicInfo]
// UpValue:		[uint32:4B]
func writeFuncs(w *bytes.Buffer, funcs []FuncProto) {
	writeUint32(w, uint32(len(funcs)))
	for _, _func := range funcs {
		writeString(w, _func.Name)
		writeUint32(w, uint32(len(_func.UpValues)))
		for _, upvalue := range _func.UpValues {
			writeUint32(w, upvalue)
		}
		writeBasicInfo(w, _func.Info)
	}
}

func readFunc(r *bufio.Reader) (_func FuncProto, err error) {
	_func.Name, err = readString(r)
	if err != nil {
		return
	}
	upvalueCnt, err := readUint32(r)
	if err != nil {
		return
	}
	for i := 0; i < int(upvalueCnt); i++ {
		idx, err := readUint32(r)
		if err != nil {
			return _func, err
		}
		_func.UpValues = append(_func.UpValues, idx)
	}
	_func.Info, err = readBasicInfo(r)
	return
}

func readFuncs(r *bufio.Reader) (funcs []FuncProto, err error) {
	funcCnt, err := readUint32(r)
	if err != nil {
		return nil, err
	}
	funcs = make([]FuncProto, 0, funcCnt)
	for i := 0; i < int(funcCnt); i++ {
		_func, err := readFunc(r)
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, _func)
	}
	return
}

// BaiscInfo:	[VaArgs:1B] [parametersCnt:4B] [parameters] [textLen:4B] [text:[]byte]
// parameter:	[nameLength:4B] [name:string] [Default:consts]
func writeBasicInfo(w *bytes.Buffer, info *BasicInfo) {
	if info.VaArgs {
		w.WriteByte(1)
	} else {
		w.WriteByte(0)
	}
	writeUint32(w, uint32(len(info.Parameters)))
	for _, par := range info.Parameters {
		writeString(w, par.Name)
		writeConst(w, par.Default)
	}
	writeText(w, info.Text)
}

func readBasicInfo(r *bufio.Reader) (info *BasicInfo, err error) {
	VaArgs, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	info = new(BasicInfo)
	info.VaArgs = VaArgs == 1
	parCnt, err := readUint32(r)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(parCnt); i++ {
		name, err := readString(r)
		if err != nil {
			return nil, err
		}
		c, err := readConst(r)
		if err != nil {
			return nil, err
		}
		info.Parameters = append(info.Parameters, ast.Parameter{
			Name:    name,
			Default: c,
		})
	}
	info.Text, err = readText(r)
	return
}

func writeText(w *bytes.Buffer, text []byte) {
	writeUint32(w, uint32(len(text)))
	w.Write(text)
}

func readText(r *bufio.Reader) ([]byte, error) {
	length, err := readUint32(r)
	if err != nil {
		return nil, err
	}
	text := make([]byte, length)
	_, err = io.ReadFull(r, text)
	return text, err
}

// Write count(4B) of constants at first, then write all constants into w,
// every kind of constant will be organized like following:
// string: [type:1B](0) [length:4B] [values:len(string)]
// int64:  [type:1B](1) [values:8B]
// float64:[type:1B](2) [values:8B]
// true:   [type:1B](3)
// false:  [type:1B](4)
func writeConsts(w *bytes.Buffer, consts []interface{}) {
	writeUint32(w, uint32(len(consts)))
	for _, constant := range consts {
		writeConst(w, constant)
	}
}

func readConsts(r *bufio.Reader) (consts []interface{}, err error) {
	length, err := readUint32(r)
	if err != nil {
		return
	}
	consts = make([]interface{}, 0, length)
	for i := 0; i < int(length); i++ {
		constant, err := readConst(r)
		if err != nil {
			return nil, err
		}
		consts = append(consts, constant)
	}
	return
}

func writeConst(w *bytes.Buffer, constant interface{}) {
	if constant == nil {
		w.WriteByte(typeNil)
		return
	}
	switch val := constant.(type) {
	case string:
		w.WriteByte(typeString)
		writeString(w, val)
	case int64:
		w.WriteByte(typeInteger)
		writeUint64(w, uint64(val))
	case float64:
		w.WriteByte(typeFloat)
		writeFloat64(w, val)
	case bool:
		if val {
			w.WriteByte(typeTrue)
		} else {
			w.WriteByte(typeFalse)
		}
	default:
		fmt.Printf("writing invalid constant type\n")
		os.Exit(0)
	}
}

func readConst(r *bufio.Reader) (interface{}, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch t {
	case typeString:
		str, err := readString(r)
		return str, err
	case typeInteger:
		num, err := readUint64(r)
		return int64(num), err
	case typeFloat:
		f, err := readFloat64(r)
		return f, err
	case typeTrue:
		return true, nil
	case typeFalse:
		return false, nil
	case typeNil:
		return nil, nil
	default:
		return nil, fmt.Errorf("reading invalid constant type")
	}
}

func writeFloat64(w *bytes.Buffer, src float64) {
	writeUint64(w, math.Float64bits(src))
}

func readFloat64(r *bufio.Reader) (float64, error) {
	u64, err := readUint64(r)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(u64), nil
}

func writeUint64(w *bytes.Buffer, src uint64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, src)
	w.Write(data)
}

func readUint64(r *bufio.Reader) (uint64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

func writeUint32(w *bytes.Buffer, src uint32) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, src)
	w.Write(data)
}

func readUint32(r *bufio.Reader) (uint32, error) {
	data := make([]byte, 4)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return 0, err
	}
	v := binary.LittleEndian.Uint32(data)
	return v, nil
}

// length:4B data:string
func writeString(w *bytes.Buffer, src string) {
	writeUint32(w, uint32(len(src)))
	w.WriteString(src)
}

func readString(r *bufio.Reader) (string, error) {
	length, err := readUint32(r)
	if err != nil {
		return "", nil
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return "", nil
	}
	return string(data), nil
}

func ReadProtos(r io.Reader) (h Header, protos []Proto, err error) {
	bufr := bufio.NewReader(r)
	h, err = readHeader(bufr)
	if err != nil {
		return
	}
	protos, err = readProtos(bufr)
	return
}

func ReadProto(r *bufio.Reader) (p Proto, err error) {
	p.FilePath, err = readString(r)
	if err != nil {
		return
	}
	p.Consts, err = readConsts(r)
	if err != nil {
		return
	}
	p.Funcs, err = readFuncs(r)
	if err != nil {
		return
	}
	p.AnonymousFuncs, err = readAnonymousFuncs(r)
	if err != nil {
		return
	}
	p.Text, err = readText(r)
	return
}

func readProtos(r *bufio.Reader) (protos []Proto, err error) {
	protoCnt, err := readUint32(r)
	if err != nil {
		return
	}
	protos = make([]Proto, 0, protoCnt)
	for i := 0; i < int(protoCnt); i++ {
		proto, err := ReadProto(r)
		if err != nil {
			return nil, err
		}
		protos = append(protos, proto)
	}
	return
}

func WriteProtosToFile(src string, protos []Proto) error {
	file, err := os.OpenFile(src, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	return WriteProtos(file, protos)
}

func IsProtoFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	if stat.Size() < 4 {
		return false
	}
	buf := make([]byte, 4)
	if _, err := io.ReadFull(file, buf); err != nil {
		return false
	}
	magic := binary.LittleEndian.Uint32(buf)
	return magic == magicNumber
}

func ReadProtosFromFile(src string) (Header, []Proto, error) {
	file, err := os.Open(src)
	if err != nil {
		return Header{}, nil, err
	}
	defer file.Close()
	return ReadProtos(file)
}
