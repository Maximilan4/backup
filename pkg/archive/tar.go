package archive

import (
	"archive/tar"
	"backup/pkg/directory"
	"compress/gzip"
	"io"
)

type (
	TarArchive struct {
		Writer   *tar.Writer
		gzWriter *gzip.Writer
	}
)

func (ta *TarArchive) Close() error {
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

func (ta *TarArchive) Add(info *directory.FileInfo) error {
	file, err := info.File()
	if err != nil {
		return err
	}
	defer file.Close()

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

func NewTarArchive(writer io.Writer) *TarArchive {
	return &TarArchive{Writer: tar.NewWriter(writer)}
}

func NewTarGzArchive(writer io.Writer) (*TarArchive, error) {
	gz, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	return &TarArchive{Writer: tar.NewWriter(gz), gzWriter: gz}, nil
}
