package parser

import (
	. "gscript/compiler/ast"
	"reflect"
	"testing"
)

func TestParseImports(t *testing.T) {
	srcs := []string{
		``,
		`
import net,http;
import system as sys`,
		`
import "../sum"
import "../crawler",socks as craw,s`,
	}
	wants := [][]Import{
		nil,
		[]Import{{2, []Lib{{true, "net", ""}, {true, "http", ""}}}, {3, []Lib{{true, "system", "sys"}}}},
		[]Import{{2, []Lib{{false, "../sum", ""}}}, {3, []Lib{{false, "../crawler", "craw"}, {true, "socks", "s"}}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		imports := NewParser(l).parseImports()
		if !reflect.DeepEqual(imports, wants[i]) {
			t.Fatalf("parseImports failed:\n%s\n", src)
		}
	}
}
