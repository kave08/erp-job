package sync

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
	"time"

	"github.com/go-resty/resty/v2"
)

type BaseData struct {
	restyClient *resty.Client
	baseURL     string
	httpClient  *http.Client
	repos       *repository.Repository
	aryan       aryan.AryanInterface
	fararavand  fararavand.FararavandInterface
}

func NewBaseData(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface, requestTimeout time.Duration) *BaseData {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &BaseData{
		restyClient: c,
		baseURL:     config.Cfg.FararavandApp.BaseURL,
		repos:       repos,
		aryan:       ar,
		fararavand:  fr,
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (b *BaseData) BaseData() error {
	var newBaseData models.BaseData

	resp, err := b.restyClient.R().SetResult(newBaseData).Get(utility.FGetBaseData)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	err = b.fararavand.SyncBaseDataWithDeliverCenter(newBaseData)
	if err != nil {
		fmt.Println("Load GetBaseData encountered an error", err.Error())
		return err
	}

	return nil
}
