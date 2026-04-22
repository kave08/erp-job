package cmd

import (
	"erp-job/internal/app"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "apply database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunMigrations(cmd.Context(), configPath)
	},
}

func init() {
	rootCMD.AddCommand(migrateCmd)
}
