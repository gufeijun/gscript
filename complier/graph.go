package complier

import (
	"fmt"
	"gscript/proto"
	"os"
	"path/filepath"
)

const (
	white = iota
	grey
	black
)

type node struct {
	color    uint32
	protoNum uint32
	proto    *proto.Proto
	pathname string
	imports  []*node
}

type graph struct {
	nodes map[string]*node
}

func newGraph() *graph {
	return &graph{
		nodes: make(map[string]*node),
	}
}

func abs(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("get absolute filepath of '%s' failed\n", path)
		fmt.Printf("\terror message: %v\n", err)
		os.Exit(0)
	}
	return path
}

func (g *graph) sortProtos() []proto.Proto {
	protos := make([]proto.Proto, len(g.nodes))
	for _, node := range g.nodes {
		protos[node.protoNum] = *node.proto
	}
	return protos
}

func (g *graph) hasCircle() bool {
	var n *node
	for _, n = range g.nodes {
		if n.color != white {
			continue
		}
		if hasCircle(n) {
			return true
		}
	}
	return false
}

func hasCircle(n *node) bool {
	n.color = grey
	for _, _import := range n.imports {
		if _import.color == grey {
			return true
		}
		if hasCircle(_import) {
			return true
		}
	}
	n.color = black
	return false
}

func chdir(dir string) {
	err := os.Chdir(dir)
	if err == nil {
		return
	}
	fmt.Printf("change work dir failed\n")
	fmt.Printf("\terror message: %v\n", err)
	os.Exit(0)
}

func getImportPath(base string, _import string) string {
	curDir := abs(".")
	chdir(filepath.Dir(base))
	res := abs(_import)
	chdir(curDir)
	return res
}

func (g *graph) insertPath(from, to string) *node {
	from = abs(from)
	to = getImportPath(from, to)
	n := g.nodes[from]
	nn, ok := g.nodes[to]
	if !ok {
		nn = &node{pathname: to, protoNum: uint32(len(g.nodes))}
		g.nodes[to] = nn
	}
	n.imports = append(n.imports, nn)
	return nn
}

func (g *graph) insert(path string) *node {
	path = abs(path)
	n := &node{pathname: path, protoNum: uint32(len(g.nodes))}
	g.nodes[path] = n
	return n
}
