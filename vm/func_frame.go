package vm

type stackFrame struct {
	prev        *stackFrame
	symbolTable *symbolTable
	wantRetCnt  int
	returnAddr  uint32
}

func newFuncFrame() *stackFrame {
	return &stackFrame{
		symbolTable: newSymbolTable(),
	}
}
