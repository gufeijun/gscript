## The Gscript Language

中文|[English](https://github.com/gufeijun/gscript/blob/master/README.md)

Gscript是一个用纯Go编写的轻量动态脚本语言。

案例:

```python
import fs
import os

# 定义函数
func writeSomeMsg(filepath) {
    # 以只写模式打开文件，如果文件不存在则创建，如果文件存在则清空文件
    let file = fs.open(filepath,"wct");

    # 往文件写入一些数据
    file.write("hello world!\n")
    file.write("gscript is a good language!")

    # 关闭文件
    file.close();
}

let filepath = "./text.txt"

try{
    writeSomeMsg(filepath)

    let stat = fs.stat(filepath)
    print("size of text.txt is " + stat.size + "B");
    
    # 读取文件所有内容
    let data = fs.readFile(filepath);
    print("message of", filepath, "is:")
    print(data.toString())

    # 删除文件
    fs.remove(filepath)
    print("success!")
}
catch(e){
    print("operation failed, error msg:",e)
    os.exit(0)
}
```

### 特性

+ 函数

  + 支持多返回值

  + 支持闭包
  + 支持函数递归

+ 脚本语言的所有标准库以静态资源的方式打包成单独一个二进制文件，gscript部署以及使用极为方便

+ 支持面向对象

+ 支持Debug模式调试代码

+ 能够生成高可读性的类汇编代码

+ 语法简单易学，特别是对于那些具有js和python经验的人

+ 多文件、模块化支持

+ 支持try, catch, throw异常处理机制

+ 编译生成字节码，使用VM执行

+ 完善的变量作用域机制

### 安装

由于gscript的编译器以及VM全部使用纯GO开发，不依赖CGO，未使用平台相关的接口，所以gscript本身是跨平台的语言。

两种安装gscript方式:

+ 源码编译

  *注意：由于使用了静态资源嵌入(embed)特性，所以保证你的GO编译器版本大于等于1.16*.

  ```shell
  git clone git@github.com:gufeijun/gscript.git
  cd gscript
  sh build.sh
  ```

  `bin/gsc`就是生成的编译器二进制文件。

+ 从[releases](https://github.com/gufeijun/gscript/releases)下载。

`gsc ` 意思是 `gscript complier`。你可以自主选择是否将这个二进制文件添加到`PATH`中。

### 快速开始

在 `main.gs`写如下代码:

```python
print("hello world");
```

运行脚本t:

```shell
gsc run main.gs
```

获得如下输出:

```
hello world
```

你也能写更复杂的代码，计算fibonacci序列:

```python
print(fib(20))

func fib(x) {
    if (x == 0) return 0;
    if (x == 1) return 1;
    return fib(x-1) + fib(x-2)
}
```

运行脚本，输出6765.

### 使用

我们上面已经演示了`gsc run`命令，你可以通过`gsc --help`获取更多使用帮助。

+ 使用 `gsc debug <source file>`  或者 `gsc debug <bytecode file>`进入调试模式，这个模式下能一行行执行字节码，并且能实时查看栈空间以及变量表的变化。

+ 使用 `gsc build <source file>` 生成字节码，可以额外指定`-o`这个参数指定输出文件的名字。

+ 使用 `gsc build -a <source file>` 或者 `gsc build -a <bytecode file>` 去生成高可读性的类汇编代码。

  以上节的hello world脚本为例，运行如下命令：

  ```shell
  gsc build -a main.gs -o output.gscasm
  ```

  会生成文件 `output.gsasm`:

  ```
  /home/xxx/gscript/main.gs(MainProto):
  	MainFunc:
  		0		LOAD_CONST "hello world"
  		9		LOAD_BUILTIN "print"
  		14		CALL 0 1
  		17		STOP
  ```

+ 使用 `gsc run <source file>` 或者 `gsc run <bytecode file>`运行脚本或者字节码。

### 参考

+ [语法](https://github.com/gufeijun/gscript/blob/master/doc/syntax_zh.md)
+ [内置函数](https://github.com/gufeijun/gscript/blob/master/doc/builtin_zh.md)
+ [标准库](https://github.com/gufeijun/gscript/blob/master/doc/std_zh.md)

