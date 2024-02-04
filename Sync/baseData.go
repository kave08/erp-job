package sync

import (
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
)

func BaseData(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) {
	err := fr.SyncBaseData()
	if err != nil {
		fmt.Println("Load GetBaseData encountered an error", err.Error())
		return
	}
}
