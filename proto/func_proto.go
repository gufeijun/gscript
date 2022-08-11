package proto

import "gscript/compiler/ast"

type BasicInfo struct {
	VaArgs     bool
	Parameters []ast.Parameter
	Text       []byte
}

type FuncProto struct {
	Info         *BasicInfo
	Name         string
	UpValues     []uint32
	UpValueTable interface{}
}

type UpValuePtr struct {
	DirectDependent bool
	Index           uint32
}

type AnonymousFuncProto struct {
	Info     *BasicInfo
	UpValues []UpValuePtr
}
