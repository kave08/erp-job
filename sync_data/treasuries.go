package syncdata

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"net/http"
)

type Treasurie struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewTreasurie(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Treasurie {
	return &Treasurie{
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (t *Treasurie) Treasuries() error {

	return nil
}
