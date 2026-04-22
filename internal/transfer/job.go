package transfer

import (
	"context"
	"fmt"

	"erp-job/internal/domain"
	"erp-job/internal/store"

	"go.uber.org/zap"
)

const defaultPageSize = 1000

type Source interface {
	FetchInvoices(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Invoices, error)
	FetchCustomers(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Customers, error)
	FetchProducts(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Products, error)
	FetchBaseData(ctx context.Context, pageNumber, pageSize, lastID int) (domain.BaseData, error)
}

type Target interface {
	PostInvoiceToSaleFactor(ctx context.Context, invoices []domain.Invoices) error
	PostProductsToGoods(ctx context.Context, products []domain.Products) error
	PostCustomerToSaleCustomer(ctx context.Context, customers []domain.Customers) error
	PostInvoiceToSaleOrder(ctx context.Context, invoices []domain.Invoices) error
	PostInvoiceToSalePayment(ctx context.Context, invoices []domain.Invoices) error
	PostInvoiceToSaleCenter(ctx context.Context, invoices []domain.Invoices) error
	PostInvoiceToSalerSelect(ctx context.Context, invoices []domain.Invoices) error
	PostInvoiceToSaleProforma(ctx context.Context, invoices []domain.Invoices) error
	PostInvoiceToSaleTypeSelect(ctx context.Context, invoices []domain.Invoices) error
	PostBaseDataToDeliverCenterSaleSelect(ctx context.Context, baseData domain.BaseData) error
}

type Job struct {
	store    store.CheckpointStore
	source   Source
	target   Target
	log      *zap.SugaredLogger
	pageSize int
}

func NewJob(store store.CheckpointStore, source Source, target Target, log *zap.SugaredLogger) *Job {
	return &Job{
		store:    store,
		source:   source,
		target:   target,
		log:      log,
		pageSize: defaultPageSize,
	}
}

func (j *Job) Run(ctx context.Context) error {
	steps := []struct {
		name string
		run  func(context.Context) error
	}{
		{name: "invoice", run: j.runInvoices},
		{name: "base data", run: j.runBaseData},
		{name: "customer", run: j.runCustomers},
		{name: "product", run: j.runProducts},
	}

	for _, step := range steps {
		if j.log != nil {
			j.log.Infow("starting transfer step", "step", step.name)
		}

		if err := step.run(ctx); err != nil {
			return fmt.Errorf("%s sync failed: %w", step.name, err)
		}

		if j.log != nil {
			j.log.Infow("completed transfer step", "step", step.name)
		}
	}

	return nil
}

func nextPage(progress store.Progress, pageSize int) (lastID int, pageNumber int) {
	if progress.LastID == 0 && progress.PageNumber == 0 {
		return 0, 0
	}

	return progress.LastID, progress.PageNumber + pageSize + 1
}

func trimAfterCheckpoint[T any](items []T, lastSyncedID int, idFn func(T) int) []T {
	if len(items) == 0 || idFn(items[len(items)-1]) <= lastSyncedID {
		return nil
	}

	for index, item := range items {
		if idFn(item) > lastSyncedID {
			return items[index:]
		}
	}

	return nil
}

func (j *Job) runInvoices(ctx context.Context) error {
	for {
		progress, err := j.store.GetInvoiceProgress()
		if err != nil {
			return err
		}

		lastID, pageNumber := nextPage(progress, j.pageSize)
		invoices, err := j.source.FetchInvoices(ctx, pageNumber, j.pageSize, lastID)
		if err != nil {
			return err
		}

		if len(invoices) == 0 {
			return nil
		}

		if err := j.syncInvoices(ctx, invoices); err != nil {
			return err
		}

		if err := j.store.SaveInvoiceProgress(store.Progress{
			LastID:     pageNumber + len(invoices),
			PageNumber: pageNumber,
		}); err != nil {
			return err
		}

		if len(invoices) < j.pageSize {
			return nil
		}
	}
}

