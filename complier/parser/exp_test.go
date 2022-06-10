package parser

import (
	. "gscript/complier/ast"
	"gscript/complier/lexer"
	"reflect"
	"testing"
)

func newLexer(src string) *lexer.Lexer {
	return lexer.NewLexer("", []byte(src))
}

func TestParseTerm12(t *testing.T) {
	srcs := []string{
		`a||c?b+d:c`,
		`a?(b?c:d):e`,
	}
	var wants = []*TernaryOpExp{
		{&BinOpExp{BINOP_LOR, &NameExp{1, "a"}, &NameExp{1, "c"}}, &BinOpExp{BINOP_ADD, &NameExp{1, "b"}, &NameExp{1, "d"}}, &NameExp{1, "c"}},
		{&NameExp{1, "a"}, &TernaryOpExp{&NameExp{1, "b"}, &NameExp{1, "c"}, &NameExp{1, "d"}}, &NameExp{1, "e"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm12(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term12 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm11(t *testing.T) {
	srcs := []string{
		`a||b||c`,
		`a||b&&c`,
	}
	var wants = []*BinOpExp{
		{BINOP_LOR, &BinOpExp{BINOP_LOR, &NameExp{1, "a"}, &NameExp{1, "b"}}, &NameExp{1, "c"}},
		{BINOP_LOR, &NameExp{1, "a"}, &BinOpExp{BINOP_LAND, &NameExp{1, "b"}, &NameExp{1, "c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm11(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term11 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm10(t *testing.T) {
	srcs := []string{
		`a&&b`,
		`a&&b|c`,
	}
	var wants = []*BinOpExp{
		{BINOP_LAND, &NameExp{1, "a"}, &NameExp{1, "b"}},
		{BINOP_LAND, &NameExp{1, "a"}, &BinOpExp{BINOP_OR, &NameExp{1, "b"}, &NameExp{1, "c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm10(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term10 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm9(t *testing.T) {
	srcs := []string{
		`a|b`,
		`a|b^c`,
	}
	var wants = []*BinOpExp{
		{BINOP_OR, &NameExp{1, "a"}, &NameExp{1, "b"}},
		{BINOP_OR, &NameExp{1, "a"}, &BinOpExp{BINOP_XOR, &NameExp{1, "b"}, &NameExp{1, "c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm9(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term9 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm8(t *testing.T) {
	srcs := []string{
		`a^b`,
		`a^b&c`,
	}
	var wants = []*BinOpExp{
		{BINOP_XOR, &NameExp{1, "a"}, &NameExp{1, "b"}},
		{BINOP_XOR, &NameExp{1, "a"}, &BinOpExp{BINOP_AND, &NameExp{1, "b"}, &NameExp{1, "c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm8(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term8 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm7(t *testing.T) {
	srcs := []string{
		`a&b`,
		`a&b==c`,
	}
	var wants = []*BinOpExp{
		{BINOP_AND, &NameExp{1, "a"}, &NameExp{1, "b"}},
		{BINOP_AND, &NameExp{1, "a"}, &BinOpExp{BINOP_EQ, &NameExp{1, "b"}, &NameExp{1, "c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm7(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term7 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm6(t *testing.T) {
	srcs := []string{
		`a==b`,
		`a!=b==c`,
		`a==b>=c`,
	}
	var wants = []*BinOpExp{
		{BINOP_EQ, &NameExp{1, "a"}, &NameExp{1, "b"}},
		{BINOP_EQ, &BinOpExp{BINOP_NE, &NameExp{1, "a"}, &NameExp{1, "b"}}, &NameExp{1, "c"}},
		{BINOP_EQ, &NameExp{1, "a"}, &BinOpExp{BINOP_GE, &NameExp{1, "b"}, &NameExp{1, "c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm6(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term6 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm5(t *testing.T) {
	srcs := []string{
		`a > b<<c`,
		`1<=3`,
		`2<4`,
		`a>=8>9`,
	}
	var wants = []*BinOpExp{
		{BINOP_GT, &NameExp{1, "a"}, &BinOpExp{BINOP_SHL, &NameExp{1, "b"}, &NameExp{1, "c"}}},
		{BINOP_LE, &NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(3)}},
		{BINOP_LT, &NumberLiteralExp{int64(2)}, &NumberLiteralExp{int64(4)}},
		{BINOP_GT, &BinOpExp{BINOP_GE, &NameExp{1, "a"}, &NumberLiteralExp{int64(8)}}, &NumberLiteralExp{int64(9)}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm5(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term5 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm4(t *testing.T) {
	srcs := []string{
		`a<<b+c`,
		`a>>b<<d`,
	}
	var wants = []*BinOpExp{
		{BINOP_SHL, &NameExp{1, "a"}, &BinOpExp{BINOP_ADD, &NameExp{1, "b"}, &NameExp{1, "c"}}},
		{BINOP_SHL, &BinOpExp{BINOP_SHR, &NameExp{1, "a"}, &NameExp{1, "b"}}, &NameExp{1, "d"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm4(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term4 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm3(t *testing.T) {
	srcs := []string{
		`a+b`,
		`a.b - c`,
		`a-b*c`,
		`(a-b)*c`,
	}
	var wants = []*BinOpExp{
		{BINOP_ADD, &NameExp{1, "a"}, &NameExp{1, "b"}},
		{BINOP_SUB, &BinOpExp{BINOP_ATTR, &NameExp{1, "a"}, &StringLiteralExp{"b"}}, &NameExp{1, "c"}},
		{BINOP_SUB, &NameExp{1, "a"}, &BinOpExp{BINOP_MUL, &NameExp{1, "b"}, &NameExp{1, "c"}}},
		{BINOP_MUL, &BinOpExp{BINOP_SUB, &NameExp{1, "a"}, &NameExp{1, "b"}}, &NameExp{1, "c"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm3(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term3 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm2(t *testing.T) {
	srcs := []string{
		`-a/b`,
		`~arr[0]*100`,
		`arr[0]++//20`,
		`--arr[0]%p`,
		`a/b*c`,
	}
	var wants = []*BinOpExp{
		{BINOP_DIV, &UnOpExp{UNOP_NEG, &NameExp{1, "a"}}, &NameExp{1, "b"}},
		{BINOP_MUL, &UnOpExp{UNOP_NOT, &BinOpExp{BINOP_ATTR, &NameExp{1, "arr"}, &NumberLiteralExp{int64(0)}}}, &NumberLiteralExp{int64(100)}},
		{BINOP_IDIV, &UnOpExp{UNOP_INC_, &BinOpExp{BINOP_ATTR, &NameExp{1, "arr"}, &NumberLiteralExp{int64(0)}}}, &NumberLiteralExp{int64(20)}},
		{BINOP_MOD, &UnOpExp{UNOP_DEC, &BinOpExp{BINOP_ATTR, &NameExp{1, "arr"}, &NumberLiteralExp{int64(0)}}}, &NameExp{1, "p"}},
		{BINOP_MUL, &BinOpExp{BINOP_DIV, &NameExp{1, "a"}, &NameExp{1, "b"}}, &NameExp{1, "c"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm2(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term2 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm1(t *testing.T) {
	srcs := []string{
		`!(false)`,
		`~++a`,
		`-a--`,
		`-arr[1]`,
		`!!~1`,
	}
	var wants = []*UnOpExp{
		{UNOP_LNOT, &FalseExp{}},
		{UNOP_NOT, &UnOpExp{UNOP_INC, &NameExp{1, "a"}}},
		{UNOP_NEG, &UnOpExp{UNOP_DEC_, &NameExp{1, "a"}}},
		{UNOP_NEG, &BinOpExp{BINOP_ATTR, &NameExp{1, "arr"}, &NumberLiteralExp{int64(1)}}},
		{UNOP_LNOT, &UnOpExp{UNOP_LNOT, &UnOpExp{UNOP_NOT, &NumberLiteralExp{int64(1)}}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm1(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term1 failed:\n%s\n", src)
		}
	}
}

func TestParseTerm0(t *testing.T) {
	srcs := []string{
		`++a`,
		`--a[0]`,
		`a++`,
		`a.b--`,
	}
	var wants = []*UnOpExp{
		{UNOP_INC, &NameExp{1, "a"}},
		{UNOP_DEC, &BinOpExp{BINOP_ATTR, &NameExp{1, "a"}, &NumberLiteralExp{int64(0)}}},
		{UNOP_INC_, &NameExp{1, "a"}},
		{UNOP_DEC_, &BinOpExp{BINOP_ATTR, &NameExp{1, "a"}, &StringLiteralExp{"b"}}},
	}

	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm0(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse term0 failed:\n%s\n", src)
		}
	}
}

func TestFuncCallOrAttrExp(t *testing.T) {
	src := `m[sum(1,2)].f()[1]`
	want := &BinOpExp{
		BinOp: BINOP_ATTR,
		Exp1: &FuncCallExp{
			Func: &BinOpExp{
				BinOp: BINOP_ATTR,
				Exp1: &BinOpExp{
					BinOp: BINOP_ATTR,
					Exp1:  &NameExp{1, "m"},
					Exp2: &FuncCallExp{
						Func: &NameExp{1, "sum"},
						Args: []Exp{
							&NumberLiteralExp{int64(1)},
							&NumberLiteralExp{int64(2)},
						},
					},
				},
				Exp2: &StringLiteralExp{"f"},
			},
			Args: nil,
		},
		Exp2: &NumberLiteralExp{int64(1)},
	}
	l := newLexer(src)
	exp := parseFuncCallOrAttrExp(NewParser(l))
	if !reflect.DeepEqual(want, exp) {
		t.Fatalf("test FuncCallOrAttr failed:\n%s\n", src)
	}
}

// name(1)()("xxx") ...
func TestFuncCallExp(t *testing.T) {
	var srcs = []string{
		`sum(1,2)`,
		`sum()(1,2)`,
	}
	var wants = []*FuncCallExp{
		{&NameExp{1, "sum"}, []Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
		{&FuncCallExp{&NameExp{1, "sum"}, nil}, []Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseFuncCallOrAttrExp(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse function call failed:\n%s\n", src)
		}
	}
}

// obj["a"][1] ...
func TestParseMapAccessExp(t *testing.T) {
	var srcs = []string{
		`map["a"]`,
		`map["b"][1]`,
	}
	var wants = []*BinOpExp{
		{BINOP_ATTR, &NameExp{1, "map"}, &StringLiteralExp{"a"}},
		{BINOP_ATTR, &BinOpExp{
			BINOP_ATTR,
			&NameExp{1, "map"},
			&StringLiteralExp{"b"},
		}, &NumberLiteralExp{int64(1)}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseFuncCallOrAttrExp(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse map access failed:\n%s\n", src)
		}
	}
}

// obj.field1.field2 ...
func TestParseAttributeAccessExp(t *testing.T) {
	var srcs = []string{
		`person.name`,
		`person.father.name`,
	}
	var wants = []*BinOpExp{
		{BINOP_ATTR, &NameExp{1, "person"}, &StringLiteralExp{"name"}},
		{BINOP_ATTR, &BinOpExp{
			BINOP_ATTR,
			&NameExp{1, "person"},
			&StringLiteralExp{"father"},
		}, &StringLiteralExp{"name"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseFuncCallOrAttrExp(NewParser(l))
		if !reflect.DeepEqual(wants[i], exp) {
			t.Fatalf("parse attribute access failed:\n%s\n", src)
		}
	}
}

func TestParseArrLiteralExp(t *testing.T) {
	var arrLiterals = []string{
		`[]`,
		`[a,b]`,
		`[1,[1,2]]`,
		`["good"]`,
		`[{"a":a}]`,
		`[
	1,
	2
]`,
		`
[
	1,
	2,
]
`,
	}
	var wants = []*ArrLiteralExp{
		{nil},
		{[]Exp{&NameExp{1, "a"}, &NameExp{1, "b"}}},
		{[]Exp{&NumberLiteralExp{int64(1)}, &ArrLiteralExp{
			[]Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}},
		}}},
		{[]Exp{&StringLiteralExp{"good"}}},
		{[]Exp{&MapLiteralExp{
			Keys: []interface{}{"a"},
			Vals: []Exp{&NameExp{1, "a"}},
		}}},
		{[]Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
		{[]Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
	}
	for i, literal := range arrLiterals {
		l := newLexer(literal)
		exp := parseArrLiteralExp(NewParser(l))
		if !reflect.DeepEqual(exp, wants[i]) {
			t.Fatalf("parse array literal failed:\n%s\n", literal)
		}
	}
}

func TestParseMapLiteralExp(t *testing.T) {
	var mapLiterals = []string{
		`{}`,
		`{interger:1,string:"b",float:1.89}`,
		`{true:true,false:false}`,
		`{1:"1",3.14:"3.14",}`,
		`{"child":nil}`,
	}
	var wants = []*MapLiteralExp{
		{nil, nil},
		{[]interface{}{"interger", "string", "float"},
			[]Exp{&NumberLiteralExp{int64(1)}, &StringLiteralExp{"b"}, &NumberLiteralExp{float64(1.89)}}},
		{[]interface{}{true, false}, []Exp{&TrueExp{}, &FalseExp{}}},
		{[]interface{}{int64(1), float64(3.14)}, []Exp{&StringLiteralExp{"1"}, &StringLiteralExp{"3.14"}}},
		{[]interface{}{"child"}, []Exp{&NilExp{}}},
	}
	for i, literal := range mapLiterals {
		l := newLexer(literal)
		exp := parseMapLiteralExp(NewParser(l))
		if !reflect.DeepEqual(exp, wants[i]) {
			t.Fatalf("parse map literal failed:\n%s\n", literal)
		}
	}
}

func TestNewObjectExp(t *testing.T) {
	var srcs = []string{
		`new people`,
		`new people()`,
		`new people("jack",12)`,
		`new student(name)`,
		`new student(name,age)`,
	}
	var wants = []*NewObjectExp{
		{"people", nil},
		{"people", nil},
		{"people", []Exp{&StringLiteralExp{"jack"}, &NumberLiteralExp{int64(12)}}},
		{"student", []Exp{&NameExp{1, "name"}}},
		{"student", []Exp{&NameExp{1, "name"}, &NameExp{1, "age"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseNewObjectExp(NewParser(l))
		if !reflect.DeepEqual(exp, wants[i]) {
			t.Fatalf("parse new obejct failed:\n%s\n", src)
		}
	}
}

func TestParseFuncLiteralExp(t *testing.T) {
	var funcLiterals = []string{
		`func(){}`,
		`func(a){}`,
		`func(a,b){}`,
		`func(a=1){}`,
		`func(a,b="xxx"){}`,
		`func(a,b='bbb',c='ccc'){}`,
		`func(a,...b){}`,
	}
	var wants = []struct {
		Pars   []Parameter
		VarArg string
	}{
		{nil, ""},
		{[]Parameter{{"a", nil}}, ""},
		{[]Parameter{{"a", nil}, {"b", nil}}, ""},
		{[]Parameter{{"a", int64(1)}}, ""},
		{[]Parameter{{"a", nil}, {"b", "xxx"}}, ""},
		{[]Parameter{{"a", nil}, {"b", "bbb"}, {"c", "ccc"}}, ""},
		{[]Parameter{{"a", nil}}, "b"},
	}
	for i := range funcLiterals {
		l := newLexer(funcLiterals[i])
		exp := parseFuncLiteralExp(NewParser(l))
		if wants[i].VarArg != exp.VaArgs {
			t.Fatalf("parse func literal vararg failed:\n%s\n", funcLiterals[i])
		}
		if !reflect.DeepEqual(wants[i].Pars, exp.Parameters) {
			t.Fatalf("parse func literal parameters failed:\n%s\n", funcLiterals[i])
		}
	}
}
