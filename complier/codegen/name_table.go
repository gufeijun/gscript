package codegen

import (
	"fmt"
	"os"
)

type NameTable struct {
	nameTable map[string]variable
	nameIdx   *uint32
	prev      *NameTable
}

type variable struct {
	idx  uint32
	line uint32
}

func NewNameTable() *NameTable {
	return newNameTable(new(uint32))
}

func newNameTable(nameIdx *uint32) *NameTable {
	return &NameTable{
		nameTable: make(map[string]variable),
		nameIdx:   nameIdx,
	}
}

func (nt *NameTable) Set(name string, line uint32) {
	if v, ok := nt.nameTable[name]; ok {
		fmt.Printf("variable '%s' already declared at line %d, but redecalres at line %d", name, v.line, line)
		os.Exit(0)
	}
	nt.nameTable[name] = variable{
		idx:  *nt.nameIdx,
		line: line,
	}
	*nt.nameIdx++
}

func (nt *NameTable) get(name string) (uint32, bool) {
	for t := nt; t != nil; t = t.prev {
		if v, ok := t.nameTable[name]; ok {
			return v.idx, true
		}
	}
	return 0, false
}
