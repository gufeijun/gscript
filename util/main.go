package main

import (
	"fmt"
	"gscript/compiler/codegen"
	"gscript/compiler/lexer"
	"gscript/compiler/parser"
	"gscript/proto"
	"gscript/std"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	codegen.SetStdLibGenMode()
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <dir>\n", os.Args[0])
		return
	}
	for lib := range std.StdLibMap {
		err := complieStdLib(lib)
		if err != nil {
			fmt.Printf("complie std libarary failed, error message: %v\n", err)
			return
		}
	}
}

func complieStdLib(stdlib string) error {
	source, target := path.Join(os.Args[1], stdlib+".gs"), path.Join(os.Args[1], stdlib+".gsproto")
	protoNum := std.StdLibMap[stdlib]

	code, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	parser := parser.NewParser(lexer.NewLexer(source, code))
	prog := parser.Parse()
	var imports []codegen.Import
	for _, _import := range prog.Imports {
		for _, lib := range _import.Libs {
			num, ok := std.StdLibMap[lib.Path]
			if !ok {
				return fmt.Errorf("invalid std libarary: %s", lib.Path)
			}
			if lib.Alias == "" {
				lib.Alias = lib.Path
			}
			imports = append(imports, codegen.Import{
				StdLib:      lib.Stdlib,
				Alias:       lib.Alias,
				ProtoNumber: num,
			})
		}
	}
	_proto := codegen.Gen(parser, prog, imports, protoNum)
	_proto.FilePath = stdlib

	file, err := os.OpenFile(target, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	return proto.WriteProto(file, &_proto)
}
