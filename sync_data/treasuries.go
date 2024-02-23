package syncdata

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"erp-job/utility"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Treasurie struct {
	restyClient *resty.Client
	baseURL     string
	httpClient  *http.Client
	repos       *repository.Repository
	aryan       aryan.AryanInterface
	fararavand  fararavand.FararavandInterface
}

func NewTreasurie(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Treasurie {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &Treasurie{
		restyClient: c,
		baseURL:     config.Cfg.FararavandApp.BaseURL,
		repos:       repos,
		aryan:       ar,
		fararavand:  fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (t *Treasurie) Treasuries() error {

	var newTreasuries []models.Treasuries

	resp, err := t.restyClient.R().SetResult(newTreasuries).Get(utility.FGetTreasuries)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	err = t.fararavand.SyncTreasuries(newTreasuries)
	if err != nil {
		fmt.Println("Load SyncTreasuries encountered an error", err.Error())
		return err
	}

	return nil
}
