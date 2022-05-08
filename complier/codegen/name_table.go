package codegen

type NameTable struct {
	nameTable map[string]uint32 // name -> index
	nameIdx   *uint32
	prev      *NameTable
}

func NewNameTable() *NameTable {
	return newNameTable(new(uint32))
}

func newNameTable(nameIdx *uint32) *NameTable {
	return &NameTable{
		nameTable: make(map[string]uint32),
		nameIdx:   nameIdx,
	}
}

func (nt *NameTable) Get(key string) (idx uint32) {
	idx, ok := nt.get(key)
	if ok {
		return idx
	}
	nt.nameTable[key] = *nt.nameIdx
	idx = *nt.nameIdx
	*nt.nameIdx++
	return
}

func (nt *NameTable) Set(name string) {
	if _, ok := nt.nameTable[name]; ok {
		panic("name exists already") // TODO
	}
	nt.nameTable[name] = *nt.nameIdx
	*nt.nameIdx++
}

func (nt *NameTable) get(name string) (uint32, bool) {
	for t := nt; t != nil; t = t.prev {
		if idx, ok := t.nameTable[name]; ok {
			return idx, true
		}
	}
	return 0, false
}
