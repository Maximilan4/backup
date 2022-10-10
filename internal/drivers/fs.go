package drivers

import (
	"backup/pkg/archive"
	"backup/pkg/filesystem"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

const (
	FS = "fs"
)

type (
	FileSystemDriver struct {
		cfg FsDriverConfig
	}
	FsDriverConfig struct {
		OutputPath string `mapstructure:"output_path"`
	}
)

func NewDirectoryDriver(cfg FsDriverConfig) *FileSystemDriver {
	return &FileSystemDriver{cfg: cfg}
}

func (dd *FileSystemDriver) createOutputDir() error {
	return os.MkdirAll(dd.cfg.OutputPath, 0700)
}

func (dd *FileSystemDriver) createFile(archiveType archive.Type) (*os.File, error) {
	if err := dd.createOutputDir(); err != nil {
		return nil, err
	}

	extension, err := archive.ExtensionOf(archiveType)
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%d.%s", time.Now().Unix(), extension)

	return os.Create(path.Join(dd.cfg.OutputPath, filename))
}

func (dd *FileSystemDriver) Backup(ctx context.Context, dir string, archiveType archive.Type) error {
	file, err := dd.createFile(archiveType)
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	archiveWriter, err := archive.NewWriter(file, archiveType)
	if err != nil {
		return err
	}

	defer func() {
		if err = archiveWriter.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	return archive.FS(ctx, archiveWriter, filesystem.NewFileScanner(dir))
}
