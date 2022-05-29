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
}