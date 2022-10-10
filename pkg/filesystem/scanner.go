package filesystem

import (
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
)

type (
	FileScanner struct {
		path string
	}

	FileInfo struct {
		fs.FileInfo
		FullPath     string
		RelativePath string
	}
)

func NewFileScanner(dir string) *FileScanner {
	return &FileScanner{dir}
}

func (fi *FileInfo) OpenFile() (*os.File, error) {
	return os.Open(fi.FullPath)
}

func (f *FileScanner) Scan() chan *FileInfo {
	files := make(chan *FileInfo, 100)

	go func() {
		if err := f.scan(files, f.path); err != nil {
			logrus.Fatal(err)
		}
	}()

	return files
}

func (f *FileScanner) scan(files chan *FileInfo, absPath string) error {
	defer close(files)
	return filepath.Walk(absPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		files <- &FileInfo{
			FileInfo:     info,
			FullPath:     path,
			RelativePath: info.Name(),
		}

		return nil
	})
}
