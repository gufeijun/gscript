package types

const MaxArrayCap = 8

type KV struct {
	Key interface{}
	Val interface{}
}

// when kv count is <= 8, use Array to stores kvs, otherwise use Map
type Object struct {
	Array []KV
	Map   map[interface{}]interface{}
}

func NewObjectN(cap int) *Object {
	obj := &Object{}
	if cap > MaxArrayCap {
		obj.Map = make(map[interface{}]interface{}, cap)
		return obj
	}
	obj.Array = make([]KV, 0, MaxArrayCap)
	return obj
}

func NewObject() *Object {
	return &Object{Array: make([]KV, 0, MaxArrayCap)}
}

func (obj *Object) Set(k, v interface{}) {
	if obj.Map != nil {
		obj.Map[k] = v
		return
	}
	if len(obj.Array) < MaxArrayCap {
		obj.Array = append(obj.Array, KV{k, v})
		return
	}
	obj.Map = make(map[interface{}]interface{}, 16)
	for _, kv := range obj.Array {
		obj.Map[kv.Key] = kv.Val
	}
	obj.Map[k] = v
}

func (obj *Object) Get(k interface{}) interface{} {
	if obj.Map != nil {
		return obj.Map[k]
	}
	for i := range obj.Array {
		if obj.Array[i].Key == k {
			return obj.Array[i].Val
		}
	}
	return nil
}

func (obj *Object) KVCount() int {
	if obj.Map != nil {
		return len(obj.Map)
	}
	return len(obj.Array)
}

func (obj *Object) Clone() *Object {
	if obj.Map != nil {
		m := make(map[interface{}]interface{}, len(obj.Map))
		for k, v := range obj.Map {
			m[k] = v
		}
		return &Object{Map: m}
	}
	arr := make([]KV, len(obj.Array), cap(obj.Array))
	copy(arr, obj.Array)
	return &Object{Array: arr}
}

func (obj *Object) Delete(key interface{}) {
	if obj.Map != nil {
		delete(obj.Map, key)
		return
	}
	for i := range obj.Array {
		if obj.Array[i].Key == key {
			obj.Array[i].Val = nil
		}
	}
}

func (obj *Object) ForEach(cb func(k, v interface{})) {
	if obj.Map != nil {
		for k, v := range obj.Map {
			cb(k, v)
		}
		return
	}
	for _, kv := range obj.Array {
		cb(kv.Key, kv.Val)
	}
}
