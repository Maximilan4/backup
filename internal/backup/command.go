package backup

import (
	"backup/internal/config"
	"backup/internal/drivers"
	"backup/internal/jobs"
	"backup/pkg/archive"
	"backup/pkg/filesystem"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultArchiveType = archive.TypeTarGz
)

var (
	cmd = &cobra.Command{
		Use:           "backup",
		Short:         "backups single dir",
		Example:       "backup this/dir s3 tgz",
		ValidArgs:     []string{"path", "driver", "archive-type"},
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

	cmd.AddCommand(jobs.ListCmd, jobs.ScheduleCmd)
}

func loadConfig() {
	err := config.Load(configPath)
	if err != nil {
		logrus.Fatal(err)
	}
}

func backup(cmd *cobra.Command, args []string) (err error) {
	cfg := config.Get()
	path, err := filesystem.NormalizePath(args[0])
	if err != nil {
		return
	}

	driverInfo, err := drivers.NewDriverInfo(args[1])
	if err != nil {
		return
	}

	driverCfg, err := cfg.Drivers.Get(driverInfo)
	if err != nil {
		return
	}

	driver, err := drivers.Load(driverInfo.Type(), driverCfg)
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

func Run(ctx context.Context) error {
	return cmd.ExecuteContext(ctx)
}
