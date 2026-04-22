package transfer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"erp-job/internal/domain"
	"erp-job/internal/observability"
	"erp-job/internal/store"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

const (
	defaultPageSize  = 1000
	cursorPageNumber = 0
)

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
	store     store.CheckpointStore
	source    Source
	target    Target
	log       *zap.SugaredLogger
	pageSize  int
	telemetry *observability.Telemetry
}

type invoiceOperation struct {
	name      string
	operation store.Operation
	post      func(context.Context, []domain.Invoices) error
	dedupeKey func(domain.Invoices) string
}

func NewJob(store store.CheckpointStore, source Source, target Target, log *zap.SugaredLogger, telemetry *observability.Telemetry) *Job {
	return &Job{
		store:     store,
		source:    source,
		target:    target,
		log:       log,
		pageSize:  defaultPageSize,
		telemetry: telemetry,
	}
}

func (j *Job) Run(ctx context.Context) (runErr error) {
	startedAt := time.Now()
	runID := observability.RunIDFromContext(ctx)
	ctx, span := j.telemetry.StartRun(ctx, runID)
	defer func() {
		result := "success"
		if runErr != nil {
			result = "failed"
			span.RecordError(runErr)
			span.SetStatus(codes.Error, runErr.Error())
			j.telemetry.RecordFailure(ctx, "job", "run", 0, "job_run_failed")
		}
		j.telemetry.RecordRun(ctx, result, time.Since(startedAt))
		span.End()
	}()

	steps := []struct {
		name   string
		entity store.Entity
		run    func(context.Context) error
	}{
		{name: "invoice", entity: store.EntityInvoice, run: j.runInvoices},
		{name: "base_data", entity: store.EntityBaseData, run: j.runBaseData},
		{name: "customer", entity: store.EntityCustomer, run: j.runCustomers},
		{name: "product", entity: store.EntityProduct, run: j.runProducts},
	}

	for _, step := range steps {
		stepCtx, stepSpan := j.telemetry.StartStep(ctx, step.name)
		stepStartedAt := time.Now()

		if j.log != nil {
			j.log.Infow("starting transfer step", "run_id", runID, "step", step.name)
		}

		if err := step.run(stepCtx); err != nil {
			stepSpan.RecordError(err)
			stepSpan.SetStatus(codes.Error, err.Error())
			j.telemetry.RecordFailure(stepCtx, "job", step.name, 0, "step_failed")
			stepSpan.End()
			return fmt.Errorf("%s sync failed: %w", step.name, err)
		}

		stepSpan.SetAttributes(attribute.Int64("duration_ms", time.Since(stepStartedAt).Milliseconds()))
		stepSpan.End()

		if j.log != nil {
			j.log.Infow("completed transfer step", "run_id", runID, "step", step.name, "duration_ms", time.Since(stepStartedAt).Milliseconds())
		}
	}

	return nil
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
		lastSourceID, err := j.store.GetSourceProgress(ctx, store.EntityInvoice)
		if err != nil {
			return err
		}

		invoices, err := j.source.FetchInvoices(ctx, cursorPageNumber, j.pageSize, lastSourceID)
		if err != nil {
			return err
		}

		if len(invoices) == 0 {
			return nil
		}

		batchLastID := invoices[len(invoices)-1].InvoiceId
		if batchLastID <= lastSourceID {
			return fmt.Errorf("invoice cursor did not advance: last_source_id=%d batch_last_id=%d", lastSourceID, batchLastID)
		}

		j.telemetry.RecordFetched(ctx, string(store.EntityInvoice), len(invoices))
		j.telemetry.RecordCheckpointLag(ctx, string(store.EntityInvoice), batchLastID-lastSourceID)

		if err := j.syncInvoices(ctx, invoices); err != nil {
			return err
		}

		if err := j.store.AdvanceSourceProgress(ctx, store.EntityInvoice, batchLastID); err != nil {
			return fmt.Errorf("advance invoice source progress: %w", err)
		}

		if len(invoices) < j.pageSize {
			return nil
		}
	}
}

