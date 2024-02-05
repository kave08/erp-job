package sync

import (
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
)

func Customer(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) {

	err := fr.SyncCustomersWithSaleCustomer()
	if err != nil {
		fmt.Println("Load SyncCustomersWithSaleCustomer encountered an error", err.Error())
		return
	}

}
