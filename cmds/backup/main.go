package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
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

	db := newDB()
	db.open(*dbPath)
	defer db.close()

	args := flag.Args()
	subCmd := strings.ToLower(args[0])
	switch subCmd {
	case "list":
		listMonitoredFiles(db)
	case "add":
		addFileToMonitoredFileList(db)
	case "remove":
		removeMonitoredFiles(db)
	}
}

func listMonitoredFiles(db db) {
	files, err := db.fileList()
	if err != nil {
		log.Printf("could not get file list: %s\n", err)
		return
	}

	for _, file := range files {
		fmt.Printf("- %s", file)
	}
}

func addFileToMonitoredFileList(db db) {
	fileNames := flag.Args()[1:]
	if len(fileNames) < 1 {
		log.Println("nothing specified, nothing added")
		return
	}
	db.addFiles(fileNames...)

	for _, fileName := range fileNames {
		fmt.Printf("+ %s", fileName)
	}
}

func removeMonitoredFiles(db db) {
	fileNames := flag.Args()[1:]
	if len(fileNames) < 1 {
		log.Println("nothing specified, nothing removed")
		return
	}

	db.removeFiles(fileNames...)
}
