package archive

import (
	"archive/zip"
	"backup/pkg/directory"
	"io"
)

type (
	ZipArchive struct {
		Writer *zip.Writer
	}
)

func (za *ZipArchive) Close() error {
	return za.Writer.Close()
}

func (za *ZipArchive) Write(info *directory.FileInfo) error {
	file, err := info.File()
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate
	header.Name = info.RelativePath
	archiveFile, err := za.Writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(archiveFile, file)
	if err != nil {
		return err
	}

	return nil
}

func NewZipWriter(writer io.Writer) *ZipArchive {
	return &ZipArchive{Writer: zip.NewWriter(writer)}
}
