package backup

import (
	md5lib "crypto/md5"
	"fmt"
	"hash"
	"os"
	"path/filepath"
)

var MD5 Hasher

func init() {
	MD5 = &md5{
		hasher: md5lib.New(),
	}
}

type Hasher interface {
	HashDir(path string) (string, error)
}

type md5 struct {
	hasher hash.Hash
}

func (m md5) HashDir(path string) (string, error) {
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		m.writeFileInfo(path, info)
		return nil
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", m.hasher.Sum(nil)), nil
}

func (m *md5) writeFileInfo(path string, fileInfo os.FileInfo) {
	m.hasher.Reset()
	fmt.Fprintf(m.hasher, "%v", path)
	fmt.Fprintf(m.hasher, "%v", fileInfo.Name())
	fmt.Fprintf(m.hasher, "%v", fileInfo.Size())
	fmt.Fprintf(m.hasher, "%v", fileInfo.Mode())
	fmt.Fprintf(m.hasher, "%v", fileInfo.ModTime())
	fmt.Fprintf(m.hasher, "%v", fileInfo.IsDir())
}
