## Gscript Syntax

### Type

Following are most frequently used types in gscript:

```python
100						# Number
100.21					# Number
"gscript"				# String
false					# Boolean
[1,"gs",{}]				# Array
{foo: 1, bar:"bar"}		# Objet
func(){}				# Closure
nil						# nil
```

**Array**

```python
let arr = [1,2,[1,2,3]]
arr[0]					# 1
arr[2][2]				# 3
arr[3]					# will panic
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

### Variables and Scopes

Use  keyword `let` to define a new local variable. Using undefined variable will make compile failed: 

```python
let a = 1

# enter a new scope
{
    print(a)			# a==1
    let a = 2
    let b = 3
    print(a)			# a==2
    print(b)			# b==3
}

print(a)				# a==1
print(b)				# invalid, complie will not pass
```

Every variable can only be used in a scope, using a variable outside the scope will make compile failed, too.  

Like many other script language, a variable can be assigned to a new value of different type. 

```python
let a = {}
a = []					# ok
```

We can define multiple variables using one `let`:

```python
let a, b = 1, "hello"	# a=1,b="hello"
let a, b = 1			# a=1,b=nil
let a = 1,"hello"		# a=1
```

Of course, we can assign to multiple variables use one `=`:

```python
let a,b
a, b = 1, 2 
```

### Operators

**Unary Operators**

| Operator |      Usage       |         Types          |
| :------: | :--------------: | :--------------------: |
|   `~`    | bitwise negation |        Integer         |
|   `-`    |     `0 - x`      | Number(Integer, Float) |
|   `!`    |       `!x`       |          all           |
|   `--`   |  `--x` or `x--`  |         Number         |
|   `++`   |  `++x` or `x++`  |         Number         |

Be careful to use `--` and `++`.

If these two operators are used as a statement:

```python
# x++ is the same as ++x
let x = 0
x++			# x==1
++x			# x==2

# x-- is the same as --x
x--			# x==1
--x			# x==0
```

But if these two operators are used as a expression:

```python
let x = 0
print(x++)	# output: 0
print(++x)	# output: 2
print(--x)	# output: 1
print(x--)	# output: 1
# now x==0
```

The behaviors of `--` and `++` are the same as language `JavaScript` and `C`.

**Binary Operators**

| Operator |     Usage      |     Types      |
| :------: | :------------: | :------------: |
|   `+`    |      add       | String, Number |
|   `-`    |      sub       |     Number     |
|   `*`    |    multiply    |     Number     |
|   `/`    |     divide     |     Number     |
|   `//`   | integer divide |     Number     |
|   `%`    |      mod       |    Integer     |
|   `&`    |  bitwise AND   |    Integer     |
|   `|`    |   bitwise OR   |    Integer     |
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
|   `||`   |   logical OR   |      all       |
|   `[]`   | `object[key]`  |     Object     |

*note: output type of  (number +  string) will be a string.*

**Ternary Operators**

```python
# condition ? true expression : false expression
let result = 1 < 2 ? 1 : 2;
print(result)		# 1
```

**Assignment Operators**

| Operator |         Example          |          Usage          |
| :------: | :----------------------: | :---------------------: |
|   `=`    | `lhs[,lhs] = rhs[,rhs]`  | `lhs[,lhs] = rhs[,rhs]` |
|   `+=`   | `lhs[,lhs] += rhs[,rhs]` |    `lhs = lhs + rhs`    |
|   `-=`   | `lhs[,lhs] -= rhs[,rhs]` |    `lhs = lhs - rhs`    |
|   `*=`   | `lhs[,lhs] *= rhs[,rhs]` |    `lhs = lhs * rhs`    |
|   `/=`   | `lhs[,lhs] /= rhs[,rhs]` |    `lhs = lhs / rhs`    |
|   `%=`   | `lhs[,lhs] %= rhs[,rhs]` |    `lhs = lhs % rhs`    |
|   `&=`   | `lhs[,lhs] &= rhs[,rhs]` |    `lhs = lhs & rhs`    |
|   `^=`   | `lhs[,lhs] ^= rhs[,rhs]` |    `lhs = lhs ^ rhs`    |
|   `|=`   | `lhs[,lhs] |= rhs[,rhs]` |    `lhs = lhs | rhs`    |

**Operator Precedences**

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
|     4      |        `|`        |
|     3      |       `&&`        |
|     2      |       `||`        |
|     1      |       `?:`        |

### Statements

**If Statement**

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

Use `{}` to mark a block. If there is only one statement in block, we can omit `{}`。 Code above is the same as following:

```python
if (a==1)
	print(1)
elif (a==2)
	print(2)
else 
	print(3)
```

**While Statement**

```python
while(true){
    # do something
}
```

If there is only one statement in block, we can omit `{}`.

**For Statement**

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

we can use keyword `break` or `continue` to control process:

```python
let sum = 0;
for(let i=0 ;; i++){
    if (i % 2 == 0) continue;
    if (i == 50) break;
    sum += i
}
print(sum)
```

**Switch statement**

switch statement is very similar to Go, it will automatically insert a  break at last of every case:

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

we can use keyword `break` to exit switch case early:

```python
switch(a){
    case 1,2,3:
    	if(a==2) break;
    	print(a)
}
```

