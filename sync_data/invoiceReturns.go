package syncdata

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"net/http"
)

type InvoiceReturn struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewInvoiceReturn(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *InvoiceReturn {
	return &InvoiceReturn{
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (i *InvoiceReturn) InvoiceReturns() error {

	return nil
}
