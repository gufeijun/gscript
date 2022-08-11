## The Gscript Language

English|[中文](https://github.com/gufeijun/gscript/blob/master/README_zh.md)

**Gscript is a light, dynamic script language written in Go**.

example:

```python
import fs
import os

# define a function
func writeSomeMsg(filepath) {
    # open, create and truncate text.txt in write-only mode
    let file = fs.open(filepath,"wct");

    # write some message into file
    file.write("hello world!\n")
    file.write("gscript is a good language!")

    # close file
    file.close();
}

let filepath = "./text.txt"

try{
    writeSomeMsg(filepath)

    let stat = fs.stat(filepath)
    print("size of text.txt is " + stat.size + "B");
    
    # read all data from text.txt
    let data = fs.readFile(filepath);
    print("message of", filepath, "is:")
    print(data.toString())

    # remove text.txt
    fs.remove(filepath)
    print("success!")
}
catch(e){
    print("operation failed, error msg:",e)
    os.exit(0)
}
```

### Fetures

+ function

  + Multiple return values

  + Closure
  + Recursive function call

+ All standard libraries are packaged into one executable, which means, to install or use gscipt it's no need to configure some  environment variables.

+ Object-oriented programming

+ Debug mode

+ Can generate human-readable assemble-like codes

+ Simple syntax, easy to learn, especially to coders with JavaScript or Python  experiences.

+ Modular support

+ Exception support: try, catch and throw

+ Complied/executed as bytecode on stack-based VM

### Install

Since compiler and VM are written in pure Go(no cgo) and does not use any platform-dependent interfaces, so Gscript is a cross-platform language.

You can install Gscript in two ways:

+ Complie from source code. 

  *note: we use the new feature "embed" of Go 1.16, so make sure your Go version is greater than or equal to 1.16*.

  ```shell
  git clone git@github.com:gufeijun/gscript.git
  cd gscript
  sh build.sh
  ```

  then the compiler will generated to `bin/gsc`. 

+ Download from [releases](https://github.com/gufeijun/gscript/releases). 

`gsc ` means `gscript compiler`. You can add the executable to `PATH` as you wish. 

Then all you need is just having fun.

### Quick Start

open file `main.gs` and write codes below:

```python
print("hello world");
```

run the srcipt:

```shell
gsc run main.gs
```

you will get following output:

```
hello world
```

you can also do something more challenging:

```python
print(fib(20))

func fib(x) {
    if (x == 0) return 0;
    if (x == 1) return 1;
    return fib(x-1) + fib(x-2)
}
```

run the script, you will get 6765.

### Usage

we demonstrated the command `run` of `gsc` above. You can use `gsc --help` for more details about how to use gsc.

+ use `gsc debug <source file>`  or `gsc debug <bytecode file>`to enter debug mode, in which we can execute virtual instruction one by one and  view stack and variable table changes in real time.

+ use `gsc build <source file>` to generate bytecode, besides, you can use `-o` flag to specific name of output file.

+ use `gsc build -a <source file>` or `gsc build -a <bytecode file>` to generate human-readable assemble-like codes.  

  Take the `main.gs` in section `Quic Start` as an example, run the following command:

  ```shell
  gsc build -a main.gs -o output.gscasm
  ```

  It will generate `output.gsasm`:

  ```
  /home/xxx/gscript/main.gs(MainProto):
  	MainFunc:
  		0		LOAD_CONST "hello world"
  		9		LOAD_BUILTIN "print"
  		14		CALL 0 1
  		17		STOP
  ```

+ use `gsc run <source file>` or `gsc run <bytecode file>` to run the script.

### References

+ [Language Syntax](https://github.com/gufeijun/gscript/blob/master/doc/syntax.md)
+ [Builtin Functions](https://github.com/gufeijun/gscript/blob/master/doc/builtin.md)
+ [Standard Library](https://github.com/gufeijun/gscript/blob/master/doc/std.md)
+ [BNF](https://github.com/gufeijun/gscript/blob/master/doc/bnf.txt)

