package codegen

type UpValue struct {
	level   uint32
	nameIdx uint32
	name    string
}

type UpValueTable struct {
	upValues []UpValue
}

func newUpValueTable() *UpValueTable {
	return &UpValueTable{}
}

func (vt *UpValueTable) get(name string) (upValueIdx uint32, level uint32, ok bool) {
	for i, v := range vt.upValues {
		if v.name == name {
			return uint32(i), v.level, true
		}
	}
	return 0, 0, false
}

func (vt *UpValueTable) set(name string, level uint32, nameIdx uint32) (upValueIdx uint32) {
	vt.upValues = append(vt.upValues, UpValue{level, nameIdx, name})
	return uint32(len(vt.upValues) - 1)
}
