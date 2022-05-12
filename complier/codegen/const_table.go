package codegen

import "fmt"

type ConstTable struct {
	Constants []interface{}
	ConsMap   map[interface{}]uint32 // constant -> constants index
	enums     map[string]uint32      // enum -> constants index
}

func newConstTable() *ConstTable {
	return &ConstTable{
		ConsMap: make(map[interface{}]uint32),
		enums:   make(map[string]uint32),
	}
}

func (ct *ConstTable) saveEnum(enum string, num int64) {
	if _, ok := ct.enums[enum]; ok {
		panic(fmt.Sprintf("enum %s is already exists", enum))
	}
	ct.enums[enum] = uint32(len(ct.Constants))
	ct.Constants = append(ct.Constants, num)
}

func (ct *ConstTable) getEnum(enum string) (idx uint32, ok bool) {
	idx, ok = ct.enums[enum]
	return
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
