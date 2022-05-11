package main

import (
	"gscript/complier/codegen"
	"gscript/complier/lexer"
	"gscript/complier/parser"
	"gscript/vm"
)

const s = `
func newCounter(){
	let count=0
	return func(){
		count++
		return func(){
			count++
			return count
		}
	}
}

let cb = newCounter()
let a = cb()	# count=1
let b= cb()		# count=2

let c = a()		# c=3 count=3
let d = b()		# d=4 count=4
`

const str = `
let count = 0;

let result = inc()	// result = count = 1

# closure
func inc(){
	count++
	return count
}
`

func main() {
	parser := parser.NewParser(lexer.NewLexer("", []byte(s)))
	_proto := codegen.Gen(parser)

	v := vm.NewVM(&_proto)
	v.Debug()
}
