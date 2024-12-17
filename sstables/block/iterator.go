package block

import (
	"encoding/binary"
	"lsm/entries"
)

type Iterator struct {
	key       entries.Key
	value     entries.Value
	offsetIdx uint16
	block     Block
}

func (it *Iterator) Key() entries.Key {
	return it.key
}

func (it *Iterator) Value() entries.Value {
	return it.value
}

func (it *Iterator) IsValid() bool {
	return it.key.RawKeySize() > 0
}

func (it *Iterator) Next() error {
	it.offsetIdx++
	it.seekToOffsetIndex(it.offsetIdx)
	return nil
}

func (it *Iterator) seekToOffsetIndex(offsetIdx uint16) {
	if offsetIdx >= uint16(len(it.block.keyValueBeginOffset)) {
		it.key = entries.EmptyKey
		it.value = entries.EmptyValue
		return
	}

	keyValueBeginOffset := it.block.keyValueBeginOffset[offsetIdx]
	it.offsetIdx = offsetIdx
	it.seekToOffset(keyValueBeginOffset)
}

func (it *Iterator) seekToOffset(offset uint16) {
	data := it.block.data[offset:]
	keySize := binary.LittleEndian.Uint16(it.block.data[:])
	key := entries.DecodeFrom(data[ReservedKeySize : uint16(ReservedKeySize)+keySize])

	valueSize := binary.LittleEndian.Uint16(data[ReservedKeySize+key.EncodedKeySizeInBytes():])
	valueOffsetStart := ReservedKeySize + key.EncodedKeySizeInBytes() + ReservedValueSize
	value := entries.NewValue(data[valueOffsetStart : uint16(valueOffsetStart)+valueSize])
	it.key = key
	it.value = value
}

func (it *Iterator) seekToEqualOrGreater(key entries.Key) {
	low := 0
	high := len(it.block.keyValueBeginOffset) - 1

	for low <= high {
		mid := low + (high-low)/2
		it.seekToOffsetIndex(uint16(mid))
		if !it.IsValid() {
			panic("Invalid iterator")
		}
		switch it.key.Compare(key) {
		case 0:
			return
		case -1:
			low = mid + 1
		case 1:
			high = mid - 1
		}
	}
	it.seekToOffset(uint16(low))
}
