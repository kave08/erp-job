package fsyncdata

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/fararavand"
	"erp-job/utility/logger"
	"net/http"

	"go.uber.org/zap"
)

type InvoiceReturn struct {
	log        *zap.SugaredLogger
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	fararavand fararavand.Interface
}

func NewInvoiceReturn(repos *repository.Repository, fr fararavand.Interface) *InvoiceReturn {
	return &InvoiceReturn{
		log:        logger.Logger(),
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (i *InvoiceReturn) InvoiceReturns() error {

	return nil
}
