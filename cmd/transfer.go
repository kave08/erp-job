package cmd

import (
	"erp-job/config"
	"erp-job/logics"
	"erp-job/repository"
	"fmt"

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

	ar := logics.NewAryan(repos)
	fr := logics.NewFararavand(repos, ar)

	_, err := fr.GetBaseData()
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	_, err = fr.GetProducts()
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	_, err = fr.GetCustomers()
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	_, err = fr.GetInvoices()
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	_, err = fr.GetTreasuries()
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	_, err = fr.GetInvoiceReturns()
	if err != nil {
		fmt.Println("err", err.Error())
	}
}

func init() {
	rootCMD.AddCommand(transferCmd)
}