func (j *Job) syncInvoices(ctx context.Context, invoices []domain.Invoices) error {
	operations := []invoiceOperation{
		{
			name:      "sale factor",
			operation: store.OperationInvoiceSaleFactor,
			post:      j.target.PostInvoiceToSaleFactor,
			dedupeKey: func(item domain.Invoices) string {
				return fmt.Sprintf("invoice:%d:product:%d", item.InvoiceId, item.ProductID)
			},
		},
		{
			name:      "sale order",
			operation: store.OperationInvoiceSaleOrder,
			post:      j.target.PostInvoiceToSaleOrder,
			dedupeKey: func(item domain.Invoices) string {
				return fmt.Sprintf("invoice:%d:product:%d", item.InvoiceId, item.ProductID)
			},
		},
		{
			name:      "sale payment",
			operation: store.OperationInvoiceSalePayment,
			post:      j.target.PostInvoiceToSalePayment,
			dedupeKey: func(item domain.Invoices) string { return fmt.Sprintf("payment:%d", item.PaymentTypeID) },
		},
		{
			name:      "saler select",
			operation: store.OperationInvoiceSalerSelect,
			post:      j.target.PostInvoiceToSalerSelect,
			dedupeKey: func(item domain.Invoices) string { return fmt.Sprintf("visitor:%s", item.VisitorCode) },
		},
		{
			name:      "sale proforma",
			operation: store.OperationInvoiceSaleProforma,
			post:      j.target.PostInvoiceToSaleProforma,
			dedupeKey: func(item domain.Invoices) string {
				return fmt.Sprintf("invoice:%d:product:%d", item.InvoiceId, item.ProductID)
			},
		},
		{
			name:      "sale center",
			operation: store.OperationInvoiceSaleCenter,
			post:      j.target.PostInvoiceToSaleCenter,
			dedupeKey: func(item domain.Invoices) string { return fmt.Sprintf("stock:%d", item.WareHouseID) },
		},
		{
			name:      "sale type select",
			operation: store.OperationInvoiceSaleTypeSelect,
			post:      j.target.PostInvoiceToSaleTypeSelect,
			dedupeKey: func(item domain.Invoices) string { return fmt.Sprintf("sale_type:%d", item.SNoePardakht) },
		},
	}

	for _, op := range operations {
		lastSyncedID, err := j.store.GetOperationCheckpoint(ctx, op.operation)
		if err != nil {
			return fmt.Errorf("get %s checkpoint: %w", op.name, err)
		}

		batch := trimAfterCheckpoint(invoices, lastSyncedID, func(item domain.Invoices) int {
			return item.InvoiceId
		})
		if len(batch) == 0 {
			continue
		}

		j.telemetry.RecordCheckpointLag(ctx, string(store.EntityInvoice), batch[len(batch)-1].InvoiceId-lastSyncedID)
		if err := j.deliverInvoiceOperation(ctx, op, batch); err != nil {
			return err
		}
	}

	return nil
}

func (j *Job) deliverInvoiceOperation(ctx context.Context, op invoiceOperation, batch []domain.Invoices) error {
	observedCtx := withDeliveryObserver(ctx, j.store, j.log, op.operation, batch,
		func(item domain.Invoices) int { return item.InvoiceId },
		op.dedupeKey,
	)

	if err := op.post(observedCtx, batch); err != nil {
		return fmt.Errorf("post %s: %w", op.name, err)
	}

	if err := j.store.MarkBatchDelivered(ctx, op.operation, batch[len(batch)-1].InvoiceId, buildDeliveredRecords(batch,
		func(item domain.Invoices) int { return item.InvoiceId },
		op.dedupeKey,
	)); err != nil {
		return fmt.Errorf("save %s checkpoint: %w", op.name, err)
	}

	j.telemetry.RecordPosted(ctx, string(op.operation), len(batch))

	if j.log != nil {
		j.log.Infow("completed target batch",
			"run_id", observability.RunIDFromContext(ctx),
			"step", "invoice",
			"operation", string(op.operation),
			"batch_size", len(batch),
			"last_source_id", batch[len(batch)-1].InvoiceId,
			"status_code", http.StatusOK,
		)
	}

	return nil
}

