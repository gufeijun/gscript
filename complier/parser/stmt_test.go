package parser

import (
	"gscript/complier/ast"
	. "gscript/complier/ast"
	"reflect"
	"testing"
)

func TestVarDeclStmt(t *testing.T) {
	srcs := []string{
		`let a="a"+"b";`,
		`let a;`,
		`const a,b = 1,2`,
		`let a,b = 1+1,1+2`,
		`let a,b,c = true,1+2,"good";`,
	}
	wants := []*VarDeclStmt{
		{false, []string{"a"}, []Exp{&BinOpExp{BINOP_ADD, &StringLiteralExp{"a"}, &StringLiteralExp{"b"}}}},
		{false, []string{"a"}, []Exp{&NilExp{}}},
		{true, []string{"a", "b"}, []Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}}},
		{false, []string{"a", "b"}, []Exp{
			&BinOpExp{BINOP_ADD, &NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(1)}},
			&BinOpExp{BINOP_ADD, &NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}},
		}},
		{false, []string{"a", "b", "c"}, []Exp{
			&TrueExp{},
			&BinOpExp{BINOP_ADD, &NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}},
			&StringLiteralExp{"good"},
		}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseVarDeclStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseVarDeclStmt failed: \n%s\n", src)
		}
	}
}

func TestVarAssignStmt(t *testing.T) {
	srcs := []string{
		"a,b += 1,1",
		"a /= 2",
		`c["ac"] %= 1 +1 `,
		`d.c.e,d[sum(1,"2")] *= 8,2`,
	}
	wants := []*ast.VarAssignStmt{
		{ASIGN_OP_ADDEQ, []Var{{"a", nil}, {"b", nil}}, []Exp{
			&NumberLiteralExp{int64(1)},
			&NumberLiteralExp{int64(1)},
		}},
		{ASIGN_OP_DIVEQ, []Var{{"a", nil}}, []Exp{&NumberLiteralExp{int64(2)}}},
		{ASIGN_OP_MODEQ, []Var{{"c", []Exp{&StringLiteralExp{"ac"}}}}, []Exp{
			&BinOpExp{BINOP_ADD, &NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(1)}},
		}},

		{ASIGN_OP_MULEQ, []Var{
			{"d", []Exp{&StringLiteralExp{"c"}, &StringLiteralExp{"e"}}},
			{"d", []Exp{&FuncCallExp{
				Func: &NameExp{"sum"},
				Args: []Exp{&NumberLiteralExp{int64(1)}, &StringLiteralExp{"2"}},
			}}},
		}, []Exp{&NumberLiteralExp{int64(8)}, &NumberLiteralExp{int64(2)}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseVarOpOrLabel()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseVarAssignStmt failed: \n%s\n", src)
		}
	}
}

func TestVarIncOrDecStmt(t *testing.T) {
	srcs := []string{
		"arr[i]++",
		`obj["key"].a--`,
	}
	wants := []*VarAssignStmt{
		{ASIGN_OP_ADDEQ, []Var{{"arr", []Exp{&NameExp{"i"}}}}, []Exp{&NumberLiteralExp{int64(1)}}},
		{ASIGN_OP_SUBEQ, []Var{{"obj", []Exp{
			&StringLiteralExp{"key"},
			&StringLiteralExp{"a"},
		}}}, []Exp{&NumberLiteralExp{int64(1)}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseVarOpOrLabel()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseVarIncOrDecStmt failed: \n%s\n", src)
		}
	}
}

func TestLabelStmt(t *testing.T) {
	src := "loop:"
	want := &LabelStmt{"loop"}

	l := newLexer(src)
	stmt := NewParser(l).parseVarOpOrLabel()
	if !reflect.DeepEqual(stmt, want) {
		t.Fatalf("parseLabelStmt failed: \n%s\n", src)
	}
}

func TestFuncCallStmt(t *testing.T) {
	srcs := []string{
		`funcs[0](a,b)`,
		`f().a["xxx"](1)()().a["xxx"](1)`,
	}
	wants := []*NamedFuncCallStmt{
		&NamedFuncCallStmt{
			Var:       Var{"funcs", []Exp{&NumberLiteralExp{int64(0)}}},
			Args:      []Exp{&NameExp{"a"}, &NameExp{"b"}},
			CallTails: nil,
		},
		&NamedFuncCallStmt{
			Var:  Var{"f", nil},
			Args: nil,
			CallTails: []CallTail{
				{[]Exp{&StringLiteralExp{"a"}, &StringLiteralExp{"xxx"}}, []Exp{&NumberLiteralExp{int64(1)}}},
				{nil, nil},
				{nil, nil},
				{[]Exp{&StringLiteralExp{"a"}, &StringLiteralExp{"xxx"}}, []Exp{&NumberLiteralExp{int64(1)}}},
			},
		},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseVarOpOrLabel()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseFuncCallStmt failed: \n%s\n", src)
		}
	}
}

