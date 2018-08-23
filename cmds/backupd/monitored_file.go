package main

import "fmt"

type monitoredFile struct {
	Path string
	Hash string
}

func (f monitoredFile) String() string {
	return fmt.Sprintf("%s [%s]", f.Path, f.Hash)
}
