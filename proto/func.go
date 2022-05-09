package proto

import "gscript/complier/ast"

type Func struct {
	Addr       uint32
	Parameters []ast.Parameter
	VaArgs     bool
}
