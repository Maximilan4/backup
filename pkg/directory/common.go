package directory

import (
	"os"
	"path"
)

type (
	Scanner interface {
		Scan() ([]*FileInfo, error)
	}
)

func Normalize(dir string) (newPath string, err error) {
	if dir[0] == '~' {
		var homeDir string
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return
		}

		newPath = path.Join(homeDir, dir[1:])
	} else if dir[0] == '.' {
		var curPath string
		curPath, err = os.Getwd()
		if err != nil {
			return
		}

		newPath = path.Join(curPath, dir[1:])
	}

	return
}
