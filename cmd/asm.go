package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gscript/proto"
	"gscript/std"
	"gscript/vm"
	"io"
	"os"
)

func WriteHumanReadableAsmToFile(target string, protos []proto.Proto) error {
	file, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	return WriteHumanReadableAsm(file, protos)
}

func WriteHumanReadableAsm(w io.Writer, protos []proto.Proto) error {
	for i := range protos {
		var buff bytes.Buffer
		if err := writeHumanReadableAsm(&buff, protos, i); err != nil {
			return err
		}
		if _, err := io.Copy(w, &buff); err != nil {
			return err
		}
	}
	return nil
}

func writeHumanReadableAsm(w *bytes.Buffer, protos []proto.Proto, idx int) error {
	p := protos[idx]
	if idx == 0 {
		fmt.Fprintf(w, "%s(MainProto):\n", p.FilePath)
	} else {
		fmt.Fprintf(w, "%s(%d):\n", p.FilePath, idx)
	}
	fmt.Fprintf(w, "\tMainFunc:\n")
	if err := writeAsm(w, p.Text, protos); err != nil {
		return err
	}
	if err := writeFuncs(w, p.Funcs, protos); err != nil {
		return err
	}
	if err := writeAnonymousFuncs(w, p.AnonymousFuncs, protos); err != nil {
		return err
	}
	// should write constant table?
	return nil
}

func writeAnonymousFuncs(w *bytes.Buffer, funcs []proto.AnonymousFuncProto, protos []proto.Proto) error {
	for i := range funcs {
		fmt.Fprintf(w, "\t%dth AnonymousFunc:\n", i)
		if err := writeAsm(w, funcs[i].Info.Text, protos); err != nil {
			return err
		}
	}
	return nil
}

func writeFuncs(w *bytes.Buffer, funcs []proto.FuncProto, protos []proto.Proto) error {
	for i := range funcs {
		fmt.Fprintf(w, "\t%dth NamedFunc \"%s\":\n", i, funcs[i].Name)
		if err := writeAsm(w, funcs[i].Info.Text, protos); err != nil {
			return err
		}
	}
	return nil
}

func getUint32(pc *int, text []byte) uint32 {
	num := binary.LittleEndian.Uint32(text[*pc:])
	*pc += 4
	return num
}

func writeEscapeString(w *bytes.Buffer, str string) {
	w.WriteByte('"')
	for _, ch := range str {
		if ch == '\r' {
			w.WriteString("\\r")
		} else if ch == '\n' {
			w.WriteString("\\n")
		} else if ch == '\t' {
			w.WriteString("\\t")
		} else {
			w.WriteRune(ch)
		}
	}
	w.WriteByte('"')
}

func writeAsm(w *bytes.Buffer, text []byte, protos []proto.Proto) error {
	getFuncName := func(protoNum uint32, idx uint32) string {
		return protos[protoNum].Funcs[idx].Name
	}
	getFileName := func(protoNum uint32) string {
		return protos[protoNum].FilePath
	}
	var pc int
	for {
		if pc >= len(text) {
			break
		}
		instruction := text[pc]
		pc++
		fmt.Fprintf(w, "\t\t")
		if int(instruction) < len(zereOpNumAsms) {
			fmt.Fprintf(w, "%s\n", zereOpNumAsms[instruction])
			continue
		}
		switch instruction {
		case proto.INS_LOAD_CONST:
			protoNum, idx := getUint32(&pc, text), getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_CONST ")
			cons := protos[protoNum].Consts[idx]
			if str, ok := cons.(string); ok {
				writeEscapeString(w, str)
			} else {
				fmt.Fprintf(w, "%v", cons)
			}
		case proto.INS_LOAD_NAME:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_NAME %d", idx)
		case proto.INS_LOAD_FUNC:
			protoNum, idx := getUint32(&pc, text), getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_FUNC \"%s\"", getFuncName(protoNum, idx))
		case proto.INS_LOAD_BUILTIN:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_BUILTIN \"%s\"", vm.GetBuiltinFuncNameByNum(idx))
		case proto.INS_LOAD_ANONYMOUS:
			_, idx := getUint32(&pc, text), getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_ANONYMOUS %d", idx)
		case proto.INS_LOAD_UPVALUE:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_UPVALUE %d", idx)
		case proto.INS_LOAD_PROTO:
			protoNum := getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_PROTO %d(%s)", protoNum, getFileName(protoNum))
		case proto.INS_LOAD_STDLIB:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "LOAD_STDLIB %s", std.GetLibNameByProtoNum(idx))
		case proto.INS_STORE_NAME:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "STORE_NAME %d", idx)
		case proto.INS_STORE_UPVALUE:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "STORE_UPVALUE %d", idx)
		case proto.INS_RESIZE_NAMETABLE:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "RESIZE_TABLE %d", idx)
		case proto.INS_SLICE_NEW:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "SLICE_NEW %d", idx)
		case proto.INS_NEW_MAP:
			idx := getUint32(&pc, text)
			fmt.Fprintf(w, "MAP_NEW %d", idx)
		case proto.INS_JUMP_REL:
			steps := getUint32(&pc, text)
			fmt.Fprintf(w, "JUMP %d", pc+int(steps))
		case proto.INS_JUMP_ABS:
			addr := getUint32(&pc, text)
			fmt.Fprintf(w, "JUMP %d", addr)
		case proto.INS_JUMP_IF:
			steps := getUint32(&pc, text)
			fmt.Fprintf(w, "JUMP_IF %d", pc+int(steps))
		case proto.INS_JUMP_LAND:
			steps := getUint32(&pc, text)
			fmt.Fprintf(w, "JUMP_LAND %d", pc+int(steps))
		case proto.INS_JUMP_LOR:
			steps := getUint32(&pc, text)
			fmt.Fprintf(w, "JUMP_LOR %d", pc+int(steps))
		case proto.INS_JUMP_CASE:
			steps := getUint32(&pc, text)
			fmt.Fprintf(w, "JUMP_CASE %d", pc+int(steps))
		case proto.INS_CALL:
			retCnt, argCnt := text[pc], text[pc+1]
			pc += 2
			fmt.Fprintf(w, "CALL %d %d", retCnt, argCnt)
		case proto.INS_RETURN:
			retCnt := getUint32(&pc, text)
			fmt.Fprintf(w, "RETURN %d", retCnt)
		case proto.INS_TRY:
			steps := getUint32(&pc, text)
			fmt.Fprintf(w, "TRY %d", pc+int(steps))
		default:
			return fmt.Errorf("invalid instruction code: %d", instruction)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
	return nil
}

var zereOpNumAsms = []string{
	"NOT", "NEG", "LNOT", "ADD", "SUB", "MUL", "DIV", "MOD", "AND",
	"XOR", "OR", "IDIV", "SHR", "SHL", "LE", "GE", "LT", "GT", "EQ",
	"NE", "LAND", "LOR", "ATTR", "LOAD_NIL", "STORE_KV", "PUSH_NIL",
	"PUSH_NAME", "COPY_STACK_TOP", "POP_TOP", "STOP", "ATTR=", "ATTR+=",
	"ATTR-=", "ATTR*=", "ATTR/=", "ATTR%=", "ATTR&=", "ATTR^=", "ATTR|=",
	"ATTR_ACCESS", "ROT_TWO", "EXPORT", "END_TRY", "NEW_EMPTY_MAP",
}
