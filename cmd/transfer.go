package cmd

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
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

	ar := aryan.NewAryan(repos)
	fr := fararavand.NewFararavand(repos, ar)

	_, err := fr.GetBaseData()
	if err != nil {
		fmt.Println("Load GetBaseData encountered an error", err.Error())
		return
	}
	err = fr.SyncProductsWithGoods()
	if err != nil {
		fmt.Println("Load SyncProductsWithGoods encountered an error", err.Error())
		return
	}
	err = fr.SyncCustomersWithSaleCustomer()
	if err != nil {
		fmt.Println("Load SyncCustomersWithSaleCustomer encountered an error", err.Error())
		return
	}
	err = fr.SyncInvoicesWithSaleFactor()
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSaleFactor encountered an error", err.Error())
		return
	}

	err = fr.SyncInvoicesWithSaleOrder()
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSaleOrder encountered an error", err.Error())
		return
	}

	err = fr.SyncInvoicesWithSalePayment()
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSalePayment encountered an error", err.Error())
		return
	}

	err = fr.SyncInvoicesWithSalerSelect()
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSalerSelect encountered an error", err.Error())
		return
	}

	err = fr.SyncInvoicesWithSaleProforma()
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSaleProforma encountered an error", err.Error())
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
