package parser

import (
	. "gscript/complier/ast"
	. "gscript/complier/lexer"
	"reflect"
	"testing"
)

func newLexer(src string) *Lexer {
	return NewLexer("", []byte(src))
}

func TestParseTerm12(t *testing.T) {
	srcs := []string{
		`a||c?b+d:c`,
		`a?(b?c:d):e`,
	}
	var wants = []*TernaryOpExp{
		{&BinOpExp{BINOP_LOR, &NameExp{"a"}, &NameExp{"c"}}, &BinOpExp{BINOP_ADD, &NameExp{"b"}, &NameExp{"d"}}, &NameExp{"c"}},
		{&NameExp{"a"}, &TernaryOpExp{&NameExp{"b"}, &NameExp{"c"}, &NameExp{"d"}}, &NameExp{"e"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm12(l)
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
		{BINOP_LOR, &BinOpExp{BINOP_LOR, &NameExp{"a"}, &NameExp{"b"}}, &NameExp{"c"}},
		{BINOP_LOR, &NameExp{"a"}, &BinOpExp{BINOP_LAND, &NameExp{"b"}, &NameExp{"c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm11(l)
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
		{BINOP_LAND, &NameExp{"a"}, &NameExp{"b"}},
		{BINOP_LAND, &NameExp{"a"}, &BinOpExp{BINOP_OR, &NameExp{"b"}, &NameExp{"c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm10(l)
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
		{BINOP_OR, &NameExp{"a"}, &NameExp{"b"}},
		{BINOP_OR, &NameExp{"a"}, &BinOpExp{BINOP_XOR, &NameExp{"b"}, &NameExp{"c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm9(l)
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
		{BINOP_XOR, &NameExp{"a"}, &NameExp{"b"}},
		{BINOP_XOR, &NameExp{"a"}, &BinOpExp{BINOP_AND, &NameExp{"b"}, &NameExp{"c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm8(l)
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
		{BINOP_AND, &NameExp{"a"}, &NameExp{"b"}},
		{BINOP_AND, &NameExp{"a"}, &BinOpExp{BINOP_EQ, &NameExp{"b"}, &NameExp{"c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm7(l)
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
		{BINOP_EQ, &NameExp{"a"}, &NameExp{"b"}},
		{BINOP_EQ, &BinOpExp{BINOP_NE, &NameExp{"a"}, &NameExp{"b"}}, &NameExp{"c"}},
		{BINOP_EQ, &NameExp{"a"}, &BinOpExp{BINOP_GE, &NameExp{"b"}, &NameExp{"c"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm6(l)
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
		{BINOP_GT, &NameExp{"a"}, &BinOpExp{BINOP_SHL, &NameExp{"b"}, &NameExp{"c"}}},
		{BINOP_LE, &NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(3)}},
		{BINOP_LT, &NumberLiteralExp{int64(2)}, &NumberLiteralExp{int64(4)}},
		{BINOP_GT, &BinOpExp{BINOP_GE, &NameExp{"a"}, &NumberLiteralExp{int64(8)}}, &NumberLiteralExp{int64(9)}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm5(l)
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
		{BINOP_SHL, &NameExp{"a"}, &BinOpExp{BINOP_ADD, &NameExp{"b"}, &NameExp{"c"}}},
		{BINOP_SHL, &BinOpExp{BINOP_SHR, &NameExp{"a"}, &NameExp{"b"}}, &NameExp{"d"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm4(l)
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
		{BINOP_ADD, &NameExp{"a"}, &NameExp{"b"}},
		{BINOP_SUB, &BinOpExp{BINOP_ATTR, &NameExp{"a"}, &StringLiteralExp{"b"}}, &NameExp{"c"}},
		{BINOP_SUB, &NameExp{"a"}, &BinOpExp{BINOP_MUL, &NameExp{"b"}, &NameExp{"c"}}},
		{BINOP_MUL, &BinOpExp{BINOP_SUB, &NameExp{"a"}, &NameExp{"b"}}, &NameExp{"c"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm3(l)
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
		{BINOP_DIV, &UnOpExp{UNOP_NEG, &NameExp{"a"}}, &NameExp{"b"}},
		{BINOP_MUL, &UnOpExp{UNOP_NOT, &BinOpExp{BINOP_ATTR, &NameExp{"arr"}, &NumberLiteralExp{int64(0)}}}, &NumberLiteralExp{int64(100)}},
		{BINOP_IDIV, &UnOpExp{UNOP_INC_, &BinOpExp{BINOP_ATTR, &NameExp{"arr"}, &NumberLiteralExp{int64(0)}}}, &NumberLiteralExp{int64(20)}},
		{BINOP_MOD, &UnOpExp{UNOP_DEC, &BinOpExp{BINOP_ATTR, &NameExp{"arr"}, &NumberLiteralExp{int64(0)}}}, &NameExp{"p"}},
		{BINOP_MUL, &BinOpExp{BINOP_DIV, &NameExp{"a"}, &NameExp{"b"}}, &NameExp{"c"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm2(l)
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
		{UNOP_NOT, &UnOpExp{UNOP_INC, &NameExp{"a"}}},
		{UNOP_NEG, &UnOpExp{UNOP_DEC_, &NameExp{"a"}}},
		{UNOP_NEG, &BinOpExp{BINOP_ATTR, &NameExp{"arr"}, &NumberLiteralExp{int64(1)}}},
		{UNOP_LNOT, &UnOpExp{UNOP_LNOT, &UnOpExp{UNOP_NOT, &NumberLiteralExp{int64(1)}}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm1(l)
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
		{UNOP_INC, &NameExp{"a"}},
		{UNOP_DEC, &BinOpExp{BINOP_ATTR, &NameExp{"a"}, &NumberLiteralExp{int64(0)}}},
		{UNOP_INC_, &NameExp{"a"}},
		{UNOP_DEC_, &BinOpExp{BINOP_ATTR, &NameExp{"a"}, &StringLiteralExp{"b"}}},
	}

	for i, src := range srcs {
		l := newLexer(src)
		exp := parseTerm0(l)
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
					Exp1:  &NameExp{"m"},
					Exp2: &FuncCallExp{
						Func: &NameExp{"sum"},
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
	exp := parseFuncCallOrAttrExp(l)
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
		{&NameExp{"sum"}, []Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
		{&FuncCallExp{&NameExp{"sum"}, nil}, []Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseFuncCallOrAttrExp(l)
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
		{BINOP_ATTR, &NameExp{"map"}, &StringLiteralExp{"a"}},
		{BINOP_ATTR, &BinOpExp{
			BINOP_ATTR,
			&NameExp{"map"},
			&StringLiteralExp{"b"},
		}, &NumberLiteralExp{int64(1)}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseFuncCallOrAttrExp(l)
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
		{BINOP_ATTR, &NameExp{"person"}, &StringLiteralExp{"name"}},
		{BINOP_ATTR, &BinOpExp{
			BINOP_ATTR,
			&NameExp{"person"},
			&StringLiteralExp{"father"},
		}, &StringLiteralExp{"name"}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseFuncCallOrAttrExp(l)
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
		{[]Exp{&NameExp{"a"}, &NameExp{"b"}}},
		{[]Exp{&NumberLiteralExp{int64(1)}, &ArrLiteralExp{
			[]Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}},
		}}},
		{[]Exp{&StringLiteralExp{"good"}}},
		{[]Exp{&MapLiteralExp{
			Keys: []interface{}{"a"},
			Vals: []Exp{&NameExp{"a"}},
		}}},
		{[]Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
		{[]Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
	}
	for i, literal := range arrLiterals {
		l := newLexer(literal)
		exp := parseArrLiteralExp(l)
		if !reflect.DeepEqual(exp, wants[i]) {
			t.Fatalf("parse array literal failed:\n%s\n", literal)
		}
	}
}

func TestParseMapLiteralExp(t *testing.T) {
	var mapLiterals = []string{
		`{}`,
		`{a}`,
		`{a,1}`,
		`{interger:1,string:"b",float:1.89}`,
		`{true:true,false:false}`,
		`{1:"1",3.14:"3.14",}`,
		`{"child":nil}`,
	}
	var wants = []*MapLiteralExp{
		{nil, nil},
		{[]interface{}{"a"}, []Exp{&NilExp{}}},
		{[]interface{}{"a", int64(1)}, []Exp{&NilExp{}, &NilExp{}}},
		{[]interface{}{"interger", "string", "float"},
			[]Exp{&NumberLiteralExp{int64(1)}, &StringLiteralExp{"b"}, &NumberLiteralExp{float64(1.89)}}},
		{[]interface{}{true, false}, []Exp{&TrueExp{}, &FalseExp{}}},
		{[]interface{}{int64(1), float64(3.14)}, []Exp{&StringLiteralExp{"1"}, &StringLiteralExp{"3.14"}}},
		{[]interface{}{"child"}, []Exp{&NilExp{}}},
	}
	for i, literal := range mapLiterals {
		l := newLexer(literal)
		exp := parseMapLiteralExp(l)
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
		{"student", []Exp{&NameExp{"name"}}},
		{"student", []Exp{&NameExp{"name"}, &NameExp{"age"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		exp := parseNewObjectExp(l)
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
		{[]Parameter{{"a", &NumberLiteralExp{int64(1)}}}, ""},
		{[]Parameter{{"a", nil}, {"b", &StringLiteralExp{"xxx"}}}, ""},
		{[]Parameter{{"a", nil}, {"b", &StringLiteralExp{"bbb"}}, {"c", &StringLiteralExp{"ccc"}}}, ""},
		{[]Parameter{{"a", nil}}, "b"},
	}
	for i := range funcLiterals {
		l := newLexer(funcLiterals[i])
		exp := parseFuncLiteralExp(l)
		if wants[i].VarArg != exp.VarArg {
			t.Fatalf("parse func literal vararg failed:\n%s\n", funcLiterals[i])
		}
		if !reflect.DeepEqual(wants[i].Pars, exp.Parameters) {
			t.Fatalf("parse func literal parameters failed:\n%s\n", funcLiterals[i])
		}
	}
}