func (j *Job) runCustomers(ctx context.Context) error {
	for {
		lastSourceID, err := j.store.GetSourceProgress(ctx, store.EntityCustomer)
		if err != nil {
			return err
		}

		customers, err := j.source.FetchCustomers(ctx, cursorPageNumber, j.pageSize, lastSourceID)
		if err != nil {
			return err
		}

		if len(customers) == 0 {
			return nil
		}

		batchLastID := customers[len(customers)-1].ID
		if batchLastID <= lastSourceID {
			return fmt.Errorf("customer cursor did not advance: last_source_id=%d batch_last_id=%d", lastSourceID, batchLastID)
		}

		j.telemetry.RecordFetched(ctx, string(store.EntityCustomer), len(customers))
		j.telemetry.RecordCheckpointLag(ctx, string(store.EntityCustomer), batchLastID-lastSourceID)

		lastSyncedID, err := j.store.GetOperationCheckpoint(ctx, store.OperationCustomerSaleCustomer)
		if err != nil {
			return err
		}

		batch := trimAfterCheckpoint(customers, lastSyncedID, func(item domain.Customers) int {
			return item.ID
		})
		if len(batch) > 0 {
			observedCtx := withDeliveryObserver(ctx, j.store, j.log, store.OperationCustomerSaleCustomer, batch,
				func(item domain.Customers) int { return item.ID },
				func(item domain.Customers) string { return fmt.Sprintf("customer:%d", item.ID) },
			)
			if err := j.target.PostCustomerToSaleCustomer(observedCtx, batch); err != nil {
				return err
			}

			if err := j.store.MarkBatchDelivered(ctx, store.OperationCustomerSaleCustomer, batch[len(batch)-1].ID, buildDeliveredRecords(batch,
				func(item domain.Customers) int { return item.ID },
				func(item domain.Customers) string { return fmt.Sprintf("customer:%d", item.ID) },
			)); err != nil {
				return err
			}

			j.telemetry.RecordPosted(ctx, string(store.OperationCustomerSaleCustomer), len(batch))
		}

		if err := j.store.AdvanceSourceProgress(ctx, store.EntityCustomer, batchLastID); err != nil {
			return fmt.Errorf("advance customer source progress: %w", err)
		}

		if len(customers) < j.pageSize {
			return nil
		}
	}
}

func (j *Job) runProducts(ctx context.Context) error {
	for {
		lastSourceID, err := j.store.GetSourceProgress(ctx, store.EntityProduct)
		if err != nil {
			return err
		}

		products, err := j.source.FetchProducts(ctx, cursorPageNumber, j.pageSize, lastSourceID)
		if err != nil {
			return err
		}

		if len(products) == 0 {
			return nil
		}

		batchLastID := products[len(products)-1].ID
		if batchLastID <= lastSourceID {
			return fmt.Errorf("product cursor did not advance: last_source_id=%d batch_last_id=%d", lastSourceID, batchLastID)
		}

		j.telemetry.RecordFetched(ctx, string(store.EntityProduct), len(products))
		j.telemetry.RecordCheckpointLag(ctx, string(store.EntityProduct), batchLastID-lastSourceID)

		lastSyncedID, err := j.store.GetOperationCheckpoint(ctx, store.OperationProductsGoods)
		if err != nil {
			return err
		}

		batch := trimAfterCheckpoint(products, lastSyncedID, func(item domain.Products) int {
			return item.ID
		})
		if len(batch) > 0 {
			observedCtx := withDeliveryObserver(ctx, j.store, j.log, store.OperationProductsGoods, batch,
				func(item domain.Products) int { return item.ID },
				func(item domain.Products) string { return fmt.Sprintf("product:%d", item.ID) },
			)
			if err := j.target.PostProductsToGoods(observedCtx, batch); err != nil {
				return err
			}

			if err := j.store.MarkBatchDelivered(ctx, store.OperationProductsGoods, batch[len(batch)-1].ID, buildDeliveredRecords(batch,
				func(item domain.Products) int { return item.ID },
				func(item domain.Products) string { return fmt.Sprintf("product:%d", item.ID) },
			)); err != nil {
				return err
			}

			j.telemetry.RecordPosted(ctx, string(store.OperationProductsGoods), len(batch))
		}

		if err := j.store.AdvanceSourceProgress(ctx, store.EntityProduct, batchLastID); err != nil {
			return fmt.Errorf("advance product source progress: %w", err)
		}

		if len(products) < j.pageSize {
			return nil
		}
	}
}

