package archive

import (
	"backup/pkg/filesystem"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

const (
	TypeZip   Type = "zip"
	TypeTar   Type = "tar"
	TypeTarGz Type = "tgz"
)

type (
	Type   string
	Writer interface {
		Write(*filesystem.FileInfo) error
		Close() error
	}
)

func FS(ctx context.Context, writer Writer, scanner filesystem.Scanner) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		entry := logrus.NewEntry(logrus.StandardLogger())
		var err error
		for info := range scanner.Scan() {
			err = writer.Write(info)
			if err != nil {
				return err
			}

			entry.WithFields(logrus.Fields{
				"path":    info.FullPath,
				"size":    info.Size(),
				"mode":    info.Mode(),
				"modTime": info.ModTime(),
			}).Info("file was added to archive")
		}
	}

	return nil
}

func NewWriter(writer io.Writer, archiveType Type) (archive Writer, err error) {
	switch archiveType {
	case TypeZip:
		archive = NewZipWriter(writer)
	case TypeTar:
		archive = NewTarWriter(writer)
	case TypeTarGz:
		archive, err = NewTarGzWriter(writer)
		if err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("unsupported archive type: %s", archiveType)
	}

	return
}

func ExtensionOf(archiveType Type) (extension string, err error) {
	switch archiveType {
	case TypeZip:
		extension = "zip"
	case TypeTar:
		extension = "tar"
	case TypeTarGz:
		extension = "tgz"
	default:
		err = fmt.Errorf("unsupported archive type: %s", archiveType)
	}

	return
}
