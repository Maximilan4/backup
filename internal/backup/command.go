package backup

import (
	"backup/internal/config"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:                        "backup",
		Short:                      "backups single dir",
		Example:                    "backup this/dir s3",
		ValidArgs:                  []string{"path", "driver"},
		Args:                       cobra.ExactArgs(2),
		RunE:                       backup,
	}
	configPath string
)

func init() {
	cobra.OnInitialize(loadConfig)
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", config.DefaultPath, "--config=config.yaml")
}

func loadConfig()  {
	err := config.Load(configPath)
	if err != nil {
		logrus.Fatal(err)
	}
}

func backup(cmd *cobra.Command, args []string) error {

	return nil
}


func Run(ctx context.Context) error {
	return cmd.ExecuteContext(ctx)
}