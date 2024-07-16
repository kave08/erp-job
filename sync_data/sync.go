package syncdata

import (
	"erp-job/config"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	asyncdata "erp-job/sync_data/aryan"
	fsyncdata "erp-job/sync_data/fararavand"

	"net/http"
)

// Sync orchestrates the synchronization of data between ERP systems.
type Sync struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.Interface
	fararavand fararavand.Interface
}

// NewSync creates a new instance of Sync with the necessary configurations and dependencies.
func NewSync(repos *repository.Repository, fr fararavand.Interface, ar aryan.Interface) *Sync {
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

// Sync synchronizes data between ERP systems by calling various data synchronization functions.
func (s *Sync) Sync() error {
	// fararavnd Sync process
	fsyncdata.NewInvoice(s.repos, s.fararavand)
	fsyncdata.NewBaseData(s.repos, s.fararavand)
	fsyncdata.NewCustomer(s.repos, s.fararavand)
	fsyncdata.NewInvoiceReturn(s.repos, s.fararavand)
	fsyncdata.NewProduct(s.repos, s.fararavand)
	fsyncdata.NewTreasurie(s.repos, s.fararavand)

	// aryan Sync process
	asyncdata.NewLogin()
	asyncdata.NewSaleFactor(s.repos, s.aryan)

	return nil
}
