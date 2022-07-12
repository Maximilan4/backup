package drivers

import (
	"backup/pkg/archive"
	"context"
)

const (
	S3Type = "s3"
)

type (
	S3Driver struct {
		cfg S3DriverConfig
	}
	S3DriverConfig struct {
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Bucket    string `mapstructure:"bucket"`
		Url       string `mapstructure:"url"`
		Region    string `mapstructure:"region"`
	}
)

func NewS3Driver(cfg S3DriverConfig) *S3Driver {
	return &S3Driver{cfg: cfg}
}

func (s3d S3Driver) Backup(ctx context.Context, dir string, archiveType archive.Type) error {
	// TODO implement me
	panic("implement me")
}
