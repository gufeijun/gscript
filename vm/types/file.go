package types

import "os"

type File struct {
	File *os.File
}

func NewFile(file *os.File) *File {
	return &File{File: file}
}
