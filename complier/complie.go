package complier

import (
	"fmt"
	"gscript/complier/codegen"
	"gscript/complier/lexer"
	"gscript/complier/parser"
	"gscript/proto"
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

func ComplieWithSrcFile(path string) (_proto []proto.Proto, err error) {
	code, err := readCode(path)
	if err != nil {
		return
	}
	return ComplieWithSrcCode(code, path)
}

func ComplieWithSrcCode(code []byte, filename string) ([]proto.Proto, error) {
	graph := newGraph()
	n := graph.insert(filename)
	if err := complie(code, n, graph); err != nil {
		return nil, err
	}
	if graph.hasCircle() {
		return nil, fmt.Errorf("import circle occurs")
	}
	protos := make([]proto.Proto, len(graph.nodes))
	for _, node := range graph.nodes {
		protos[node.protoNum] = *node.proto
	}
	return protos, nil
}

func complie(code []byte, n *node, graph *graph) error {
	parser := parser.NewParser(lexer.NewLexer(n.pathname, code))
	prog := parser.Parse()

	var imports []codegen.Import
	var nodes []*node
	for _, _import := range prog.Imports {
		for _, lib := range _import.Libs {
			if lib.Stdlib {
				continue // TODO
			}
			filepath := lib.Path + ".gs"
			nn := graph.insertPath(n.pathname, filepath)
			alias := lib.Alias
			if alias == "" {
				alias = path.Base(lib.Path)
			}
			nodes = append(nodes, nn)
			imports = append(imports, codegen.Import{
				ProtoNumber: nn.protoNum,
				Alias:       alias,
			})
		}
	}

	mainProto := n.protoNum == 0
	proto := codegen.Gen(parser, prog, imports, mainProto)
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
