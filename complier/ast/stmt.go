package ast

type Stmt interface{}

type Var struct {
	Prefix string
	Attrs  []Exp
}

// let a,b,c = 1, "hello", add(1,2)
type VarDeclStmt struct {
	Line   int
	Lefts  []string
	Rights []Exp
}

// number, map["key"] = 2,"value"
// obj.total, i += 1,2
// i++ ==> i += 1
type VarAssignStmt struct {
	Line     int
	AssignOp int
	Lefts    []Var
	Rights   []Exp
}

// function call may like this:
// arr[1].FuncMap["Handlers"]("Sum")(1,2)
type NamedFuncCallStmt struct {
	Prefix    string
	CallTails []CallTail
}

type CallTail struct {
	Attrs []Exp
	Args  []Exp
}

type LabelStmt struct {
	Name string
}

// func sum(...arr)
// func add(a,b=1)
type FuncDefStmt struct {
	Name string
	FuncLiteral
}

type FuncLiteral struct {
	Line       int
	Parameters []Parameter
	VaArgs     string
	Block      Block
}

type Parameter struct {
	Name    string      // parameter name
	Default interface{} // string, number, true or false
}

type AnonymousFuncCallStmt struct {
	FuncLiteral
	CallTails []CallTail
}

type BreakStmt struct {
	Line int
}

type ContinueStmt struct {
	Line int
}

type GotoStmt struct {
	Line  int
	Label string
}

type FallthroughStmt struct {
	Line int
}

type WhileStmt struct {
	Condition Exp
	Block     Block
}

type ForStmt struct {
	// only one of these two is not nil
	AsgnStmt *VarAssignStmt
	DeclStmt *VarDeclStmt

	Condition Exp
	ForTail   *VarAssignStmt
	Block     Block
}

type LoopStmt struct {
	Key      string
	Val      string
	Iterator Exp
	Block    Block
}

// else ==> elif(true)
type IfStmt struct {
	Conditions []Exp
	Blocks     []Block
}

type ClassStmt struct {
	Name        string
	AttrName    []string
	AttrValue   []Exp
	Constructor *FuncLiteralExp
}

type EnumStmt struct {
	Names  []string
	Lines  []int
	Values []int64
}

type SwitchStmt struct {
	Value   Exp
	Cases   [][]Exp
	Blocks  [][]BlockStmt
	Default []BlockStmt
}

type ReturnStmt struct {
	Args []Exp
}

type TryCatchStmt struct {
	TryBlocks   []BlockStmt
	CatchValue  string
	CatchLine   int
	CatchBlocks []BlockStmt
}
