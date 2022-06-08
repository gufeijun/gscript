package complier

import (
	"fmt"
	"gscript/complier/codegen"
	"gscript/complier/lexer"
	"gscript/complier/parser"
	"gscript/proto"
	"gscript/std"
	"io/ioutil"
	"os"
	"path"
)

const maxLimitFileSize = 10 << 20 // 10MB

func readCode(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.Size() > maxLimitFileSize {
		return nil, fmt.Errorf("%s: file size exceeds 10MB limit", path)
	}
	return ioutil.ReadFile(path)
}

func ComplieWithSrcFile(path string) (protos []proto.Proto, err error) {
	code, err := readCode(path)
	if err != nil {
		return
	}
	return ComplieWithSrcCode(code, path)
}

func ComplieWithSrcCode(code []byte, filename string) (protos []proto.Proto, err error) {
	graph := newGraph()
	n := graph.insert(filename)
	if err = complie(code, n, graph); err != nil {
		return
	}
	if graph.hasCircle() {
		return nil, fmt.Errorf("import circle occurs")
	}
	return graph.sortProtos(), nil
}

func complie(code []byte, n *node, graph *graph) error {
	parser := parser.NewParser(lexer.NewLexer(n.pathname, code))
	prog := parser.Parse()

	var imports []codegen.Import
	var nodes []*node
	for _, _import := range prog.Imports {
		for _, lib := range _import.Libs {
			var protoNumber uint32
			if lib.Stdlib {
				var ok bool
				protoNumber, ok = std.StdLibs[lib.Path]
				if !ok {
					return fmt.Errorf("invalid std libarary: %s", lib.Path)
				}
			} else {
				nn := graph.insertPath(n.pathname, lib.Path+".gs")
				nodes = append(nodes, nn)
				protoNumber = nn.protoNum
			}
			alias := lib.Alias
			if alias == "" {
				alias = path.Base(lib.Path)
			}
			imports = append(imports, codegen.Import{
				ProtoNumber: protoNumber,
				Alias:       alias,
				StdLib:      lib.Stdlib,
			})
		}
	}

	proto := codegen.Gen(parser, prog, imports, n.protoNum)
	proto.FilePath = n.pathname
	n.proto = &proto

	for _, n = range nodes {
		// if file has been already complied
		if n.proto != nil {
			continue
		}
		code, err := readCode(n.pathname)
		if err != nil {
			return err
		}
		if err := complie(code, n, graph); err != nil {
			return err
		}
	}
	return nil
}
