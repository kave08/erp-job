package cmd

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"erp-job/sync"

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

	sync.BaseData(repos, fr, ar)
	sync.SyncCustomer(repos, fr, ar)
	sync.SyncInvoice(repos, fr, ar)
	sync.InvoiceReturns(repos, fr, ar)
	sync.SyncProducts(repos, fr, ar)
	sync.Treasuries(repos, fr, ar)

}

func init() {
	rootCMD.AddCommand(transferCmd)
}
