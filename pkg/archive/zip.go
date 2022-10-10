package archive

import (
	"archive/zip"
	"backup/pkg/filesystem"
	"io"
)

type (
	ZipWriter struct {
		Writer *zip.Writer
	}
)

func (zw *ZipWriter) Close() error {
	return zw.Writer.Close()
}

func (zw *ZipWriter) Write(info *filesystem.FileInfo) error {
	file, err := info.OpenFile()
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
	archiveFile, err := zw.Writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(archiveFile, file)
	if err != nil {
		return err
	}

	return nil
}

func NewZipWriter(writer io.Writer) *ZipWriter {
	return &ZipWriter{Writer: zip.NewWriter(writer)}
}