func TestFuncDefStmt(t *testing.T) {
	srcs := []string{
		`
func A(a,b=1){
	let a=b
	print(a)
}
`,
		`
func A(a,...vararg){}
`,
		`
func A(a=1,b="good"){
	return a,b
}
`,
	}
	wants := []*FuncDefStmt{
		{"A", FuncLiteral{
			Parameters: []Parameter{{"a", nil}, {"b", &NumberLiteralExp{int64(1)}}},
			VaArgs:     "",
			Block: Block{Blocks: []BlockStmt{
				NewParser(newLexer("let a=b;")).parseVarDeclStmt(),
				NewParser(newLexer("print(a)")).parseVarOpOrLabel(),
			}},
		}},
		{"A", FuncLiteral{
			Parameters: []Parameter{{"a", nil}},
			VaArgs:     "vararg",
		}},
		{"A", FuncLiteral{
			Parameters: []Parameter{{"a", &NumberLiteralExp{int64(1)}}, {"b", &StringLiteralExp{"good"}}},
			VaArgs:     "",
			Block: Block{Blocks: []BlockStmt{
				&ReturnStmt{
					Args: []Exp{&NameExp{"a"}, &NameExp{"b"}},
				},
			}},
		}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseFunc()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseFuncDefStmt failed: \n%s\n", src)
		}
	}
}

func TestAnonymousFuncCall(t *testing.T) {
	srcs := []string{
		`
func(a,b){return a+b;}(1,2)
 `,
		`
func(num){
	return {
		show:func(){
			print(num)
		}
	}
}(1).show()
`,
		`
func(num){
	return {show:func(){print(num)}}
}(1)[echo("show")]()
`,
	}
	wants := []*AnonymousFuncCallStmt{
		&AnonymousFuncCallStmt{
			FuncLiteral: FuncLiteral{
				Parameters: []Parameter{{"a", nil}, {"b", nil}},
				Block: Block{Blocks: []BlockStmt{&ReturnStmt{
					Args: []Exp{&BinOpExp{BINOP_ADD, &NameExp{"a"}, &NameExp{"b"}}}}}}},
			CallArgs:  []Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}},
			CallTails: nil,
		},
		&AnonymousFuncCallStmt{
			FuncLiteral: FuncLiteral{
				Parameters: []Parameter{{"num", nil}},
				Block: Block{[]BlockStmt{&ReturnStmt{Args: []Exp{&MapLiteralExp{
					Keys: []interface{}{"show"},
					Vals: []Exp{parseFuncLiteralExp(NewParser(newLexer("func(){print(num)}")))},
				}}}}}},
			CallArgs:  []Exp{&NumberLiteralExp{int64(1)}},
			CallTails: []CallTail{{[]Exp{&StringLiteralExp{"show"}}, nil}},
		},
		&AnonymousFuncCallStmt{
			FuncLiteral: FuncLiteral{
				Parameters: []Parameter{{"num", nil}},
				Block: Block{[]BlockStmt{&ReturnStmt{[]Exp{&MapLiteralExp{
					Keys: []interface{}{"show"},
					Vals: []Exp{parseFuncLiteralExp(NewParser(newLexer("func(){print(num)}")))},
				}}}}}},
			CallArgs: []Exp{&NumberLiteralExp{int64(1)}},
			CallTails: []CallTail{
				{[]Exp{parseFuncCallOrAttrExp(NewParser(newLexer(`echo("show")`)))}, nil}},
		},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseFunc()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseAnonymousFuncCallStmt failed: \n%s\n", src)
		}
	}
}

