package archive

import (
	"backup/pkg/directory"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"sync"
)

const (
	TypeZip   Type = "zip"
	TypeTar   Type = "tar"
	TypeTarGz Type = "tgz"
)

type (
	Type    string
	Archive interface {
		Add(*directory.FileInfo) error
		Close() error
	}
	Archiver struct {
		Archive Archive
		scanner directory.Scanner
		mut     *sync.Mutex
	}
)

func (za *Archiver) Create(ctx context.Context) (err error) {
	files, err := za.scanner.Scan()
	if err != nil {
		return err
	}

	err = za.addFiles(ctx, files)

	return
}

func (za *Archiver) addFiles(ctx context.Context, files []*directory.FileInfo) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, file := range files {
		fileInfo := file
		group.Go(func() error {
			return za.addFile(groupCtx, fileInfo)
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func (za *Archiver) addFile(ctx context.Context, info *directory.FileInfo) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		za.mut.Lock()
		defer za.mut.Unlock()

		logrus.Infof("archiving file: %s", info.FullPath)
		err := za.Archive.Add(info)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewArchiver(archive Archive, scanner directory.Scanner) *Archiver {
	return &Archiver{
		Archive: archive,
		scanner: scanner,
		mut:     &sync.Mutex{},
	}
}

func NewArchive(writer io.Writer, archiveType Type) (archive Archive, err error) {
	switch archiveType {
	case TypeZip:
		archive = NewZipArchive(writer)
	case TypeTar:
		archive = NewTarArchive(writer)
	case TypeTarGz:
		archive, err = NewTarGzArchive(writer)
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
