package filesystem

import (
	"errors"
	"os"
	"path"
)

type (
	Scanner interface {
		Scan() chan *FileInfo
	}
)

func NormalizePath(dir string) (newPath string, err error) {
	if len(dir) == 0 {
		err = errors.New("unable to normalize empty path")
		return
	}

	if dir[0] == '~' {
		var homeDir string
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return
		}

		newPath = path.Join(homeDir, dir[1:])
	} else if dir[0] == '.' && (len(dir) == 1 || dir[2] == '/') {
		var curPath string
		curPath, err = os.Getwd()
		if err != nil {
			return
		}

		newPath = path.Join(curPath, dir[1:])
	} else {
		newPath = dir
	}

	return
}
