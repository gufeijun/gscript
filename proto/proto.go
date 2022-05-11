package proto

type Proto struct {
	Text           []Instruction
	Consts         []interface{}
	Funcs          []FuncProto
	AnonymousFuncs []AnonymousFuncProto
}