func TestIncOrDecVarStmt(t *testing.T) {
	srcs := []string{
		"++arr[i]",
		`--obj["key"].a`,
	}
	wants := []*VarAssignStmt{
		{ASIGN_OP_ADDEQ, []Var{{"arr", []Exp{&NameExp{"i"}}}}, []Exp{&NumberLiteralExp{int64(1)}}},
		{ASIGN_OP_SUBEQ, []Var{{"obj", []Exp{
			&StringLiteralExp{"key"},
			&StringLiteralExp{"a"},
		}}}, []Exp{&NumberLiteralExp{int64(1)}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseIncOrDecVar()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseIncOrDecVarStmt failed: \n%s\n", src)
		}
	}
}

func TestEnumStmt(t *testing.T) {
	srcs := []string{
		`enum {}`,
		`
enum {
	a=1,
	b,c=8,d,
	e,f=19,
	g
}
`,
	}
	wants := []*EnumStmt{
		{nil, nil},
		{
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			[]int64{1, 2, 8, 9, 10, 19, 20},
		},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseEnumStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseEnumStmt failed: \n%s\n", src)
		}
	}
}

func TestSwitchStmt(t *testing.T) {
	srcs := []string{
		`switch(a){}`,
		`
switch(a)
{
case 1:
	i++
	let b,c = 1,2
case 2,3:
	i--
	call(a,b)
}`,
		`
switch(a){
case "hello":
	++i
default:
	--i
}
`,
	}
	wants := []*SwitchStmt{
		{
			Value: &NameExp{"a"},
		},
		{
			Value: &NameExp{"a"},
			Cases: [][]Exp{{&NumberLiteralExp{int64(1)}}, {&NumberLiteralExp{int64(2)}, &NumberLiteralExp{int64(3)}}},
			Blocks: [][]BlockStmt{
				{NewParser(newLexer("i++")).parseVarOpOrLabel(), NewParser(newLexer("let b,c = 1,2")).parseVarDeclStmt()},
				{NewParser(newLexer("i--")).parseVarOpOrLabel(), NewParser(newLexer("call(a,b)")).parseVarOpOrLabel()},
			},
		},
		{
			Value:   &NameExp{"a"},
			Cases:   [][]Exp{{&StringLiteralExp{"hello"}}},
			Blocks:  [][]BlockStmt{{NewParser(newLexer("++i")).parseIncOrDecVar()}},
			Default: []BlockStmt{NewParser(newLexer("--i")).parseIncOrDecVar()},
		},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseSwitchStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseSwitchStmt failed: \n%s\n", src)
		}
	}
}

func TestLoopStmt(t *testing.T) {
	srcs := []string{
		`
loop(let v:m["arr"])
	print(v)
`,
		`
loop(let k,v:arr)
	print(k,v)
`,
		`
loop(let k,v:arr){
	print(k,v)
}
`,
	}
	wants := []*LoopStmt{
		{"", "v", parseExp(NewParser(newLexer(`m["arr"]`))), Block{
			[]BlockStmt{NewParser(newLexer("print(v)")).parseVarOpOrLabel()},
		}},
		{"k", "v", parseExp(NewParser(newLexer("arr"))), Block{
			[]BlockStmt{NewParser(newLexer("print(k,v)")).parseVarOpOrLabel()},
		}},
		{"k", "v", parseExp(NewParser(newLexer("arr"))), Block{
			[]BlockStmt{NewParser(newLexer("print(k,v)")).parseVarOpOrLabel()},
		}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseLoopStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseLoopStmt failed: \n%s\n", src)
		}
	}
}

