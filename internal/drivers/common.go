package drivers

import (
	"backup/pkg/archive"
	"context"
	"errors"
	"fmt"
	"strings"
)

type (
	Driver interface {
		Backup(ctx context.Context, dir string, archiveType archive.Type) error
	}

	DriverInfo [2]string
)

func Load(driverType string, cfg any) (Driver, error) {
	var driver Driver
	switch driverType {
	case FS:
		driver = NewDirectoryDriver(*cfg.(*FsDriverConfig))

	case S3:
		driver = NewS3Driver(*cfg.(*S3DriverConfig))

	default:
		return nil, fmt.Errorf("unsupported driver given: %s", driverType)
	}

	return driver, nil
}

func NewDriverInfo(v string) (DriverInfo, error) {
	driverInfo := strings.Split(v, ".")
	var driverName string
	var err error
	switch {
	case len(v) == 0:
		err = errors.New("empty driver name not allowed")
		return DriverInfo{}, err
	case len(driverInfo) == 1 || driverInfo[1] == "":
		driverName = "default"
	default:
		driverName = driverInfo[1]
	}

	return DriverInfo{driverInfo[0], driverName}, nil
}

func (di DriverInfo) Type() string {
	return di[0]
}

func (di DriverInfo) Name() string {
	return di[1]
}
