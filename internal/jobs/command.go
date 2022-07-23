package jobs

import (
	"backup/internal/config"
	"backup/internal/drivers"
	"backup/pkg/logging"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var (
	ListCmd = &cobra.Command{
		Use:           "jobs",
		Short:         "returns list of jobs from cfg",
		Example:       "returns list of jobs from cfg",
		Args:          cobra.NoArgs,
		RunE:          runTasksCmd,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	ScheduleCmd = &cobra.Command{
		Use:     "schedule",
		Short:   "subcommand for backup background jobs execution",
		Example: "subcommand for backup background jobs execution",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	ScheduleStartCmd = &cobra.Command{
		Use:           "start",
		Short:         "subcommand for backup background jobs execution",
		Example:       "subcommand for backup background jobs execution",
		Args:          cobra.NoArgs,
		RunE:          runScheduleStartCmd,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	ScheduleCmd.AddCommand(ScheduleStartCmd)
}

func runScheduleStartCmd(cmd *cobra.Command, args []string) (err error) {
	cfg := config.Get()
	schedule := cron.New(
		cron.WithLogger(logging.NewCronAdapter()),
		cron.WithLocation(time.Now().Location()),
	)

	var driverInfo drivers.DriverInfo
	var driverCfg any
	var driver drivers.Driver
	var entryID cron.EntryID
	for _, jobCfg := range cfg.Jobs {
		driverInfo, err = jobCfg.Driver()
		if err != nil {
			return
		}

		driverCfg, err = cfg.Drivers.Get(driverInfo)
		if err != nil {
			return
		}

		driver, err = drivers.Load(driverInfo.Type(), driverCfg)
		if err != nil {
			return
		}

		entryID, err = schedule.AddJob(jobCfg.Schedule(), NewJob(jobCfg, driver))
		if err != nil {
			return
		}
		logrus.Infof("Added job %d", entryID)
	}

	schedule.Start()
	<-make(chan struct{})
	return
}

func runTasksCmd(cmd *cobra.Command, args []string) (err error) {
	cfg := config.Get()

	encoder := yaml.NewEncoder(os.Stdout)
	if err = encoder.Encode(cfg.Jobs); err != nil {
		return err
	}

	return
}
