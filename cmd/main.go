package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	p := "./test"
	fmt.Println(filepath.Base(p))
	fmt.Println(filepath.Dir(p))
}
