package main

import (
	"encoding/json"
	"log"

	"github.com/matryer/filedb"
)

const defaultHash = "not hashed yet"

type db interface {
	open(url string) error
	close()
	getFileList() ([]*monitoredFile, error)
	addFile(file monitoredFile) error
	removeFile(file monitoredFile) error
}

type fileDB struct {
	columnName string
	session    *filedb.DB
}

func (db *fileDB) open(url string) error {
	var err error
	db.session, err = filedb.Dial(url)
	return err
}

func (db *fileDB) close() {
	db.session.Close()
}

func (db fileDB) getFileList() ([]*monitoredFile, error) {
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

func (db fileDB) addFile(fileNames ...string) error {
	fileCollection, err := db.session.C(db.columnName)
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		file := &monitoredFile{
			path: fileName,
			hash: defaultHash,
		}

		if err := fileCollection.InsertJSON(file); err != nil {
			log.Printf("filedb could not insert file into db: %s\n", err)
			continue
		}
	}

	return nil
}

func (db fileDB) removeFiles(fileNames ...string) error {
	fileCollection, err := db.session.C(db.columnName)
	if err != nil {
		return err
	}

	fileCollection.RemoveEach(func(i int, data []byte) (bool, bool) {
		var file monitoredFile
		if err := json.Unmarshal(data, &file); err != nil {
			log.Printf("filedb could not unmarshal json: %s\n", err)
			return false, true
		}

		for _, fileName := range fileNames {
			if fileName == file.path {
				return true, false
			}
		}

		return false, false
	})

	return nil
}
