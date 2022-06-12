## Standard Library - os

```python
import os
```

### Functions

+ `chdir(dir:String) exception`: changes the current working directory to the named directory.
+ `exit(code:Integer)`: causes the current program to exit with the given status code.
+ `getEnv(key:String) => String`: retrieves the value of the environment variable named by the key.
+ `setEnv(key:String,val:String) exception`: sets the value of the environment variable named by the key.
+ `args() => Array<String>`: returns command-line arguments, starting with the program name.
+ `getegid() => Integer`: returns the numeric effective group id of the caller.
+ `geteuid() => Integer`:  returns the numeric effective user id of the caller.
+ `getgid() => Integer`: returns the numeric group id of the caller.
+ `getpid() => Integer`: returns the process id of the caller.
+ `getppid() => Integer`: returns the process id of the caller's parent.
+ `getuid() => Integer`: returns the numeric user id of the caller.
+ `exec(cmd:String, ...args:String) => String`: execute the named program with the given arguments. return the output with protgram.

