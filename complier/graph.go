package complier

import (
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
		panic("") // TODO
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

func getImportPath(base string, _import string) string {
	curDir := abs(".")
	if err := os.Chdir(filepath.Dir(base)); err != nil {
		panic(err)
	}
	res := abs(_import)
	if err := os.Chdir(curDir); err != nil {
		panic(err)
	}
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

func (g *graph) getNode(path string) (n *node, ok bool) {
	n, ok = g.nodes[abs(path)]
	return
}
