package syncdata

import (
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
)

type Sync struct {
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.Interface
}

func NewSync(repos *repository.Repository, fr fararavand.Interface, ar aryan.AryanInterface) *Sync {
	return &Sync{
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
	}
}

func (s *Sync) Sync() error {
	steps := []struct {
		name string
		run  func() error
	}{
		{name: "invoice", run: NewInvoice(s.repos, s.fararavand, s.aryan).Invoices},
		{name: "base data", run: NewBaseData(s.repos, s.fararavand, s.aryan).BaseData},
		{name: "customer", run: NewCustomer(s.repos, s.fararavand, s.aryan).Customers},
		{name: "product", run: NewProduct(s.repos, s.fararavand, s.aryan).Products},
	}

	for _, step := range steps {
		if err := step.run(); err != nil {
			return fmt.Errorf("%s sync failed: %w", step.name, err)
		}
	}

	return nil
}
