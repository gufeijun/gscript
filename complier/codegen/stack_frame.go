package codegen

type StackFrame struct {
	prev        *StackFrame
	nt          *NameTable
	vt          *UpValueTable
	bs          *blockStack
	validLabels map[string]label
	returnAtEnd bool
	text        []byte

	nowParsingAnonymous int
	curTryLevel         int
}

func newStackFrame() *StackFrame {
	return &StackFrame{
		nt:          NewNameTable(),
		bs:          newBlockStack(),
		validLabels: map[string]label{},
		vt:          newUpValueTable(),
	}
}
