## Gscript 文法

### 类型

下列是`gscript`最常用的类型：

```python
100					# Number
100.21					# Number
"gscript"				# String
false					# Boolean
[1,"gs",{}]				# Array
{foo: 1, bar:"bar"}			# Objet
func(){}				# Closure
nil					# nil
```

**Array**

```python
let arr = [1,2,[1,2,3]]
arr[0]					# 1
arr[2][2]				# 3
arr[3]					# out of range, will panic
```

**Object**

```python
let obj = {
    foo:"foo",
    bar:{
        1:"1",
        arr:[1,2,nil],
    },
}

obj["foo"]				# "foo"
obj.foo					# "foo"
obj.bar[1]				# "1"
obj.bar.arr				# [1,2,nil]
```

### 变量以及作用域

使用关键字`let`去声明一个局部变量。使用未定义的变量会让编译失败。

```python
let a = 1

# 进入新的作用域
{
    print(a)				# a==1
    let a = 2
    let b = 3
    print(a)				# a==2
    print(b)				# b==3
}

print(a)				# a==1
print(b)				# 非法，编译失败
```

在变量的作用域外使用一个变量会编译失败。

和很多脚本语言一样，一个变量被赋予了一种类型值后，能够重新被赋予新的类型。

```python
let a = {}
a = []					# ok
```

可以只用一个`let`声明多个变量：

```python
let a, b = 1, "hello"	# a=1,b="hello"
let a, b = 1			# a=1,b=nil
let a = 1,"hello"		# a=1
```

当然也可以同时为多个变量赋值：

```python
let a,b
a, b = 1, 2 
```

### 操作符

**单目操作符**

| 操作符 |      用法      |        适用类型        |
| :----: | :------------: | :--------------------: |
|  `~`   |    按位取反    |        Integer         |
|  `-`   |    `0 - x`     | Number(Integer, Float) |
|  `!`   |      `!x`      |          all           |
|  `--`  | `--x` or `x--` |         Number         |
|  `++`  | `++x` or `x++` |         Number         |

注意`--`以及`++`的使用规范。如果这两个被用作statement：

```python
# x++ is the same as ++x
let x = 0
x++			# x==1
++x			# x==2

# x-- is the same as --x
x--			# x==1
--x			# x==0
```

如果被用作expression：

```python
let x = 0
print(x++)	# output: 0
print(++x)	# output: 2
print(--x)	# output: 1
print(x--)	# output: 1
# now x==0
```

 `--` 和 `++`的行为和 `JavaScript` and `C`非常相似。

**双目操作符**

| Operator |     Usage      |     Types      |
| :------: | :------------: | :------------: |
|   `+`    |      add       | String, Number |
|   `-`    |      sub       |     Number     |
|   `*`    |    multiply    |     Number     |
|   `/`    |     divide     |     Number     |
|   `//`   | integer divide |     Number     |
|   `%`    |      mod       |    Integer     |
|   `&`    |  bitwise AND   |    Integer     |
|   `\|`   |   bitwise OR   |    Integer     |
|   `^`    |  bitwise XOR   |    Integer     |
|   `>>`   |  shrift right  |    Integer     |
|   `<<`   |  shrift left   |    Integer     |
|   `<=`   |       LE       | String, Number |
|   `<`    |       LT       | String, Number |
|   `>=`   |       GE       | String, Number |
|   `>`    |       GT       | String, Number |
|   `==`   |       EQ       | String, Number |
|   `!=`   |       NE       | String, Number |
|   `&&`   |  logical AND   |      all       |
|  `\|\|`  |   logical OR   |      all       |
|   `[]`   | `object[key]`  |     Object     |

*注意：字符串和数相加的结果是字符串.*

**三目操作符**

```python
# condition ? true expression : false expression
let result = 1 < 2 ? 1 : 2;
print(result)		# 1
```

**赋值操作符**

