package vm

import (
	"gscript/vm/types"
)

type stackFrame struct {
	prev        *stackFrame
	symbolTable *symbolTable
	wantRetCnt  int
	pc          uint32
	upValues    []*types.GsValue
	text        []byte
	tryInfos    []tryInfo
}

func newFuncFrame() *stackFrame {
	return &stackFrame{
		symbolTable: newSymbolTable(),
	}
}

type tryInfo struct {
	curVarCnt uint32
	catchAddr uint32
}

func (sf *stackFrame) pushTryInfo(addr uint32, varCnt uint32) {
	sf.tryInfos = append(sf.tryInfos, tryInfo{
		curVarCnt: varCnt,
		catchAddr: addr,
	})
}

func (sf *stackFrame) popTryInfo() (addr, varCnt uint32) {
	last := len(sf.tryInfos) - 1
	info := sf.tryInfos[last]
	sf.tryInfos = sf.tryInfos[:last]
	return info.catchAddr, info.curVarCnt
}
