package main

import (
	"fmt"
	"log"
)

type path struct {
	Path string
	Hash string
}

func (p path) String() string {
	return fmt.Sprintf("%s [%s]", p.Path, p.Hash)
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatal(fatalErr)
		}
	}()
}
