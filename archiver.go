package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

var ZIP Archiver

func init() {
	ZIP = new(zipper)
}

type Archiver interface {
	Archive(src, dest string) error
	Extension() string
}

type zipper struct {
	writer *zip.Writer
}

func (z zipper) Archive(src, dest string) error {
	destFile, err := createDestinationFile(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	z.writer = zip.NewWriter(destFile)
	defer z.writer.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		if err := z.archive(path); err != nil {
			return err
		}
		return nil
	})
}

func createDestinationFile(dest string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return nil, err
	}

	return os.Create(dest)
}

func (z zipper) archive(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	writer, err := z.writer.Create(src)
	if err != nil {
		return err
	}

	io.Copy(writer, srcFile)
	return nil
}

func (z zipper) Extension() string {
	return "zip"
}
