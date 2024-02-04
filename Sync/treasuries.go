package sync

import (
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
)

func Treasuries(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) {
	err := fr.SyncTreasuries()
	if err != nil {
		fmt.Println("Load SyncTreasuries encountered an error", err.Error())
		return
	}
}
