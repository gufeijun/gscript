package proto

import "gscript/complier/ast"

type BasicInfo struct {
	VaArgs     bool
	Addr       uint32
	Parameters []ast.Parameter
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