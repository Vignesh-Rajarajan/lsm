package wal

import (
	"bufio"
	"fmt"
	"lsm/entries"
	"os"
	"path/filepath"
	"sync"
)

type Wal struct {
	file       *os.File
	buffer     *bufio.Writer
	batchSize  int
	writeCount int
	mu         sync.Mutex
}

func NewWal(id int64, dir string) (*Wal, error) {
	path := filepath.Join(dir, fmt.Sprintf("%d.wal", id))
	_, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Wal{
		file:       file,
		buffer:     bufio.NewWriterSize(file, 4096),
		batchSize:  1000,
		writeCount: 0,
		mu:         sync.Mutex{},
	}, nil
}

func (w *Wal) Close() error {
	if err := w.buffer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}
	return w.file.Close()
}

func (w *Wal) Path() (string, error) {
	return filepath.Abs(w.file.Name())
}

func (w *Wal) Delete() error {
	return os.Remove(w.file.Name())
}

func (w *Wal) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writeCount++
	if w.writeCount >= w.batchSize {
		if err := w.buffer.Flush(); err != nil {
			return fmt.Errorf("failed to flush buffer: %w", err)
		}
		if err := w.file.Sync(); err != nil {
			return fmt.Errorf("failed to sync file: %w", err)
		}
		w.writeCount = 0
	}
	return nil
}

func (w *Wal) Append(key entries.Key, value entries.Value) {
	w.mu.Lock()
	defer w.mu.Unlock()

}
