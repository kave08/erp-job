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
	fmt.Println("step 0", configPath)

	mdb := config.LoadConfig(configPath)

	// fmt.Println("step 0", config.Cfg.BaseURL)
	repos := repository.NewRepository(mdb.MysqlConnection)
	// fmt.Println("step 1 -----", repos)

	ar := logics.NewAryan(repos)
	// fmt.Println("step 2 -----", ar)
	fr := logics.NewFararavand(repos, ar)
	// fmt.Println("step 3 -----", repos, ar)

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