| Operator |          Example          |          Usage          |
| :------: | :-----------------------: | :---------------------: |
|   `=`    |  `lhs[,lhs] = rhs[,rhs]`  | `lhs[,lhs] = rhs[,rhs]` |
|   `+=`   | `lhs[,lhs] += rhs[,rhs]`  |    `lhs = lhs + rhs`    |
|   `-=`   | `lhs[,lhs] -= rhs[,rhs]`  |    `lhs = lhs - rhs`    |
|   `*=`   | `lhs[,lhs] *= rhs[,rhs]`  |    `lhs = lhs * rhs`    |
|   `/=`   | `lhs[,lhs] /= rhs[,rhs]`  |    `lhs = lhs / rhs`    |
|   `%=`   | `lhs[,lhs] %= rhs[,rhs]`  |    `lhs = lhs % rhs`    |
|   `&=`   | `lhs[,lhs] &= rhs[,rhs]`  |    `lhs = lhs & rhs`    |
|   `^=`   | `lhs[,lhs] ^= rhs[,rhs]`  |    `lhs = lhs ^ rhs`    |
|  `\|=`   | `lhs[,lhs] \|= rhs[,rhs]` |   `lhs = lhs \| rhs`    |

**操作符优先级**

| Precedence |     Operator      |
| :--------: | :---------------: |
|     13     |     `++` `--`     |
|     12     |    `-` `~` `!`    |
|     11     | `/` `*` `%` `//`  |
|     10     |      `+` `-`      |
|     9      |     `>>` `<<`     |
|     8      | `>` `<` `>=` `<=` |
|     7      |     `==` `!=`     |
|     6      |        `&`        |
|     5      |        `^`        |
|     4      |       `\|`        |
|     3      |       `&&`        |
|     2      |      `\|\|`       |
|     1      |       `?:`        |

### 语句Statements

**If语句**

```python
if (a==1){
	print(1)
}
elif (a==2){
	print(2)
}
else {
	print(3)
}
```

使用`{}`标记一个块。如果块中只包含一个语句的话，我们能省略掉`{}`。上述的代码和下面的一样：

```python
if (a==1)
	print(1)
elif (a==2)
	print(2)
else 
	print(3)
```

这点和`C`语言一致。

**While语句**

```python
while(true){
    # do something
}
```

同理，块里只有一个语句的话，可省略`{}`。

**For语句**

```python
for(let low,high=0,100;low<high;low,high=low+1,high-1) 
    print("hello")

for(let low,high=0,100;low<high;low,high=low+1,high-1) {
    print("hello")
}

let low,high;
for(low,high=0,100;low<high;low,high=low+1,high-1) {
    print("hello")
}
```

能够使用`break`和`continue`关键字进行流程控制：

```python
let sum = 0;
for(let i=0 ;; i++){
    if (i % 2 == 0) continue;
    if (i == 50) break;
    sum += i
}
print(sum)
```

**Switch语句**

switch语句和Go语言很像，编译器会自动在每个case后面添加上break：

```python
switch(a) {
    case 1:
    	print("1")
    case 2:
    	print("2")
    case "str":
    	print("str")
    default:
    	print("others")
}
```

我们能用 `break`关键字提前退出case：

```python
switch(a){
    case 1,2,3:
    	if(a==2) break;
    	print(a)
}
```

能够利用 `fallthrough` 关键字跳转到另外一个case:

```python
switch(a){
    case 1:
    	print("1")
    	fallthrough
    case 2:
    	print("2")
}
```

**Goto语句**

和很多语言一样，`gscript`支持goto关键字：

```python
	let i,sum = 0,0
label:
	sum += i++
    if (i<=100) 
    	goto label
    
    print(sum)		# 5050
```

*注意: 使用goto从一个大作用域跳转到小作用域是不允许的*. 下列代码会编译失败:

```python
# illegal
goto label			# goto statement is in a larger scope
{	
	label:			# label is in a smaller scope
}
```

*注意：try语句的块中也不能使用goto跳转到块之外的label上来*. 

```python
# illegal
try{
    goto label	
}catch(){}

label:
```

尽管在语法层面支持了`goto`，但并不推荐使用。

**Enum Statement**

使用 `enum` 关键字进行枚举，这点很类似于C语言：

```python
enum {
    KW_GOTO,		# 0
    KW_LET		# 1
    KW_FOR=5,		# 5
    KW_WHILE,		# 6
}
```

枚举值的起始值为0，之后的每一个+1，当然你也可以通过枚举赋值改变这个行为如`KW_FOR`。

### 函数

在`gscript`中，函数是一等公民。它能被调用，当作函数参数，也能作为对象的成员。

有两种定义函数的方式：

```python
# 方式1：有名函数
func foo() {
    # do something
}

# 方式2：将匿名函数赋值给一个变量
let foo = func() {
    # do something
}
```

在第一种方式中，foo函数具有全局的作用域，能在当前文件的任何地方被调用。然而，第二种方式只能被用作局部函数。

