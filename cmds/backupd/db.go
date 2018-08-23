package main

import (
	"encoding/json"
	"log"

	"github.com/matryer/filedb"
)

type db interface {
	open(url string) error
	close()
	fileList() ([]*monitoredFile, error)
	updateFilesIfUpdated(oldHashs map[string]string) error
}

func newDB() db {
	return newFileDB()
}

type fileDB struct {
	columnName string
	session    *filedb.DB
}

func newFileDB() *fileDB {
	return &fileDB{
		columnName: "monitored_files",
	}
}

func (db *fileDB) open(url string) error {
	var err error
	db.session, err = filedb.Dial(url)
	return err
}

func (db *fileDB) close() {
	db.session.Close()
}

func (db fileDB) fileList() ([]*monitoredFile, error) {
	files := make([]*monitoredFile, 0)
	fileCollection, err := db.session.C(db.columnName)
	if err != nil {
		return files, err
	}

	fileCollection.ForEach(func(_ int, data []byte) bool {
		var file monitoredFile
		if err := json.Unmarshal(data, &file); err != nil {
			log.Printf("filedb could not umarshal json: %s\n", err)
			return true
		}

		files = append(files, &file)
		return false
	})

	return files, nil
}

func (db fileDB) updateFilesIfUpdated(oldHashs map[string]string) error {
	fileCollection, err := db.session.C(db.columnName)
	if err != nil {
		return err
	}

	fileCollection.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
		var file monitoredFile
		if err := json.Unmarshal(data, &file); err != nil {
			log.Printf("faild to unmarshal json: %s\n", err)
			return true, data, false
		}

		file.Hash = oldHashs[file.Path]
		newData, err := json.Marshal(file)
		if err != nil {
			log.Printf("faild to marshal json: %s\n", err)
			return true, data, false
		}

		return true, newData, false
	})

	return nil
}
