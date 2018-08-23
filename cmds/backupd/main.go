package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tomocy/backup"
)

type path struct {
	Path string
	Hash string
}

func main() {
	dest := flag.String("dest", "./archive", "the path to folder backup file is archived")
	dbPath := flag.String("db", "./db", "the path to db")
	interval := flag.Duration("interval", 10, "the interval of backup")
	flag.Parse()

	monitor := backup.NewMonitor(backup.ZIP, *dest)

	db := newDB()
	if err := db.open(*dbPath); err != nil {
		log.Println(err)
		return
	}
	defer db.close()

	if err := setHashsInMonitor(monitor, db); err != nil {
		log.Println(err)
		return
	}
	startMonitoring(*interval, monitor, db)
}

func setHashsInMonitor(monitor *backup.Monitor, db db) error {
	files, err := db.fileList()
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return errors.New("no paths specified. add paths with backup cmd")
	}

	for _, file := range files {
		monitor.Hashs[file.Path] = file.Hash
	}

	return nil
}

func startMonitoring(interval time.Duration, monitor *backup.Monitor, db db) {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	check(monitor, db)
	for {
		select {
		case <-ticker.C:
			check(monitor, db)
		case <-signalCh:
			fmt.Println()
			fmt.Println("stopped backupping")
			return
		}
	}
}

func check(m *backup.Monitor, db db) {
	log.Println("check backup")
	backupCnt, err := m.MonitorAndArchive()
	if err != nil {
		log.Panicf("faild to backup: %s\n", err)
	}
	if backupCnt < 1 {
		log.Println("nothing")
		return
	}

	log.Printf("%d files archived\n", backupCnt)
	db.updateFilesIfUpdated(m.Hashs)
}
