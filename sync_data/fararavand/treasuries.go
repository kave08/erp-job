package fsyncdata

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/fararavand"
	"erp-job/utility/logger"
	"net/http"

	"go.uber.org/zap"
)

type Treasurie struct {
	log        *zap.SugaredLogger
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	fararavand fararavand.Interface
}

func NewTreasurie(repos *repository.Repository, fr fararavand.Interface) *Treasurie {
	return &Treasurie{
		log:        logger.Logger(),
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (t *Treasurie) Treasuries() error {

	return nil
}
