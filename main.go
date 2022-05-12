package main

import (
	"gscript/complier/codegen"
	"gscript/complier/lexer"
	"gscript/complier/parser"
	"gscript/vm"
)

const code = `
class people{
	__self(age,name){
		this.age = age
		this.name = name
	}
	birthday(){
		this.age++
	}
	location = "wuhan"
}

let b = new people(10,"Mike")
b.birthday()  # b.age=11, b.name="Mike", b.location="wuhan"
`

func main() {
	parser := parser.NewParser(lexer.NewLexer("", []byte(code)))
	_proto := codegen.Gen(parser)

	v := vm.NewVM(&_proto)
	v.Debug()
}
