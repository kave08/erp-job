package cmd

import (
	"erp-job/internal/app"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer will fetch data from fararavnd and insert it to aryan",
	RunE: func(cmd *cobra.Command, args []string) error {
		return transfer(cmd)
	},
}

func transfer(cmd *cobra.Command) error {
	return app.RunTransfer(cmd.Context(), configPath)
}

func init() {
	rootCMD.AddCommand(transferCmd)
}
