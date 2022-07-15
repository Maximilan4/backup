package backup

import (
	"backup/internal/config"
	"backup/internal/drivers"
	"backup/pkg/archive"
	"backup/pkg/directory"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

const (
	defaultArchiveType = archive.TypeTarGz
)

var (
	cmd = &cobra.Command{
		Use:           "backup",
		Short:         "backups single dir",
		Example:       "backup this/dir s3",
		ValidArgs:     []string{"path", "driver"},
		Args:          cobra.MinimumNArgs(2),
		RunE:          backup,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	configPath string
)

func init() {
	cobra.OnInitialize(loadConfig)
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", config.DefaultPath, "--config=config.yaml")
}

func loadConfig() {
	err := config.Load(configPath)
	if err != nil {
		logrus.Fatal(err)
	}
}

func backup(cmd *cobra.Command, args []string) (err error) {
	path, err := directory.Normalize(args[0])
	if err != nil {
		return
	}

	driverInfo := strings.Split(args[1], ".")
	var driverName string
	switch {
	case len(driverInfo) == 0:
		err = errors.New("empty driver name not allowed")
		return
	case len(driverInfo) == 1 || driverInfo[1] == "":
		driverName = "default"
	default:
		driverName = driverInfo[1]
	}

	driver, err := loadDriver(driverInfo[0], driverName)
	if err != nil {
		return
	}

	var archiveType archive.Type
	if len(args) == 2 {
		archiveType = defaultArchiveType
	} else {
		archiveType = archive.Type(args[2])
	}

	return driver.Backup(cmd.Context(), path, archiveType)
}

func loadDriver(driverType, driverName string) (drivers.Driver, error) {
	cfg := config.Get()
	var driver drivers.Driver
	switch driverType {
	case drivers.DirectoryType:
		if dirCfg, ok := cfg.Drivers.Dir[driverName]; ok && dirCfg != nil {
			driver = drivers.NewDirectoryDriver(*dirCfg)
		} else {
			return nil, fmt.Errorf("unable to find configuration %s for %s driver", driverName, driverType)
		}
	case drivers.S3Type:
		if s3Cfg, ok := cfg.Drivers.S3[driverName]; ok && s3Cfg != nil {
			driver = drivers.NewS3Driver(*s3Cfg)
		} else {
			return nil, fmt.Errorf("unable to find configuration %s for %s driver", driverName, driverType)
		}
	default:
		return nil, fmt.Errorf("unsupported driver given: %s", driverType)
	}

	return driver, nil
}

func Run(ctx context.Context) error {
	return cmd.ExecuteContext(ctx)
}
