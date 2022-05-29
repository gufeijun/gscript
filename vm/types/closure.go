package types

import "gscript/proto"

type GsValue struct {
	Value interface{}
}

type Closure struct {
	Info     *proto.BasicInfo
	UpValues []*GsValue
}
