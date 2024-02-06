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

type InvoiceReturn struct {
	restyClient *resty.Client
	baseURL     string
	httpClient  *http.Client
	repos       *repository.Repository
	aryan       aryan.AryanInterface
	fararavand  fararavand.FararavandInterface
}

func NewInvoiceReturn(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface, requestTimeout time.Duration) *InvoiceReturn {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &InvoiceReturn{
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

func (i *InvoiceReturn) InvoiceReturns() error {
	var newInvoiceReturn []models.InvoiceReturn

	resp, err := i.restyClient.R().SetResult(newInvoiceReturn).Get(utility.FGetInvoiceReturns)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	err = i.fararavand.SyncInvoiceReturns(newInvoiceReturn)
	if err != nil {
		fmt.Println("Load SyncInvoiceReturns encountered an error", err.Error())
	}

	return nil
}
