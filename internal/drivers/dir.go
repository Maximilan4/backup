package drivers

import (
	"backup/pkg/archive"
	"backup/pkg/directory"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

const (
	DirectoryType = "dir"
)

type (
	DirectoryDriver struct {
		cfg DirectoryDriverConfig
	}
	DirectoryDriverConfig struct {
		OutputPath string `mapstructure:"output_path"`
	}
)

func NewDirectoryDriver(cfg DirectoryDriverConfig) *DirectoryDriver {
	return &DirectoryDriver{cfg: cfg}
}

func (dd *DirectoryDriver) createOutputDir() error {
	return os.MkdirAll(dd.cfg.OutputPath, 0700)
}

func (dd *DirectoryDriver) createFile(archiveType archive.Type) (*os.File, error) {
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

func (dd *DirectoryDriver) Backup(ctx context.Context, dir string, archiveType archive.Type) error {
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

	return archive.Directory(ctx, archiveWriter, directory.NewFileScanner(dir))
}
