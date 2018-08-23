package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"

	"github.com/matryer/filedb"
	"github.com/tomocy/backup"
)

type path struct {
	Path string
	Hash string
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatalln(fatalErr)
		}
	}()

	dest := flag.String("dest", "./archive", "the path to folder backup file is archived")
	dbPath := flag.String("db", "./db", "the path to db")
	flag.Parse()

	monitor := &backup.Monitor{
		Dest:     *dest,
		Archiver: backup.ZIP,
		Hashs:    make(map[string]string),
	}

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

	paths.ForEach(func(i int, data []byte) bool {
		var path path
		if err := json.Unmarshal(data, &path); err != nil {
			fatalErr = err
			return true
		}

		monitor.Hashs[path.Path] = path.Hash
		return false
	})
	if fatalErr != nil {
		return
	}
	if len(monitor.Hashs) < 1 {
		fatalErr = errors.New("no paths specified. add paths with backup cmd")
		return
	}
}
