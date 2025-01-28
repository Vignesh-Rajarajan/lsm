package block

import (
	"encoding/binary"
	"lsm/entries"
)

// Block data: [key1:value1][key2:value2][key3:value3]
// Offsets: [0,10,20]
// To find key2:value2:
// Read offset[1] = 10
// Jump directly to position 10 in block
// Read key-value pair without scanning through key1
// Offsets: [0,100,200,...,9900]
// Searching for key500:
// - Check middle offset (position 4500)
// - If key500 > middle, search upper half
// - If key500 < middle, search lower half
type Block struct {
	data                []byte
	keyValueBeginOffset []uint16
}

func NewBlock(data []byte, keyValueBeginOffset []uint16) Block {
	return Block{data: data, keyValueBeginOffset: keyValueBeginOffset}
}

// Encode encodes the block into a byte slice.
// it looks like this:
// ---------------------------------------------------------------------------------------------------
// | encoded key-value pairs | encoded key-value pairs |...| 0 | 48 | 120| ... |3088|    2 bytes     |
// ---------------------------------------------------------------------------------------------------
// <-------------Encoded data------------------------------><-Begin of offset keys--><--No of Offsets->
func (b Block) Encode() []byte {
	data := b.data
	data = append(data, b.encodeKeyValueBeginOffset()...)
	numberOfKeyValOffset := make([]byte, Uint16Size)
	binary.LittleEndian.PutUint16(numberOfKeyValOffset, uint16(len(b.keyValueBeginOffset)))
	data = append(data, numberOfKeyValOffset...)
	return data
}

// DecodeBlockData : take the last 2 bytes to get the number of offsets
// The size of keyValueBeginOffset is the number of keyValueBeginOffsets * 2 bytes (offset is in uint16)
// the (length of the data block) minus  (size of the keyValueBeginOffsets) minus (2 bytes for the number of offsets)
// is the start of the data or the end of key-value encoded data
func (b Block) DecodeBlockData(data []byte) Block {
	numberOfOffsets := binary.LittleEndian.Uint16(data[len(data)-Uint16Size:])
	startOfOffsets := uint16(len(data)) - uint16(Uint16Size) - numberOfOffsets*uint16(Uint16Size)
	offsetBuffer := data[startOfOffsets : len(data)-Uint16Size]

	keyValBeginOffset := make([]uint16, 0, numberOfOffsets)
	for i := 0; i < len(offsetBuffer); i += Uint16Size {
		keyValBeginOffset = append(keyValBeginOffset, binary.LittleEndian.Uint16(offsetBuffer[i:]))
	}
	return Block{data: data[:startOfOffsets], keyValueBeginOffset: keyValBeginOffset}
}

func (b Block) encodeKeyValueBeginOffset() []byte {
	//Each offset takes only 2 bytes
	offsetBuffer := make([]byte, Uint16Size*len(b.keyValueBeginOffset))
	offsetIndex := 0
	for _, offset := range b.keyValueBeginOffset {
		binary.LittleEndian.PutUint16(offsetBuffer[offsetIndex:], offset)
		offsetIndex += Uint16Size
	}
	return offsetBuffer
}

func (b Block) SeekToFirst() *Iterator {
	it := &Iterator{
		block:     b,
		offsetIdx: 0,
	}
	it.seekToOffsetIndex(it.offsetIdx)
	return it
}

func (b Block) Seek(key entries.Key) *Iterator {
	it := &Iterator{
		block: b,
	}
	it.seekToEqualOrGreater(key)
	return it
}