func (j *Job) runBaseData(ctx context.Context) error {
	for {
		lastSourceID, err := j.store.GetSourceProgress(ctx, store.EntityBaseData)
		if err != nil {
			return err
		}

		baseData, err := j.source.FetchBaseData(ctx, cursorPageNumber, j.pageSize, lastSourceID)
		if err != nil {
			return err
		}

		if len(baseData.PaymentTypes) == 0 {
			return nil
		}

		batchLastID := baseData.PaymentTypes[len(baseData.PaymentTypes)-1].ID
		if batchLastID <= lastSourceID {
			return fmt.Errorf("base-data cursor did not advance: last_source_id=%d batch_last_id=%d", lastSourceID, batchLastID)
		}

		j.telemetry.RecordFetched(ctx, string(store.EntityBaseData), len(baseData.PaymentTypes))
		j.telemetry.RecordCheckpointLag(ctx, string(store.EntityBaseData), batchLastID-lastSourceID)

		lastSyncedID, err := j.store.GetOperationCheckpoint(ctx, store.OperationBaseDataDeliverCenter)
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
			observedCtx := withDeliveryObserver(ctx, j.store, j.log, store.OperationBaseDataDeliverCenter, paymentTypes,
				func(item struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}) int {
					return item.ID
				},
				func(item struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}) string {
					return fmt.Sprintf("payment_type:%d", item.ID)
				},
			)
			if err := j.target.PostBaseDataToDeliverCenterSaleSelect(observedCtx, domain.BaseData{
				PaymentTypes: paymentTypes,
			}); err != nil {
				return err
			}

			if err := j.store.MarkBatchDelivered(ctx, store.OperationBaseDataDeliverCenter, paymentTypes[len(paymentTypes)-1].ID, buildDeliveredRecords(paymentTypes,
				func(item struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}) int {
					return item.ID
				},
				func(item struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}) string {
					return fmt.Sprintf("payment_type:%d", item.ID)
				},
			)); err != nil {
				return err
			}

			j.telemetry.RecordPosted(ctx, string(store.OperationBaseDataDeliverCenter), len(paymentTypes))
		}

		if err := j.store.AdvanceSourceProgress(ctx, store.EntityBaseData, batchLastID); err != nil {
			return fmt.Errorf("advance base-data source progress: %w", err)
		}

		if len(baseData.PaymentTypes) < j.pageSize {
			return nil
		}
	}
}

func withDeliveryObserver[T any](ctx context.Context, checkpointStore store.CheckpointStore, log *zap.SugaredLogger, operation store.Operation, items []T, idFn func(T) int, dedupeFn func(T) string) context.Context {
	return observability.WithAttemptObserver(ctx, func(attempt observability.HTTPAttempt) {
		recordCtx := observability.WithRunID(context.Background(), observability.RunIDFromContext(ctx))
		for _, item := range items {
			status := store.DeliveryStatusSucceeded
			errorMessage := ""
			if attempt.Error != nil {
				status = store.DeliveryStatusFailed
				errorMessage = attempt.Error.Error()
			}

			if err := checkpointStore.RecordDeliveryAttempt(recordCtx, store.DeliveryAttempt{
				Operation:   operation,
				SourceID:    idFn(item),
				DedupeKey:   dedupeFn(item),
				Status:      status,
				HTTPStatus:  attempt.StatusCode,
				Error:       errorMessage,
				AttemptedAt: time.Now().UTC(),
			}); err != nil && log != nil {
				log.Warnw("failed to persist delivery attempt",
					"run_id", observability.RunIDFromContext(ctx),
					"operation", string(operation),
					"source_id", idFn(item),
					"error", err.Error(),
				)
			}
		}
	})
}

func buildDeliveredRecords[T any](items []T, idFn func(T) int, dedupeFn func(T) string) []store.DeliveredRecord {
	records := make([]store.DeliveredRecord, 0, len(items))
	now := time.Now().UTC()
	for _, item := range items {
		records = append(records, store.DeliveredRecord{
			SourceID:    idFn(item),
			DedupeKey:   dedupeFn(item),
			HTTPStatus:  http.StatusOK,
			DeliveredAt: now,
		})
	}
	return records
}
