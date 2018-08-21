package backup

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
)

func hashDir(path string) (string, error) {
	hasher := md5.New()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Fprintf(hasher, "%v", path)
		fmt.Fprintf(hasher, "%v", info.Name())
		fmt.Fprintf(hasher, "%v", info.Size())
		fmt.Fprintf(hasher, "%v", info.Mode())
		fmt.Fprintf(hasher, "%v", info.ModTime())
		fmt.Fprintf(hasher, "%v", info.IsDir())
		return nil
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
