## Standard Library - fs

```python
import fs
```

### Functions

+ `open(path:String, flag="r":String, mode=0664:Integer) exception => class File `: open named file with flag and mode. 
  + flag(`r`): read-only.
  + flag(`w`): write-only.
  + flag(`wr`): read write.
  + flag(`c`): if file do not exist, create file.
  + flag(`a`): append.
  + flag(`t`): truncate.
  + flag(`e`): excl .
+ `create(path:String, mode=0664:Integer) exception => class File`: create file with mode.
+ `stat(path:String) exception => class stat`: get file stat of named path. member of class stat:
  + `is_dir:Bool`. If the file is a directory.
  + `mode:Integer`. File mode.
  + `name:String`. File name.
  + `size:Integer`. File Size.
  + `mod_time`. File modify unix time.
+ `remove(path:String) exception`: remove the named file.
+ `readFile(path:String) exception => class Buffer`:  read all the data from named path to Buffer.
+ `mkdir(dir:String, mode=0664:Integer) exception`: make directory with named dir.
+ `chmod(path:String, mode:Integer) exception`: changes the mode of the named file to mode.
+ `chown(path:String, uid:Integer, gid:Integer) exception`: changes the numeric uid and gid of the named file.
+ `rename(oldpath:String, newpath:String) exception`: renames (moves) oldpath to newpath.
+ `readDir(path:String) exception => Array<class stat>`: reads the named directory, returning all its directory entries sorted by filename.

### File

method:

+ `read(buf:class Buffer, size:Integer) exception => Integer`: read @size bytes from file to Buffer. If size is not specific, it will try to fill up Buffer. It returns the number of bytes read.

+ `write(data:string or class Buffer, size=-1:Integer) exception => Integer`: write @size bytes from data to file. If size==-1, write all bytes of data to file.

+ `close() exception`: close file.

+ `seek(offset:Integer, whence="cur":String) exception => Integer` : Seek sets the offset for the next Read or Write on file to offset, interpreted according to whence: "start" means relative to the origin of the file, "cur" means relative to the current offset, and "end" means relative to the end. It returns the new offset .

+ `chmod(mode) exception `: changes the mode of the file to mode.

+ `chown(uid, gid) exception`: changes the mode of the file to mode.

+ `chdir() exception`: changes the current working directory to the file.

+ `stat() exception => class stat`: get file stat. 

+ `isDir() exception => Boolean`: if the file is directory.

+ `readDir(n:Integer) exception => Array<class stat>`:  reads the contents of the directory associated with the file f and returns a slice of DirEntry values in directory order. Subsequent calls on the same file will yield later DirEntry records in the directory.
  + If n > 0, ReadDir returns at most n DirEntry records. In this case, if ReadDir returns an empty slice, it will return an error explaining why. At the end of a directory, the error is io.EOF.
  + If n <= 0, ReadDir returns all the DirEntry records remaining in the directory.
