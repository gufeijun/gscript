## Builtin Functions

### print

```python
print(1,"2",[1,2,3],{foo:"bar"})	# 1 2 Array[1, 2, 3] Object{foo: bar}
```

+ parameter
  + count: `>=0`
  + type: any

### len

```python
let l = len([1,2,3])      # l = 3
l = len({foo:"bar"})      # l = 1
l = len("hello")          # l = 5
```

+ parameter
  + count: `1`
  + type: `String`, `Array`, `Object` or `Buffer`

+ return
  + count: `1`
  + type: `Number`

### append

```python
let arr = []
append(arr,1,2)		# arr = [1,2]
```

+ parameter
  + count: `>=2`
  + type: `arg0(Array)`,`arg1(any)`, `arg2(any)`, ... ,`argn(any)`

### sub

sub(src, start[,end])

```python
let str = "1234"
let s1 = sub(str,1)         # s1 = "234"
let s2 = sub(str,1,2)       # s2 = "2"

let arr = [1,2,3,4]
ler a1 = sub(arr,1)         # a1 = [2,3,4]
let a2 = sub(arr,1,2)       # a2 = [2]
```

+ description: get sub array or sub string. 
+ parameter
  + count: `2` or `3`
  + type: `arg0(String or Array)`, `start(Integer)`, `end(Integer)`

+ return
  + count: 1
  + type: `String` or `Array`

### type

```python
print(type(""))         # String
print(type(func(){}))   # closure
print(type(type))       # Builtin
print(type({}))         # Object
print(type([]))         # Array
print(type(false))      # Boolean
print(type(nil))        # Nil
```

### delete

```python
let obj = {foo:"bar"}
delete(obj,"foo")
print(obj.foo)				# <nil>
```

+ parameter
  + count: `2`
  + type: `arg0(Object)`, `arg1(any)`

### clone

```python
let arr1 = [1,2,3]
let arr2, arr3 = arr1, clone(arr1)
arr1[0] = -1
print(arr1)				# [-1,2,3]
print(arr2)				# [-1,2,3]
print(arr3)				# [1,2,3]
```

+ parameter
  + count: `1`
  + type: `Array`, `Object`, `Buffer`
+ return
  + count: `1`
  + type: `Array`, `Object`, `Buffer`

### throw

```python
try{
    let arr = [1,2,3]
    throw(arr)
}
catch(e) {
    print(e)		# [1,2,3]
}
```

### others

there are many other builtin functions, but always dangerous to use. So we wrap these apis to standard libraries, please use these libraries instead. 
