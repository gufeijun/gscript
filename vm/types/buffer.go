package types

type Buffer struct {
	Data []byte
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{Data: data}
}

func NewBufferFromString(str string) *Buffer {
	return &Buffer{
		Data: []byte(str),
	}
}

func NewBufferN(cap int) *Buffer {
	return &Buffer{Data: make([]byte, cap)}
}
