package vm

type funcFrame struct {
	prev        *funcFrame
	symbolTable *symbolTable
	wantRetCnt  int
	returnAddr  uint32
}

func newFuncFrame() *funcFrame {
	return &funcFrame{
		symbolTable: newSymbolTable(),
	}
}
