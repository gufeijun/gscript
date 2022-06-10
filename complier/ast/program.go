package ast

// AST root
type Program struct {
	File       string // source file
	Imports    []Import
	BlockStmts []BlockStmt
	Export     Export
}

// import net,http as n,h
type Import struct {
	Line int
	Libs []Lib
}

type Lib struct {
	Stdlib bool
	Path   string
	Alias  string
}

// Block or Stmt
type BlockStmt interface{}

type Block struct {
	Blocks []BlockStmt
}

type Export struct {
	Exp Exp
}
