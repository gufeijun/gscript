package ast

type Stmt interface{}

type Var struct {
	Prefix string
	Attrs  []Exp
}

// const|let a,b,c = 1, "hello", add(1,2)
type VarDeclStmt struct {
	Const  bool
	DeepEq bool
	Lefts  []string
	Rights []Exp
}

// number, map["key"] = 2,"value"
// obj.total, i += 1,2
// i++ ==> i += 1
type VarAssignStmt struct {
	AssignOp int
	Lefts    []Var
	Rights   []Exp
}

// function call may like this:
// arr[1].FuncMap["Handlers"]("Sum")(1,2)
type NamedFuncCallStmt struct {
	Var       Var
	Args      []Exp
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

// TODO return
type FuncLiteral struct {
	Parameters []Parameter
	VarArg     string
	Block      Block
}

type Parameter struct {
	Name    string // parameter name
	Default Exp    // default value
}

type AnonymousFuncCallStmt struct {
	FuncLiteral
	CallArgs  []Exp
	CallTails []CallTail
}

type BreakStmt struct{}

type ContinueStmt struct{}

type GotoStmt struct {
	Label string
}

type FallthroughStmt struct{}

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
	Name      string
	AttrName  []string
	AttrValue []Exp
}

type EnumStmt struct {
	Names  []string
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
