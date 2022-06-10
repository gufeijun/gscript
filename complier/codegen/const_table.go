package codegen

import (
	"fmt"
	"os"
)

type ConstTable struct {
	Constants []interface{}
	ConsMap   map[interface{}]uint32 // constant -> constants index
	enums     map[string]enum        // enum -> constants index
}

type enum struct {
	idx  uint32
	line uint32
}

func newConstTable() *ConstTable {
	return &ConstTable{
		ConsMap: make(map[interface{}]uint32),
		enums:   make(map[string]enum),
	}
}

func (ct *ConstTable) saveEnum(name string, line int, num int64) {
	if exists, ok := ct.enums[name]; ok {
		fmt.Printf("enum name %s already defines at line %d, but redeclares at line %d", name, exists.line, line)
		os.Exit(0)
	}
	ct.enums[name] = enum{
		idx:  uint32(len(ct.Constants)),
		line: uint32(line),
	}
	ct.Constants = append(ct.Constants, num)
}

func (ct *ConstTable) getEnum(name string) (idx uint32, ok bool) {
	enum, ok := ct.enums[name]
	return enum.idx, ok
}

func (ct *ConstTable) Get(key interface{}) uint32 {
	if idx, ok := ct.ConsMap[key]; ok {
		return idx
	}
	ct.Constants = append(ct.Constants, key)
	idx := uint32(len(ct.Constants) - 1)
	ct.ConsMap[key] = idx
	return idx
}
