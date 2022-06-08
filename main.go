package main

import (
	"gscript/complier"
	"gscript/std"
	"gscript/vm"
)

func main() {
	protos, err := complier.ComplieWithSrcFile("./test.gs")
	if err != nil {
		panic(err)
	}

	stdlibs, err := std.ReadProtos()
	if err != nil {
		panic(err)
	}

	v := vm.NewVM(protos, stdlibs)
	v.Run()
	// v.Debug()
}
