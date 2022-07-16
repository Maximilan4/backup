package tasks

import (
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:           "tasks",
		Short:         "tasks list for backup",
		Example:       "tasks list for backup",
		ValidArgs:     []string{"path", "driver"},
		Args:          cobra.MinimumNArgs(2),
		RunE:          runTasksCmd,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func runTasksCmd(cmd *cobra.Command, args []string) error {
	return nil
}
