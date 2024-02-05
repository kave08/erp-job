package sync

import (
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
)


func Invoice(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) {

	err := fr.SyncInvoicesWithSaleFactor()
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

}
