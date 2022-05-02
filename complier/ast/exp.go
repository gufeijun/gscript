package ast

type Exp interface{}

type MapLiteralExp struct {
	Keys []interface{} // true, false, STRING, NUMBER
	Vals []Exp
}

type ArrLiteralExp struct {
	Vals []Exp
}

type StringLiteralExp struct {
	Value string
}

type NumberLiteralExp struct {
	Value interface{}
}

type FuncLiteralExp struct {
	FuncLiteral
}

type FalseExp struct{}

type TrueExp struct{}

type NilExp struct{}

type NewObjectExp struct {
	Name string
	Args []Exp
}

type NameExp struct {
	Name string
}

type UnOpExp struct {
	Op  int
	Exp Exp
}

// exp++
type IncExp struct {
	Exp Exp
}

// exp--
type DecExp struct {
	Exp Exp
}

// . ==> []
// a op b or a[b] or a.b
type BinOpExp struct {
	BinOp int
	Exp1  Exp
	Exp2  Exp
}

// a==b ? a:b
type TernaryOpExp struct {
	Exp1, Exp2, Exp3 Exp
}

type FuncCallExp struct {
	Func Exp
	Args []Exp
}
