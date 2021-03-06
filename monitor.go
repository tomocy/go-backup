package backup

import (
	"fmt"
	"path/filepath"
	"time"
)

type Monitor struct {
	Hashs    map[string]string
	Archiver Archiver
	Dest     string
}

func NewMonitor(archiver Archiver, dest string) *Monitor {
	return &Monitor{
		Hashs:    make(map[string]string),
		Archiver: archiver,
		Dest:     dest,
	}
}

func (m *Monitor) MonitorAndArchive() (int, error) {
	backupCnt := 0
	for path, lastHash := range m.Hashs {
		newHash, err := MD5.Hash(path)
		if err != nil {
			return backupCnt, err
		}
		if newHash == lastHash {
			continue
		}

		if err := m.archive(path); err != nil {
			return backupCnt, err
		}
		m.Hashs[path] = newHash
		backupCnt++
	}

	return backupCnt, nil
}

func (m Monitor) archive(path string) error {
	dirName := filepath.Base(path)
	fileName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), m.Archiver.Extension())
	return m.Archiver.Archive(path, filepath.Join(m.Dest, dirName, fileName))
}
