package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct {
	dir string
}

func New(dir string) (*Storage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &Storage{dir: dir}, nil
}

func (s *Storage) Save(docID, filename string, data []byte) error {
	ext := filepath.Ext(filename)
	path := filepath.Join(s.dir, docID+ext)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (s *Storage) Delete(docID, ext string) error {
	path := filepath.Join(s.dir, docID+ext)
	return os.Remove(path)
}
