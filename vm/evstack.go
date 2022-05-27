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

func (s *evalStack) top(n int) interface{} {
	return s.Buf[len(s.Buf)-n]
}

func (s *evalStack) popN(n int) {
	last := len(s.Buf)
	for i := 0; i < n; i++ {
		last--
		s.Buf[last] = nil
	}
	s.Buf = s.Buf[:last]
}

func (s *evalStack) Pop() {
	last := len(s.Buf) - 1
	s.Buf[last] = nil
	s.Buf = s.Buf[:last]
}

func (s *evalStack) pop() (val interface{}) {
	last := len(s.Buf) - 1
	val = s.Buf[last]
	s.Buf[last] = nil
	s.Buf = s.Buf[:last]
	return
}

func (s *evalStack) Push(v interface{}) {
	s.Buf = append(s.Buf, v)
}
