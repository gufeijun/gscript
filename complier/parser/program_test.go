package parser

import (
	. "gscript/complier/ast"
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
		[]Import{{[]Lib{{true, "net", ""}, {true, "http", ""}}}, {[]Lib{{true, "system", "sys"}}}},
		[]Import{{[]Lib{{false, "../sum", ""}}}, {[]Lib{{false, "../crawler", "craw"}, {true, "socks", "s"}}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		imports := parseImports(l)
		if !reflect.DeepEqual(imports, wants[i]) {
			t.Fatalf("parseImports failed:\n%s\n", src)
		}
	}
}