```python
# compile pass
foo()
func foo(){}		# global function

# compile failed
foo()
let foo = func(){}

# compile pass
let foo = func(){}
foo()
```

*注意：只有在本模块的最顶端作用域定义的有名函数才具有全局的作用域*. Example:

```python
func foo() {}		# foo can be used anywhere

{
    func foo(){}	# 这个会被转化为"let foo = func(){}"
}
```

**多返回值**

和Go语言一样，`gscript`支持多返回值：

```python
func foo(a,b){
    return a//b, a%b
}

let a,b = foo(9,2) 		# a=4,b=1
let c = foo(9,2)		# c=4			抛弃多余的返回值
let d,e,f = foo(9,2)		# d=4,e=1,f=nil		不够的部分用nil填补
```

注意一些特殊情况:

```python
let a,b,c = 1,foo(9,2)		# a=1,b=4,c=1
let a,b,c = foo(9,2),9		# a=4,b=9,c=nil
print(foo(9,2))			# output: 4
```

这些情况很类似与 `lua`语言。

如果我们将一个函数调用的返回值当作另外一个函数的参数时，只会把第一个返回值作为参数传递，其他的返回值会被抛弃掉。

**闭包**

`gscript` 支持闭包：

```python
func foo() {
    let i = 0;
    # 返回一个闭包
    return func(){		
        return ++i;		# 捕获外部作用域的标量
    }
}

let f = foo();

print(f())				# output: 1
print(f())				# output: 2
```

**默认参数**

`gscript` 允许函数参数具有默认值。如果函数调用时没有指定参数，则会以默认值进行填充

```
func foo(name="jack",age=10) {
	print(namea,age)
}

foo()			# output: jack 10
foo("rose")		# output: rose 10
foo("rose",20)		# output: rose 20
```

**可变参数**

`gscript` 允许函数接收可变参数：

```python
func foo(a,...args) {
    print(a,len(args),args)
}

foo(1,2,3)			# output: 1 2 [2,3]
foo(1)				# output: 1 0 []
foo(1,2,3,4)			# output: 1 3 [2,3,4]
foo()				# output: nil,0,[]
```

### Class

`gscript` 支持面向对象：

```python
class People{
    # 构造函数
    __self(name,age) {
       	this.name = name
        this.age = age
    }
    # 方法
    show(){
        print("age of",this.name,"is",this.age)
        print(this.name,"lives at",this.planet)
    }
    
    # 另外一种定义方法的方式
    change_name = func(name) {
        this.name = name
    }
    
    # 给planet成员设置默认值
    planet = "earth"
    
    # 下面两个会被编译器忽略掉，但写出来可以帮助别人知道该对象有哪些成员
    name
    age			
}
```

+ `__self`时构造函数
+ 使用`this`代表对象

如何使用对象：

```python
let p = new People("Jack",18);			# use keyword new
p.show();					# call method
p.change_name("Rose");				# call method
p.age = 20;					# access member
p.show();					# call method
```

### 异常处理

不像Go语言，`gscript`使用try-catch机制处理异常：

```python
import fs

try{
    # readFile may throw an exception
    let data = fs.readFile("foo.txt")
    print(data.toString())
}
catch(e){
    print("operation failed:",e)
}
```

使用内置函数`throw`抛出异常。

```python
try{
    throw("this is a exception")
}
catch(e){
    print(e)
}
```

如果异常在向上抛的过程中，没有被任何try catch捕获，则整个程序会崩掉。

### 模块

模块是`gscript`中最基本的编译单元。一个模块可以使用`import`关键字引用另外一个模块：

假设我们的工作目录如下:

```
.
├── lib
│   └── xxx.gs
├── main.gs
└── module.gs
```

`main.gs` 是主模块和程序入口。 `lib/lib.gs` 和 `module.gs` 是待导入的自定义模块。

`main.gs`:

```python
import "module"		# 导入自定义模块module.gs
import "lib/xxx"	# 导入自定义模块lib/xxx.gs
import fs		# 导入标准库fs

module.foo();	# 调用模块的方法
```

导入的模块名是否用引号标注代表了导入的模块是不是标准库。

我们能给导入的模块取别名：

```python
import "module" as mod
```

模块文件可以使用`export`关键字返回一个表达式，例`module.gs`：

```python
func foo(){
    print("foo")
}

# syntax: export Expression
# 这里导出了一个对象
export {
    foo: foo,
}
```

### Comments

目前只支持行注释：

```python
# line comments
```

