package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
)

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatal(fatalErr)
		}
	}()

	dbPath := flag.String("db", "./backup", "the path to db folder")
	flag.Parse()
	fmt.Println(*dbPath)

	args := flag.Args()
	if len(args) < 1 {
		fatalErr = errors.New("specify an option")
		return
	}
}
