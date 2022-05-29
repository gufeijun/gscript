package types

type Array struct {
	Data []interface{}
}

func NewArray(data []interface{}) *Array {
	return &Array{Data: data}
}
