package vm

import "gscript/proto"

type protoFrame struct {
	topFrame       *stackFrame
	frame          *stackFrame
	stack          *evalStack
	constTable     []interface{}
	funcTable      []proto.FuncProto
	anonymousTable []proto.AnonymousFuncProto
	filepath       string
	prev           *protoFrame
}

func newProtoFrame(_proto proto.Proto) *protoFrame {
	topFrame := newFuncFrame()
	topFrame.text = _proto.Text
	frame := &protoFrame{
		topFrame:       topFrame,
		frame:          topFrame,
		stack:          newEvalStack(),
		constTable:     _proto.Consts,
		funcTable:      _proto.Funcs,
		anonymousTable: _proto.AnonymousFuncs,
		filepath:       _proto.FilePath,
	}
	return frame
}