we can jump to another case by using keyword `fallthrough` :

```python
switch(a){
    case 1:
    	print("1")
    	fallthrough
    case 2:
    	print("2")
}
```

**Goto Statement**

Like many languages, gscript can use keyword "goto":

```python
	let i,sum = 0,0
label:
	sum += i++
    if (i<=100) 
    	goto label
    
    print(sum)		# 5050
```

*note: use "goto" to jump inside the scope of a label from outside is now allowed*. Following code will compile failed:

```python
# illegal
goto label			# goto statement is in a larger scope
{	
	label:			# label is in a smaller scope
}
```

*note: use "goto" to jump outside of try block is now allowed too*. 

```python
# illegal
try{
    goto label	
}catch(){}

label:
```

Although we support goto at the syntactic level, it is not recommended for use in any case.

**Enum Statement**

We can use keyword `enum` to enumerate some constants, very similar to `C`:

```python
enum {
    KW_GOTO,		# 0
    KW_LET			# 1
    KW_FOR=5,		# 5
    KW_WHILE,		# 6
}
```

enum value starts at 0. Their values are determined at compile time.

### Function

Functions are first class citizens in`gscript`. It can be called, passed as a parameter, or as a member of an object.

There are two ways to define a function:

```python
# way 1: named function
func foo() {
    # do something
}

# way 2: assign an anonymous function to a variable
let foo = func() {
    # do something
}
```

what's the differences?

In first way, function `foo` has a global scope, it can be used anywhere. However, `foo` in second way is only used as a local variable. 

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

*note: only named function at top of scope in the source file is global*. Example:

```python
func foo() {}		# foo can be used anywhere

{
    func foo(){}	# only can used in this scope, equals to "let foo = func(){}"
}
```

**Multiple Return Values**

Similar to Go, a function can return Multiple values:

```python
func foo(a,b){
    return a//b, a%b
}

let a,b = foo(9,2) 		# a=4,b=1
let c = foo(9,2)		# c=4				will discard extra return values
let d,e,f = foo(9,2)	# d=4,e=1,f=nil		the missing part will be filled with nil
```

Pay attention to some special cases:

```python
let a,b,c = 1,foo(9,2)	# a=1,b=4,c=1
let a,b,c = foo(9,2),9	# a=4,b=9,c=nil
print(foo(9,2))			# output: 4
```

These are much like language `lua`.

If we pass return values of a function call as another function call arguments, only first return value will be taken, others will be discard.

**closure**

`gscript` supports closure, which is widely used in script language.

```python
func foo() {
    let i = 0;
    # return a closure
    return func(){		
        return ++i;		# capture value outside of scope
    }
}

let f = foo();

print(f())				# output: 1
print(f())				# output: 2
```

**Default arguments**

`gscript` allows function arguments to have default values. If the function is called without the argument, the argument gets its default value.

```
func foo(name="jack",age=10) {
	print(namea,age)
}

foo()			# output: jack 10
foo("rose")		# output: rose 10
foo("rose",20)	# output: rose 20
```

**varargs**

`gscript` allows function receive a variable number of arguments.

```python
func foo(a,...args) {
    print(a,len(args),args)
}

foo(1,2,3)			# output: 1 2 [2,3]
foo(1)				# output: 1 0 []
foo(1,2,3,4)		# output: 1 3 [2,3,4]
foo()				# output: nil,0,[]
```

### Class

`gscript` support Object-oriented programming.

```python
class People{
    # construction function
    __self(name,age) {
       	this.name = name
        this.age = age
    }
    # method
    show(){
        print("age of",this.name,"is",this.age)
        print(this.name,"lives at",this.planet)
    }
    
    # another way to define a method
    change_name = func(name) {
        this.name = name
    }
    
    # default value for member planet
    planet = "earth"
    
    # name and age will be ignored by compiler, but useful for people to kown which members this class has
    name
    age			
}
```

+ `__self` is the constructor of a class
+ use `this` to access object.

how to create an object?

```python
let p = new People("Jack",18);		# use keyword new
p.show();							# call method
p.change_name("Rose");				# call method
p.age = 20;							# access member
p.show();							# call method
```

### Exception Handle

Unlike Go, `gscript` use keyword `try`, `catch` to capture exception and builtin function`throw` to throw an exception.

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

Throw an exception:

```python
try{
    throw("this is a exception")
}
catch(e){
    print(e)
}
```

If an exception is not caught, the whole program will crash.

### Module

Module is the basic compilation unit in `gscript`. A module can import another module using `import` keyword. 

Consume our work directory is like this:

```
.
├── lib
│   └── xxx.gs
├── main.gs
└── module.gs
```

`main.gs` is our main module and general entry of program. `lib/lib.gs` and `module.gs` are local modules to be imported.

`main.gs`:

```python
import "module"		# import local module "module.gs"
import "lib/xxx"	# import local module "lib/xxx.gs"
import fs			# import standard library "fs"

module.foo();
```

Whether the imported module is quoted to indicate whether it is a standard library or a local module.

Use keyword `export` to return a value, `module.gs`:

```python
func foo(){
    print("foo")
}

# export Expression
# here we export an object with one member foo
export {
    foo: foo,
}
```

### Comments

For now, we only support line comments:

```python
# line comments
```

