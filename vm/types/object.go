package types

type Object struct {
	Data map[interface{}]interface{}
}

func NewObjectN(cap int) *Object {
	return &Object{
		Data: make(map[interface{}]interface{}, cap),
	}
}

func NewObject() *Object {
	return &Object{Data: map[interface{}]interface{}{}}
}
