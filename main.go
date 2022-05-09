package main

import (
	"gscript/complier/codegen"
	"gscript/complier/lexer"
	"gscript/complier/parser"
	"gscript/vm"
)

const str = `
func fib(a){
	if (a==0) return 0;
	if (a==1) return 1;
	return fib(a-1)+fib(a-2)
}

let a=fib(35)
`

func main() {
	parser := parser.NewParser(lexer.NewLexer("", []byte(str)))
	text, consts, ft := codegen.Gen(parser)

	v := vm.NewVM(text, consts, ft)
	v.Debug()
}
