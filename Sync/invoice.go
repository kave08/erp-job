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

type Invoice struct {
	restyClient *resty.Client
	baseURL     string
	httpClient  *http.Client
	repos       *repository.Repository
	aryan       aryan.AryanInterface
	fararavand  fararavand.FararavandInterface
}

func NewInvoice(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface, requestTimeout time.Duration) *Invoice {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &Invoice{
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

func (i *Invoice) Invoices() error {
	var newInvoices []models.Invoices

	resp, err := i.restyClient.R().SetResult(newInvoices).Get(utility.FGetInvoices)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	err = i.fararavand.SyncInvoicesWithSaleFactor(newInvoices)
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSaleFactor encountered an error", err.Error())
		return err
	}

	err = i.fararavand.SyncInvoicesWithSaleOrder(newInvoices)
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSaleOrder encountered an error", err.Error())
		return err
	}

	err = i.fararavand.SyncInvoicesWithSalePayment(newInvoices)
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSalePayment encountered an error", err.Error())
		return err
	}

	err = i.fararavand.SyncInvoicesWithSalerSelect(newInvoices)
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSalerSelect encountered an error", err.Error())
		return err
	}

	err = i.fararavand.SyncInvoicesWithSaleProforma(newInvoices)
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSaleProforma encountered an error", err.Error())
		return err
	}
	err = i.fararavand.SyncInvoicesWithSaleCenter(newInvoices)
	if err != nil {
		fmt.Println("Load SyncInvoicesWithSaleCenter encountered an error", err.Error())
		return err
	}

	return nil
}
