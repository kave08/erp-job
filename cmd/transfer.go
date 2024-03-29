package cmd

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	syncdata "erp-job/sync_data"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer will fetch data from fararavnd and insert it to aryan",
	Run: func(cmd *cobra.Command, args []string) {
		transfer()
	},
}

func transfer() {

	mdb := config.LoadConfig(configPath)

	repos := repository.NewRepository(mdb.MysqlConnection)

	ar := aryan.NewAryan(repos)
	fr := fararavand.NewFararavand(repos, ar)

	syncdata.NewSync(repos, fr, ar)

}

func init() {
	rootCMD.AddCommand(transferCmd)
}
