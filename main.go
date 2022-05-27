package main

import (
	"gscript/complier"
	"gscript/vm"
)

func main() {
	protos, err := complier.ComplieWithSrcFile("./test.gs")
	if err != nil {
		panic(err)
	}
	v := vm.NewVM(protos)
	v.Debug()
}
