package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matryer/filedb"
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

	monitor := &backup.Monitor{
		Dest:     *dest,
		Archiver: backup.ZIP,
		Hashs:    make(map[string]string),
	}

	dbSession, err := filedb.Dial(*dbPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer dbSession.Close()
	paths, err := dbSession.C("paths")
	if err != nil {
		log.Println(err)
		return
	}

	paths.ForEach(func(i int, data []byte) bool {
		var path path
		if err := json.Unmarshal(data, &path); err != nil {
			log.Println(err)
			return true
		}

		monitor.Hashs[path.Path] = path.Hash
		return false
	})
	if len(monitor.Hashs) < 1 {
		log.Println("no paths specified. add paths with backup cmd")
		return
	}

	check(monitor, paths)
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT)
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			check(monitor, paths)
		case <-signalCh:
			fmt.Println()
			fmt.Println("stopped backupping")
			return
		}
	}
}

func check(m *backup.Monitor, col *filedb.C) {
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
	col.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
		var path path
		if err := json.Unmarshal(data, &path); err != nil {
			log.Printf("faild to unmarshal json: %s\n", err)
			return true, data, false
		}

		path.Hash = m.Hashs[path.Path]
		newData, err := json.Marshal(path)
		if err != nil {
			log.Printf("faild to marshal json: %s\n", err)
			return true, data, false
		}

		return true, newData, false
	})
}
