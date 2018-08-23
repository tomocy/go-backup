package main

import "fmt"

type monitoredFile struct {
	path string
	hash string
}

func (f monitoredFile) String() string {
	return fmt.Sprintf("%s [%s]", f.path, f.hash)
}
