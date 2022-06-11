## Standard Library - Buffer

```python
import Buffer
```

### Functions

+ `alloc(cap:Integer) => class Buffer`: make a bytes Buffer with capacity.
+ `from(str:String) => class Buffer`: make a bytes Buffer from string.
+ `concat(buf1:class Buffer, buf2:class Buffer) => class Buffer` : concat two Buffer.
+ `copy(buf1:class Buffer, buf2:class Buffer, length:Integer, start1=0:Integer, start2=0:Integer)`: copy @length bytes from buf2 to buf1.

### Buffer

method:

+ `cap() => Integer`: return capacity of  Buffer
+ `toString() => String`: bytes data to string
+ `slice(offset:Integer, length:Integer) => class Buffer`: get slice of Buffer. Start from offset, ends at offset+length. if length is not specific, ends at Buffer.cap().
+ `readInt8(offset:Integer) => Integer`.
+ `readUint8(offset:Integer) => Integer`.
+ `readInt16BE(offset:Integer) => Integer`.
+ `readInt16LE(offset:Integer) => Integer`.
+ `readUint16LE(offset:Integer) => Integer`.
+ `readUint16BE(offset:Integer) => Integer`.
+ `readInt32LE(offset:Integer) => Integer`.
+ `readInt32BE(offset:Integer) => Integer`.
+ `readUint32LE(offset:Integer) => Integer`.
+ `readUint32BE(offset:Integer) => Integer`.
+ `readInt64LE(offset:Integer) => Integer`.
+ `readInt64BE(offset:Integer) => Integer`.
+ `readUint64LE(offset:Integer) => Integer`.
+ `readUint64BE(offset:Integer) => Integer`.
+ `write8(offset:Integer, number:Number)`.
+ `write16BE(offset:Integer, number:Number)`.
+ `write16LE(offset:Integer, number:Number)`.
+ `write32BE(offset:Integer, number:Number)`.
+ `write32LE(offset:Integer, number:Number)`.
+ `write64LE(offset:Integer, number:Number)`.
+ `write64BE(offset:Integer, number:Number)`.