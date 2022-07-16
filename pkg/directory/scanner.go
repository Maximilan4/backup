package directory

import (
	"io/fs"
	"os"
	"path"
)

type (
	FileScanner struct {
		dir     string
		exclude []string
	}

	FileInfo struct {
		fs.FileInfo
		FullPath     string
		RelativePath string
	}
)

func NewFileScanner(dir string, excludes ...string) *FileScanner {
	return &FileScanner{dir, excludes}
}

func (fi *FileInfo) OpenFile() (*os.File, error) {
	return os.Open(fi.FullPath)
}

func (fs *FileScanner) Scan() ([]*FileInfo, error) {
	return fs.scan(fs.dir, "")
}

func (fs *FileScanner) scan(absPath, relPath string) ([]*FileInfo, error) {
	dirEntries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	files := make([]*FileInfo, 0, len(dirEntries))
	var absolutePath, relativePath string
	for _, dirEntry := range dirEntries {
		absolutePath = path.Join(absPath, dirEntry.Name())
		relativePath = path.Join(relPath, dirEntry.Name())
		if dirEntry.IsDir() {
			nestedFiles, err := fs.scan(absolutePath, relativePath)
			if err != nil {
				return nil, err
			}

			files = append(files, nestedFiles...)
			continue
		}

		info, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		files = append(files, &FileInfo{
			FileInfo:     info,
			FullPath:     absolutePath,
			RelativePath: relativePath,
		})
	}

	return files, nil
}