func TestWhileStmt(t *testing.T) {
	srcs := []string{
		`
while(1){
	let conn = accept(listener)
	thread(handleConn,conn)
}
`,
		`
while(i--)
	print(i)
`,
	}

	wants := []*WhileStmt{
		{&NumberLiteralExp{int64(1)}, Block{
			[]BlockStmt{
				NewParser(newLexer("let conn = accept(listener)")).parseVarDeclStmt(),
				NewParser(newLexer("thread(handleConn,conn)")).parseVarOpOrLabel(),
			},
		}},
		{parseExp(NewParser(newLexer("i--"))), Block{[]BlockStmt{NewParser(newLexer("print(i)")).parseVarOpOrLabel()}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseWhileStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseWhileStmt failed: \n%s\n", src)
		}
	}
}

func TestForStmt(t *testing.T) {
	srcs := []string{
		`
for(let low,high=0,len(arr)-1;low<high;low,high = low+1,high-1){
	arr[low],arr[high]=arr[high],arr[low]
}
`,
		`
for(i=0;i<10;i++)
	print(i)
`,
	}
	wants := []*ForStmt{
		{nil, NewParser(newLexer("let low,high=0,len(arr)-1")).parseVarDeclStmt(),
			parseExp(NewParser(newLexer("low<high"))),
			NewParser(newLexer("low,high = low+1,high-1")).parseVarOpOrLabel().(*VarAssignStmt),
			Block{[]BlockStmt{NewParser(newLexer("arr[low],arr[high]=arr[high],arr[low]")).parseVarOpOrLabel()}},
		},
		{NewParser(newLexer("i=0")).parseVarOpOrLabel().(*VarAssignStmt), nil,
			parseExp(NewParser(newLexer("i<10"))), NewParser(newLexer("i++")).parseVarOpOrLabel().(*VarAssignStmt),
			Block{[]BlockStmt{NewParser(newLexer("print(i)")).parseVarOpOrLabel()}},
		},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseForStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseForStmt failed: \n%s\n", src)
		}
	}
}

func TestIfStmt(t *testing.T) {
	srcs := []string{
		`
if(a){
	print(true)
}`,
		`
if(a)
	print(true)
`,
		`
if(a) print(1)
elif(b) print(2);
elif(c) print(3)
`,
		`
if(a) print(1)
else print(2)
`,
	}
	wants := []*IfStmt{
		{[]Exp{&NameExp{"a"}}, []Block{
			Block{[]BlockStmt{NewParser(newLexer("print(true)")).parseVarOpOrLabel()}},
		}},
		{[]Exp{&NameExp{"a"}}, []Block{
			Block{[]BlockStmt{NewParser(newLexer("print(true)")).parseVarOpOrLabel()}},
		}},
		{[]Exp{&NameExp{"a"}, &NameExp{"b"}, &NameExp{"c"}}, []Block{
			Block{[]BlockStmt{NewParser(newLexer("print(1)")).parseVarOpOrLabel()}},
			Block{[]BlockStmt{NewParser(newLexer("print(2)")).parseVarOpOrLabel()}},
			Block{[]BlockStmt{NewParser(newLexer("print(3)")).parseVarOpOrLabel()}},
		}},
		{[]Exp{&NameExp{"a"}, &TrueExp{}}, []Block{
			Block{[]BlockStmt{NewParser(newLexer("print(1)")).parseVarOpOrLabel()}},
			Block{[]BlockStmt{NewParser(newLexer("print(2)")).parseVarOpOrLabel()}},
		}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseIfStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseIfStmt failed: \n%s\n", src)
		}
	}
}

func TestParseReturnStmt(t *testing.T) {
	srcs := []string{
		`return;`,
		`
return {a:1,b:2};
`,
		`
		return a,b
		`,
	}
	wants := []*ReturnStmt{
		{nil},
		{Args: []Exp{&MapLiteralExp{
			Keys: []interface{}{"a", "b"},
			Vals: []Exp{&NumberLiteralExp{int64(1)}, &NumberLiteralExp{int64(2)}},
		}}},
		{Args: []Exp{&NameExp{"a"}, &NameExp{"b"}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseReturnStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseReturnStmt failed: \n%s\n", src)
		}
	}
}

func TestClassStmt(t *testing.T) {
	srcs := []string{
		// `class A{}`,
		`
class A{
	__self(){
		this.age = 10;
	}
	name="jack"
	age
	show = func(){
		#...
	}
}`,
	}
	wants := []*ClassStmt{
		// {"A", nil, nil},
		{"A", []string{"__self", "name", "show"}, []Exp{&FuncLiteralExp{
			FuncLiteral{nil, "", Block{[]BlockStmt{NewParser(newLexer("this.age = 10")).parseVarOpOrLabel()}}}},
			&StringLiteralExp{"jack"}, &FuncLiteralExp{FuncLiteral{nil, "", Block{}}}}},
	}
	for i, src := range srcs {
		l := newLexer(src)
		stmt := NewParser(l).parseClassStmt()
		if !reflect.DeepEqual(stmt, wants[i]) {
			t.Fatalf("parseReturnStmt failed: \n%s\n", src)
		}
	}
}
