package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/matryer/filedb"
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

	dbPath := flag.String("db", "./backup", "the path to db folder")
	flag.Parse()

	dbSession, err := filedb.Dial(*dbPath)
	if err != nil {
		fatalErr = err
		return
	}
	defer dbSession.Close()

	paths, err := dbSession.C("paths")
	if err != nil {
		fatalErr = err
		return
	}

	args := flag.Args()
	subCmd := strings.ToLower(args[0])
	switch subCmd {
	case "list":
		paths.ForEach(func(i int, data []byte) bool {
			var path path
			err := json.Unmarshal(data, &path)
			if err != nil {
				fatalErr = err
				return true
			}

			fmt.Printf("= %s\n", path)
			return false
		})
	case "add":
		newPaths := flag.Args()[1:]
		if len(newPaths) < 1 {
			fatalErr = errors.New("nothing specified, nothing added")
			return
		}
		for _, newPath := range newPaths {
			pathJSON := &path{
				Path: newPath,
				Hash: "not hashed yet",
			}
			if err := paths.InsertJSON(pathJSON); err != nil {
				fatalErr = err
				return
			}

			fmt.Printf("+ %s\n", pathJSON)
		}
	}
}
