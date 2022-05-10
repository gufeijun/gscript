package proto

import "gscript/complier/ast"

type Closure struct {
	Addr         uint32
	Parameters   []ast.Parameter
	VaArgs       bool
	UpValues     [][2]uint32
	UpValueTable []*interface{}
}
