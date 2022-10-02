package drivers

import (
	"backup/pkg/archive"
	"backup/pkg/directory"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"path"
	"time"
)

const (
	S3Type = "s3"
)

type (
	S3Driver struct {
		cfg S3DriverConfig
	}
	S3DriverConfig struct {
		Profile string `mapstructure:"profile"`
		Bucket  string `mapstructure:"bucket"`
		Url     string `mapstructure:"url"`
		Path    string `mapstructure:"path"`
	}
)

func NewS3Driver(cfg S3DriverConfig) *S3Driver {
	return &S3Driver{cfg: cfg}
}

func (s3d *S3Driver) ConfigureClient(ctx context.Context) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               s3d.cfg.Url,
			SigningRegion:     region,
			Source:            aws.EndpointSourceCustom,
			HostnameImmutable: false,
		}, nil
	})

	optionsFuncs := make([]func(options *config.LoadOptions) error, 0, 2)
	optionsFuncs = append(optionsFuncs, config.WithSharedConfigProfile(s3d.cfg.Profile))
	if s3d.cfg.Url != "" {
		optionsFuncs = append(optionsFuncs, config.WithEndpointResolverWithOptions(customResolver))
	}

	cfg, err := config.LoadDefaultConfig(ctx, optionsFuncs...)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func (s3d *S3Driver) ConfigureUploader(client *s3.Client) *manager.Uploader {
	return manager.NewUploader(client, func(uploader *manager.Uploader) {
		uploader.LeavePartsOnError = false
		uploader.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(32 * 1024 * 1024)
	})
}

func (s3d *S3Driver) GetFilename(archiveType archive.Type) string {
	return path.Join(s3d.cfg.Path, fmt.Sprintf("%d.%s", time.Now().Unix(), archiveType))
}

func (s3d *S3Driver) Backup(ctx context.Context, dir string, archiveType archive.Type) error {
	client, err := s3d.ConfigureClient(ctx)
	if err != nil {
		return err
	}

	uploader := s3d.ConfigureUploader(client)
	reader, writer := io.Pipe()
	defer func() {
		if err = reader.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer func() {
			if err = writer.Close(); err != nil {
				logrus.Fatal(err)
			}
		}()

		archiveWriter, err := archive.NewWriter(writer, archiveType)
		if err != nil {
			return err
		}

		defer func() {
			if err = archiveWriter.Close(); err != nil {
				logrus.Fatal(err)
			}
		}()

		return archive.Directory(gCtx, archiveWriter, directory.NewFileScanner(dir))
	})

	name := s3d.GetFilename(archiveType)
	logrus.Infof("Created a new file: %s", name)
	logrus.WithField("bucket", s3d.cfg.Bucket).Infof("Uploading file: %s", name)
	output, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3d.cfg.Bucket),
		Key:    aws.String(name),
		Body:   reader,
	})

	if err != nil {
		return err
	}

	logrus.WithField("bucket", s3d.cfg.Bucket).Infof("File %s uploaded: %s", name, output.UploadID)

	return nil
}
