package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	dbPath := flag.String("db", "./backup", "the path to db folder")
	flag.Parse()

	db := newDB()
	if err := db.open(*dbPath); err != nil {
		log.Println(err)
		return
	}
	defer db.close()

	subCmd := strings.ToLower(flag.Arg(0))
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
	if len(files) < 1 {
		fmt.Println("no files to be monitored")
		return
	}

	for _, file := range files {
		fmt.Printf("= %s\n", file)
	}
}

func addFileToMonitoredFileList(db db) {
	fileNames := flag.Args()[1:]
	if len(fileNames) < 1 {
		log.Println("nothing specified, nothing added")
		return
	}

	db.addFiles(fileNames...)
}

func removeMonitoredFiles(db db) {
	fileNames := flag.Args()[1:]
	if len(fileNames) < 1 {
		log.Println("nothing specified, nothing removed")
		return
	}

	db.removeFiles(fileNames...)
}
