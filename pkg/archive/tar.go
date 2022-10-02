package archive

import (
	"archive/tar"
	"backup/pkg/filesystem"
	"compress/gzip"
	"github.com/sirupsen/logrus"
	"io"
)

type (
	TarWriter struct {
		Writer   *tar.Writer
		gzWriter *gzip.Writer
	}
)

func (ta *TarWriter) Close() error {
	if err := ta.Writer.Close(); err != nil {
		return err
	}

	if ta.gzWriter != nil {
		if err := ta.gzWriter.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (ta *TarWriter) Write(info *filesystem.FileInfo) error {
	file, err := info.OpenFile()
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	header, err := tar.FileInfoHeader(info, info.RelativePath)
	if err != nil {
		return err
	}

	header.Name = info.RelativePath
	if err = ta.Writer.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(ta.Writer, file)
	if err != nil {
		return err
	}

	return nil
}

func NewTarWriter(writer io.Writer) *TarWriter {
	return &TarWriter{Writer: tar.NewWriter(writer)}
}

func NewTarGzWriter(writer io.Writer) (*TarWriter, error) {
	gz, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	return &TarWriter{Writer: tar.NewWriter(gz), gzWriter: gz}, nil
}
