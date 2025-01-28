package block

import (
	"github.com/stretchr/testify/assert"
	"lsm/entries"
	"testing"
)

type TestCase struct {
	name           string
	builderEntries []struct {
		key   entries.Key
		value entries.Value
	}
	seekKey        entries.Key
	expectedValid  []bool
	expectedValues []entries.Value
}

func TestBlock_Seek(t *testing.T) {
	tests := []TestCase{
		{
			name: "Seek to first",
			builderEntries: []struct {
				key   entries.Key
				value entries.Value
			}{
				{entries.NewStringKeyWithTimestamp("consensus", 4), entries.NewValue([]byte("raft"))},
				{entries.NewStringKeyWithTimestamp("etcd", 4), entries.NewValue([]byte("kv"))},
			},
			seekKey:       entries.EmptyKey,
			expectedValid: []bool{true, true, false},
			expectedValues: []entries.Value{
				entries.NewStringValue("raft"),
				entries.NewStringValue("kv"),
			},
		},
		{
			name: "Seek to matching key",
			builderEntries: []struct {
				key   entries.Key
				value entries.Value
			}{
				{entries.NewStringKeyWithTimestamp("consensus", 10), entries.NewStringValue("raft")},
				{entries.NewStringKeyWithTimestamp("etcd", 5), entries.NewStringValue("kv")},
			},
			seekKey:       entries.NewStringKeyWithTimestamp("etcd", 5),
			expectedValid: []bool{true, false},
			expectedValues: []entries.Value{
				entries.NewStringValue("kv"),
			},
		},
		{
			name: "Seek to key with less timestamp",
			builderEntries: []struct {
				key   entries.Key
				value entries.Value
			}{
				{entries.NewStringKeyWithTimestamp("consensus", 10), entries.NewStringValue("raft")},
				{entries.NewStringKeyWithTimestamp("etcd", 5), entries.NewStringValue("kv")},
			},
			seekKey:       entries.NewStringKeyWithTimestamp("etcd", 6),
			expectedValid: []bool{true, false},
			expectedValues: []entries.Value{
				entries.NewStringValue("kv"),
			},
		},
		{
			name: "Seek to key matching followed by next",
			builderEntries: []struct {
				key   entries.Key
				value entries.Value
			}{
				{entries.NewStringKeyWithTimestamp("consensus", 5), entries.NewStringValue("raft")},
				{entries.NewStringKeyWithTimestamp("etcd", 5), entries.NewStringValue("kv")},
			},
			seekKey:       entries.NewStringKeyWithTimestamp("consensus", 5),
			expectedValid: []bool{true, true, false},
			expectedValues: []entries.Value{
				entries.NewStringValue("raft"),
				entries.NewStringValue("kv"),
			},
		},
		{
			name: "Seek to key greater than specified followed by next",
			builderEntries: []struct {
				key   entries.Key
				value entries.Value
			}{
				{entries.NewStringKeyWithTimestamp("consensus", 5), entries.NewStringValue("raft")},
				{entries.NewStringKeyWithTimestamp("etcd", 6), entries.NewStringValue("kv")},
				{entries.NewStringKeyWithTimestamp("foundationDB", 7), entries.NewStringValue("distributed-kv")},
			},
			seekKey:       entries.NewStringKeyWithTimestamp("distributed", 8),
			expectedValid: []bool{true, true, false},
			expectedValues: []entries.Value{
				entries.NewStringValue("kv"),
				entries.NewStringValue("distributed-kv"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockBuilder := NewBuilder(4096)
			for _, entry := range tt.builderEntries {
				blockBuilder.Add(entry.key, entry.value)
			}
			block := blockBuilder.Build()

			var iterator *Iterator
			if tt.seekKey.EqualTo(entries.EmptyKey) {
				iterator = block.SeekToFirst()
			} else {
				iterator = block.Seek(tt.seekKey)
			}
			defer iterator.Close()

			for i, expectedValid := range tt.expectedValid {
				assert.Equal(t, expectedValid, iterator.IsValid())
				if expectedValid {
					assert.Equal(t, tt.expectedValues[i], iterator.Value())
					iterator.Next()
				}
			}
		})
	}
}
