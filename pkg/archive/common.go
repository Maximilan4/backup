package archive

import (
	"backup/pkg/directory"
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
		Write(*directory.FileInfo) error
		Close() error
	}
)

func Directory(ctx context.Context, writer Writer, scanner directory.Scanner) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		infos, err := scanner.Scan()
		if err != nil {
			return err
		}

		entry := logrus.NewEntry(logrus.StandardLogger())
		for _, info := range infos {
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
