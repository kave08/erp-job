package asyncdata

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/utility/logger"
	"net/http"

	"go.uber.org/zap"
)

// SaleFactor represents the structure needed to interact with sale factors in the Aryan ERP system.
type SaleFactor struct {
	log        *zap.SugaredLogger
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.Interface
}

// NewSaleFactor initializes a new SaleFactor instance with the necessary configurations and dependencies.
func NewSaleFactor(repos *repository.Repository, ar aryan.Interface) *SaleFactor {
	return &SaleFactor{
		log:     logger.Logger(),
		baseURL: config.Cfg.AryanApp.BaseURL,
		repos:   repos,
		aryan:   ar,
		httpClient: &http.Client{
			Timeout: config.Cfg.AryanApp.Timeout,
		},
	}
}

func (i *SaleFactor) SaleFactors() error {
	return nil
}
