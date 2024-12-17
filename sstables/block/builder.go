package block

import (
	"encoding/binary"
	"lsm/entries"
	"unsafe"
)

var ReservedKeySize = int(unsafe.Sizeof(uint16(0)))
var ReservedValueSize = int(unsafe.Sizeof(uint16(0)))
var Uint16Size = int(unsafe.Sizeof(uint16(0)))
var Uint32Size = int(unsafe.Sizeof(uint32(0)))

const (
	kb               uint = 1024
	DefaultBlockSize      = 4 * kb
)

type Builder struct {
	blockSize           uint
	data                []byte
	firstKey            entries.Key
	keyValueBeginOffset []uint16
}

func NewBuilder(blockSize uint) *Builder {
	return &Builder{
		blockSize: blockSize,
		data:      make([]byte, 0, blockSize),
	}
}

func (b *Builder) Add(key entries.Key, value entries.Value) bool {
	if uint(b.size()+key.EncodedKeySizeInBytes()+value.Size()+Uint16Size*2) > b.blockSize {
		return false
	}

	if b.firstKey.RawKeySize() == 0 {
		b.firstKey = key
	}

	b.keyValueBeginOffset = append(b.keyValueBeginOffset, uint16(len(b.data)))
	buffer := make([]byte, ReservedKeySize+ReservedValueSize+key.EncodedKeySizeInBytes()+value.Size())
	binary.LittleEndian.PutUint16(buffer[:], uint16(key.EncodedKeySizeInBytes()))
	copy(buffer[ReservedKeySize:], key.EncodedBytes())
	binary.LittleEndian.PutUint16(buffer[ReservedKeySize+key.EncodedKeySizeInBytes():], uint16(value.Size()))
	copy(buffer[ReservedKeySize+key.EncodedKeySizeInBytes()+ReservedValueSize:], value.Bytes())
	b.data = append(b.data, buffer...)
	return true

}

func (b *Builder) isEmpty() bool {
	return len(b.keyValueBeginOffset) == 0 || len(b.data) == 0
}

func (b *Builder) Build() Block {
	if b.isEmpty() {
		panic("block is empty")
	}

	return NewBlock(b.data, b.keyValueBeginOffset)
}

func (b *Builder) size() int {
	return len(b.data) + len(b.keyValueBeginOffset)*Uint16Size + Uint16Size
}
