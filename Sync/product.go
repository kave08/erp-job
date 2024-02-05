package sync

import (
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
)

func Products(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) {
	err := fr.SyncProductsWithGoods()
	if err != nil {
		fmt.Println("Load SyncProductsWithGoods encountered an error", err.Error())
		return
	}
}
