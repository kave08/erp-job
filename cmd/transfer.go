package cmd

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	syncdata "erp-job/sync_data"
	"erp-job/utility/logger"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer will fetch data from fararavnd and insert it to aryan",
	RunE: func(cmd *cobra.Command, args []string) error {
		return transfer()
	},
}

func transfer() error {
	setup, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}
	defer setup.MysqlConnection.Close()

	logger.Initialize()

	repos := repository.NewRepository(setup.MysqlConnection)

	ar := aryan.NewAryan(repos)
	fr := fararavand.NewFararavand(repos, ar)

	return syncdata.NewSync(repos, fr, ar).Sync()
}

func init() {
	rootCMD.AddCommand(transferCmd)
}
