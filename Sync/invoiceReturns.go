package sync

import (
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
)

func InvoiceReturns(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) {
	err := fr.SyncInvoiceReturns()
	if err != nil {
		fmt.Println("Load SyncInvoiceReturns encountered an error", err.Error())
	}
}
