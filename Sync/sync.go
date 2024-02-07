package sync

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Sync struct {
	restyClient *resty.Client
	baseURL     string
	httpClient  *http.Client
	repos       *repository.Repository
	aryan       aryan.AryanInterface
	fararavand  fararavand.FararavandInterface
}

func NewSync(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface, requestTimeout time.Duration) *Sync {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &Sync{
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

func (s *Sync) Sync() error {
	NewBaseData(s.repos, s.fararavand, s.aryan, s.httpClient.Timeout)
	NewCustomer(s.repos, s.fararavand, s.aryan, s.httpClient.Timeout)
	NewInvoice(s.repos, s.fararavand, s.aryan, s.httpClient.Timeout)
	NewInvoiceReturn(s.repos, s.fararavand, s.aryan, s.httpClient.Timeout)
	NewProduct(s.repos, s.fararavand, s.aryan, s.httpClient.Timeout)
	NewTreasurie(s.repos, s.fararavand, s.aryan, s.httpClient.Timeout)

	return nil
}
