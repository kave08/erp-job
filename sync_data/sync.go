package syncdata

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"net/http"
)

type Sync struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewSync(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Sync {
	return &Sync{

		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (s *Sync) Sync() error {
	NewInvoice(s.repos, s.fararavand, s.aryan)
	NewBaseData(s.repos, s.fararavand, s.aryan)
	NewCustomer(s.repos, s.fararavand, s.aryan)
	NewInvoiceReturn(s.repos, s.fararavand, s.aryan)
	NewProduct(s.repos, s.fararavand, s.aryan)
	NewTreasurie(s.repos, s.fararavand, s.aryan)

	return nil
}
