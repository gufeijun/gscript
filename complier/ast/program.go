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
	Libs []Lib
}

type Lib struct {
	Stdlib bool
	Path   string
	Alia   string
}

// Block or Stmt
type BlockStmt interface{}

type Block struct {
	Blocks []BlockStmt
}

type Export struct {
	Exp Exp
}
