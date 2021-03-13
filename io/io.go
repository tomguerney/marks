package io

import (
	"io/ioutil"
	"os"
)

type ReaderWriter struct{}

func NewReaderWriter() *ReaderWriter {
	return &ReaderWriter{}
}

func (rw *ReaderWriter) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (rw *ReaderWriter) WriteFile(filename string, data []byte, perm uint32) error {
	return ioutil.WriteFile(filename, data, os.FileMode(perm))
}
