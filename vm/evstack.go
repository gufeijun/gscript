package vm

type evalStack struct {
	Buf []interface{}
}

func newEvalStack() *evalStack {
	return &evalStack{}
}

func (s *evalStack) Replace(v interface{}) {
	s.Buf[len(s.Buf)-1] = v
}

func (s *evalStack) Top() interface{} {
	return s.Buf[len(s.Buf)-1]
}

func (s *evalStack) Pop() {
	last := len(s.Buf) - 1
	s.Buf[last] = nil
	s.Buf = s.Buf[:last]
}

func (s *evalStack) Push(v interface{}) {
	s.Buf = append(s.Buf, v)
}
