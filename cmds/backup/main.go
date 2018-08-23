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
			var file monitoredFile
			if err := json.Unmarshal(data, &file); err != nil {
				fatalErr = err
				return true
			}

			fmt.Printf("= %s\n", file)
			return false
		})
	case "add":
		newPaths := flag.Args()[1:]
		if len(newPaths) < 1 {
			fatalErr = errors.New("nothing specified, nothing added")
			return
		}
		for _, newPath := range newPaths {
			pathJSON := &monitoredFile{
				path: newPath,
				hash: "not hashed yet",
			}
			if err := paths.InsertJSON(pathJSON); err != nil {
				fatalErr = err
				return
			}

			fmt.Printf("+ %s\n", pathJSON)
		}
	case "remove":
		paths.RemoveEach(func(i int, data []byte) (bool, bool) {
			var file monitoredFile
			if err := json.Unmarshal(data, &file); err != nil {
				fatalErr = err
				return false, true
			}

			pathsToBeRemoved := flag.Args()[1:]
			for _, pathToBeRemoved := range pathsToBeRemoved {
				if pathToBeRemoved == file.path {
					return true, false
				}
			}

			return false, false
		})
	}
}