func (j *Job) syncInvoices(ctx context.Context, invoices []domain.Invoices) error {
	operations := []struct {
		name     string
		getLast  func() (int, error)
		saveLast func(int) error
		post     func(context.Context, []domain.Invoices) error
	}{
		{
			name:     "sale factor",
			getLast:  j.store.GetInvoiceToSaleFactor,
			saveLast: j.store.SaveInvoiceToSaleFactor,
			post:     j.target.PostInvoiceToSaleFactor,
		},
		{
			name:     "sale order",
			getLast:  j.store.GetInvoiceToSaleOrder,
			saveLast: j.store.SaveInvoiceToSaleOrder,
			post:     j.target.PostInvoiceToSaleOrder,
		},
		{
			name:     "sale payment",
			getLast:  j.store.GetInvoiceToSalePayment,
			saveLast: j.store.SaveInvoiceToSalePayment,
			post:     j.target.PostInvoiceToSalePayment,
		},
		{
			name:     "saler select",
			getLast:  j.store.GetInvoiceToSalerSelect,
			saveLast: j.store.SaveInvoiceToSalerSelect,
			post:     j.target.PostInvoiceToSalerSelect,
		},
		{
			name:     "sale proforma",
			getLast:  j.store.GetInvoiceToSaleProforma,
			saveLast: j.store.SaveInvoiceToSaleProforma,
			post:     j.target.PostInvoiceToSaleProforma,
		},
		{
			name:     "sale center",
			getLast:  j.store.GetInvoiceToSaleCenter,
			saveLast: j.store.SaveInvoiceToSaleCenter,
			post:     j.target.PostInvoiceToSaleCenter,
		},
		{
			name:     "sale type select",
			getLast:  j.store.GetInvoiceToSaleTypeSelect,
			saveLast: j.store.SaveInvoiceToSaleTypeSelect,
			post:     j.target.PostInvoiceToSaleTypeSelect,
		},
	}

	for _, op := range operations {
		lastSyncedID, err := op.getLast()
		if err != nil {
			return fmt.Errorf("get %s checkpoint: %w", op.name, err)
		}

		batch := trimAfterCheckpoint(invoices, lastSyncedID, func(item domain.Invoices) int {
			return item.InvoiceId
		})
		if len(batch) == 0 {
			continue
		}

		if err := op.post(ctx, batch); err != nil {
			return fmt.Errorf("post %s: %w", op.name, err)
		}

		if err := op.saveLast(batch[len(batch)-1].InvoiceId); err != nil {
			return fmt.Errorf("save %s checkpoint: %w", op.name, err)
		}
	}

	return nil
}

func (j *Job) runCustomers(ctx context.Context) error {
	for {
		progress, err := j.store.GetCustomerProgress()
		if err != nil {
			return err
		}

		lastID, pageNumber := nextPage(progress, j.pageSize)
		customers, err := j.source.FetchCustomers(ctx, pageNumber, j.pageSize, lastID)
		if err != nil {
			return err
		}

		if len(customers) == 0 {
			return nil
		}

		lastSyncedID, err := j.store.GetCustomerToSaleCustomer()
		if err != nil {
			return err
		}

		batch := trimAfterCheckpoint(customers, lastSyncedID, func(item domain.Customers) int {
			return item.ID
		})
		if len(batch) > 0 {
			if err := j.target.PostCustomerToSaleCustomer(ctx, batch); err != nil {
				return err
			}

			if err := j.store.SaveCustomerToSaleCustomer(batch[len(batch)-1].ID); err != nil {
				return err
			}
		}

		if err := j.store.SaveCustomerProgress(store.Progress{
			LastID:     pageNumber + len(customers),
			PageNumber: pageNumber,
		}); err != nil {
			return err
		}

		if len(customers) < j.pageSize {
			return nil
		}
	}
}

func (j *Job) runProducts(ctx context.Context) error {
	for {
		progress, err := j.store.GetProductProgress()
		if err != nil {
			return err
		}

		lastID, pageNumber := nextPage(progress, j.pageSize)
		products, err := j.source.FetchProducts(ctx, pageNumber, j.pageSize, lastID)
		if err != nil {
			return err
		}

		if len(products) == 0 {
			return nil
		}

		lastSyncedID, err := j.store.GetProductsToGoods()
		if err != nil {
			return err
		}

		batch := trimAfterCheckpoint(products, lastSyncedID, func(item domain.Products) int {
			return item.ID
		})
		if len(batch) > 0 {
			if err := j.target.PostProductsToGoods(ctx, batch); err != nil {
				return err
			}

			if err := j.store.SaveProductsToGoods(batch[len(batch)-1].ID); err != nil {
				return err
			}
		}

		if err := j.store.SaveProductProgress(store.Progress{
			LastID:     pageNumber + len(products),
			PageNumber: pageNumber,
		}); err != nil {
			return err
		}

		if len(products) < j.pageSize {
			return nil
		}
	}
}

func (j *Job) runBaseData(ctx context.Context) error {
	for {
		progress, err := j.store.GetBaseDataProgress()
		if err != nil {
			return err
		}

		lastID, pageNumber := nextPage(progress, j.pageSize)
		baseData, err := j.source.FetchBaseData(ctx, pageNumber, j.pageSize, lastID)
		if err != nil {
			return err
		}

		if len(baseData.PaymentTypes) == 0 {
			return nil
		}

		lastSyncedID, err := j.store.GetBaseDataToDeliverCenter()
		if err != nil {
			return err
		}

		paymentTypes := trimAfterCheckpoint(baseData.PaymentTypes, lastSyncedID, func(item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}) int {
			return item.ID
		})
		if len(paymentTypes) > 0 {
			if err := j.target.PostBaseDataToDeliverCenterSaleSelect(ctx, domain.BaseData{
				PaymentTypes: paymentTypes,
			}); err != nil {
				return err
			}

			if err := j.store.SaveBaseDataToDeliverCenter(paymentTypes[len(paymentTypes)-1].ID); err != nil {
				return err
			}
		}

		if err := j.store.SaveBaseDataProgress(store.Progress{
			LastID:     pageNumber + len(baseData.PaymentTypes),
			PageNumber: pageNumber,
		}); err != nil {
			return err
		}

		if len(baseData.PaymentTypes) < j.pageSize {
			return nil
		}
	}
}
