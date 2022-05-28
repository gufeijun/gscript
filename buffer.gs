class Buffer{
    __self(cap, str) {
        if (cap == nil) return;
        if (str != nil) this._buffer = __buffer_from(str);
        else this._buffer = __buffer_new(cap);
        this._cap = cap;
    }
    cap() {
        return this._cap;
    }
    toString(offset=0, length=-1) {
        if (length == -1) 
            length = this._cap - offset;
        return __buffer_toString(this._buffer, offset, length);
    }
    slice(offset, length) {
        if (length == nil) 
            length = len(this._buffer) - offset;
        let buf = new Buffer;
        buf._buffer = __buffer_slice(this._buffer, offset, length);
        buf._cap = length;
        return buf;
    }
    readInt8(offset) {
        return __buffer_readNumber(this._buffer, offset, 1, true, false, false);
    }
    readUint8(offset) {
        return __buffer_readNumber(this._buffer, offset, 1, false, false, false);
    }
    readInt16BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 2, true, false, false);
    }
    readInt16LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 2, true, true, false);
    }
    readUint16BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 2, false, false, false);
    }
    readUint16LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 2, false, true, false);
    }
    readInt32BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 4, true, false, false);
    }
    readInt32LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 4, true, true, false);
    }
    readUint32BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 4, false, false, false);
    }
    readUint32LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 4, false, true, false);
    }
    readInt64BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 8, true, false, false);
    }
    readInt64LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 8, true, true, false);
    }
    readUint64BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 8, false, false, false);
    }
    readUint64LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 8, false, true, false);
    }
    readFloat32LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 4, false, true, true);
    }
    readFloat32BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 4, false, false, true);
    }
    readFloat64LE(offset) {
        return __buffer_readNumber(this._buffer, offset, 8, false, true, true);
    }
    readFloat64BE(offset) {
        return __buffer_readNumber(this._buffer, offset, 8, false, false, true);
    }
    write8(offset, number) {
        __buffer_writeNumber(this._buffer, offset, 1, false, number);
    }
    write16BE(offset, number) {
        __buffer_writeNumber(this._buffer, offset, 2, false, number);
    }
    write16LE(offset, number) {
        __buffer_writeNumber(this._buffer, offset, 2, true, number);
    }
    write32BE(offset, number) {
        __buffer_writeNumber(this._buffer, offset, 4, false, number);
    }
    write32LE(offset, number) {
        __buffer_writeNumber(this._buffer, offset, 4, true, number);
    }
    write64BE(offset, number) {
        __buffer_writeNumber(this._buffer, offset, 8, false, number);
    }
    write64LE(offset, number) {
        __buffer_writeNumber(this._buffer, offset, 8, true, number);
    }
}

export {
    alloc: func(cap) {
        return new Buffer(cap);
    },
    from: func(str) {
        return new Buffer(len(str),str);
    },
    concat: func(buf1, buf2) {
        let buf = new Buffer;
        buf._buffer = __buffer_concat(buf1._buffer, buf2._buffer);
        buf._cap = buf1.cap() + buf2.cap();
        return buf;
    },
    copy: func(buf1, buf2, length, start1=0, start2=0) {
        __buffer_copy(buf1._buffer, buf2._buffer, length, start1, start2);
    }
}
