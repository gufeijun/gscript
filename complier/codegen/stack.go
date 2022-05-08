package codegen

type forBlock struct {
	nameCnt   uint32
	prev      interface{}
	breaks    []int
	continues []int
}

type switchBlock struct {
	nameCnt      uint32
	prev         interface{}
	breaks       []int
	_fallthrough *int
}

type blockStack struct {
	cur interface{}
}

func newBlockStack() *blockStack {
	return &blockStack{}
}

func (bs *blockStack) pop() {
	if b, ok := bs.cur.(*forBlock); ok {
		bs.cur = b.prev
		return
	}
	if b, ok := bs.cur.(*switchBlock); ok {
		bs.cur = b.prev
		return
	}
}

func (bs *blockStack) pushSwitch(nameCnt uint32) {
	b := &switchBlock{nameCnt: nameCnt, prev: bs.cur}
	bs.cur = b
}

func (bs *blockStack) pushFor(nameCnt uint32) {
	b := &forBlock{nameCnt: nameCnt, prev: bs.cur}
	bs.cur = b
}

func (bs *blockStack) latestFor() *forBlock {
	cur := bs.cur
	for cur != nil {
		if fb, ok := cur.(*forBlock); ok {
			return fb
		}
		cur = cur.(*switchBlock).prev
	}
	panic("unmatched continue")
}

func (bs *blockStack) latestSwitch() *switchBlock {
	cur := bs.cur
	for cur != nil {
		if fb, ok := cur.(*switchBlock); ok {
			return fb
		}
		cur = cur.(*forBlock).prev
	}
	panic("unmatched fallthrough")
}

func (bs *blockStack) top() interface{} {
	if bs.cur == nil {
		panic("unmatched continue or break")
	}
	return bs.cur
}
