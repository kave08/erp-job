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
		fmt.Println("Load GetBaseData encountered an error", err.Error())
		return
	}
	_, err = fr.GetProductsToGoods()
	if err != nil {
		fmt.Println("Load GetProductsToGoods encountered an error", err.Error())
		return
	}
	_, err = fr.GetCustomers()
	if err != nil {
		fmt.Println("Load GetCustomers encountered an error", err.Error())
		return
	}
	_, err = fr.GetInvoicesForSaleFactor()
	if err != nil {
		fmt.Println("Load GetInvoicesForSaleFactor encountered an error", err.Error())
		return
	}

	_, err = fr.GetInvoicesForSaleOrder()
	if err != nil {
		fmt.Println("Load GetInvoicesForSaleOrder encountered an error", err.Error())
		return
	}

	_, err = fr.GetInvoicesForSalePayment()
	if err != nil {
		fmt.Println("Load GetInvoicesForSalePayment encountered an error", err.Error())
		return
	}

	_, err = fr.GetInvoicesForSalerSelect()
	if err != nil {
		fmt.Println("Load GetInvoicesForSalerSelect encountered an error", err.Error())
		return
	}

	_, err = fr.GetTreasuries()
	if err != nil {
		fmt.Println("Load GetTreasuries encountered an error", err.Error())
		return
	}
	_, err = fr.GetInvoiceReturns()
	if err != nil {
		fmt.Println("Load GetInvoiceReturns encountered an error", err.Error())
	}
}

func init() {
	rootCMD.AddCommand(transferCmd)
}
