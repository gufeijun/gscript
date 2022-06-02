import "./Buffer";

class File{
    __self(file) {
        this._file = file;
    }
    read(buf,size) {
        if(size == nil) size = buf.cap();
        return __read(this._file, buf._buffer, size);
    }
    # data is a Buffer or String
    write(data, size=-1) {
        if (type(data) != "String") {
            data = data._buffer;
        }
        size = size == -1 ? len(data) : size;
        return __write(this._file, data, size);
    }
    close() {
        __close(this._file);
    }
    # whence = cur, end or start
    seek(offset, whence="cur") {
        return __seek(this._file, offset, whence);
    }
    chmod(mode) {
        __fchmod(this._file, mode);
    }
    chown(uid, gid) {
        __fchown(this._file, uid, gid);
    }
    chdir() {
        __fchdir(this._file);
    }
    stat() {
        if (this._stat == nil)
            this._stat = __fstat(this._file); 
        return this._stat;
    }
    isDir() {
        if this._stat == nil 
            this.stat();
        return this._stat.is_dir;
    }
    readDir(n) {
        return __freaddir(this._file, n);
    }
}

export {
    open: func(path, flag="r", mode=0664) {
        let file = __open(path, flag, mode);
        return new File(file);
    },
    create: func(path, mode=0664) {
        let file = __open(path, "crw", mode);
        return new File(file);
    },
    stat: func(path) {
        return __stat(path);
    },
    remove: func(path) {
        __remove(path);
    },
    readFile: func(path) {
        let size = __stat(path).size;
        let buf = Buffer.alloc(size);
        let file = __open(path, "r", 0);
        let n = __read(file, buf._buffer, size);
        __close(file);
        return buf, n;
    },
    mkdir: func(dir, mode=0664) {
        __mkdir(dir, mode);
    },
    chmod: func(path, mode) {
        __chmod(path, mode);
    },
    chown: func(path, uid, gid) {
        __chown(path, uid, gid);
    },
    rename: func(oldpath, newpath) {
        __rename(oldpath, newpath);
    },
    readDir: func(path) {
        return __readdir(path);
    },
}