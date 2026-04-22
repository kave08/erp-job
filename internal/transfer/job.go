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
	entityKey func(domain.Invoices) (string, error)
}

type preparedBatch[T any] struct {
	pending          []deliveryCandidate[T]
	checkpointCursor int
	lastSourceCursor int
}

type deliveryCandidate[T any] struct {
	item         T
	sourceCursor int
	entityKey    string
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
		{name: "base_data", entity: store.EntityBaseData, run: j.runBaseData},
		{name: "customer", entity: store.EntityCustomer, run: j.runCustomers},
		{name: "product", entity: store.EntityProduct, run: j.runProducts},
		{name: "invoice", entity: store.EntityInvoice, run: j.runInvoices},
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
	return runEntitySync(j, ctx, store.EntityInvoice, j.source.FetchInvoices, func(item domain.Invoices) int {
		return item.InvoiceId
	}, j.syncInvoices)
}

func (j *Job) syncInvoices(ctx context.Context, invoices []domain.Invoices) error {
	operations := []invoiceOperation{
		{
			name:      "sale payment",
			operation: store.OperationInvoiceSalePayment,
			post:      j.target.PostInvoiceToSalePayment,
			entityKey: noErrEntityKey(func(item domain.Invoices) string { return fmt.Sprintf("payment:%d", item.SNoePardakht) }),
		},
		{
			name:      "sale center",
			operation: store.OperationInvoiceSaleCenter,
			post:      j.target.PostInvoiceToSaleCenter,
			entityKey: noErrEntityKey(func(item domain.Invoices) string { return fmt.Sprintf("stock:%d", item.WareHouseID) }),
		},
		{
			name:      "saler select",
			operation: store.OperationInvoiceSalerSelect,
			post:      j.target.PostInvoiceToSalerSelect,
			entityKey: func(item domain.Invoices) (string, error) {
				visitorID, err := domain.ParseVisitorCode(item.VisitorCode)
				if err != nil {
					return "", err
				}
				return fmt.Sprintf("visitor:%d", visitorID), nil
			},
		},
		{
			name:      "sale type select",
			operation: store.OperationInvoiceSaleTypeSelect,
			post:      j.target.PostInvoiceToSaleTypeSelect,
			entityKey: noErrEntityKey(func(item domain.Invoices) string { return fmt.Sprintf("sale_type:%d", item.SNoePardakht) }),
		},
		{
			name:      "sale factor",
			operation: store.OperationInvoiceSaleFactor,
			post:      j.target.PostInvoiceToSaleFactor,
			entityKey: noErrEntityKey(func(item domain.Invoices) string {
				return fmt.Sprintf("invoice:%d:product:%d", item.InvoiceId, item.ProductID)
			}),
		},
		{
			name:      "sale order",
			operation: store.OperationInvoiceSaleOrder,
			post:      j.target.PostInvoiceToSaleOrder,
			entityKey: noErrEntityKey(func(item domain.Invoices) string {
				return fmt.Sprintf("invoice:%d:product:%d", item.InvoiceId, item.ProductID)
			}),
		},
		{
			name:      "sale proforma",
			operation: store.OperationInvoiceSaleProforma,
			post:      j.target.PostInvoiceToSaleProforma,
			entityKey: noErrEntityKey(func(item domain.Invoices) string {
				return fmt.Sprintf("invoice:%d:product:%d", item.InvoiceId, item.ProductID)
			}),
		},
	}

	for _, op := range operations {
		prepared, err := prepareBatch(j, ctx, op.operation, invoices, func(item domain.Invoices) int {
			return item.InvoiceId
		}, op.entityKey)
		if err != nil {
			return fmt.Errorf("prepare %s batch: %w", op.name, err)
		}
		if prepared.lastSourceCursor == 0 {
			continue
		}

		j.telemetry.RecordCheckpointLag(ctx, string(store.EntityInvoice), prepared.lastSourceCursor-prepared.checkpointCursor)

		if err := deliverBatch(j, ctx, "invoice", op.name, op.operation, prepared, op.post); err != nil {
			return err
		}
	}

	return nil
}

func (j *Job) runCustomers(ctx context.Context) error {
	return runEntitySync(j, ctx, store.EntityCustomer, j.source.FetchCustomers, func(item domain.Customers) int {
		return item.ID
	}, func(ctx context.Context, customers []domain.Customers) error {
		return syncOperation(j, ctx, "customer", store.EntityCustomer, "sale customer", store.OperationCustomerSaleCustomer, customers,
			func(item domain.Customers) int { return item.ID },
			noErrEntityKey(func(item domain.Customers) string { return fmt.Sprintf("customer:%d", item.ID) }),
			j.target.PostCustomerToSaleCustomer,
		)
	})
}

func (j *Job) runProducts(ctx context.Context) error {
	return runEntitySync(j, ctx, store.EntityProduct, j.source.FetchProducts, func(item domain.Products) int {
		return item.ID
	}, func(ctx context.Context, products []domain.Products) error {
		return syncOperation(j, ctx, "product", store.EntityProduct, "goods", store.OperationProductsGoods, products,
			func(item domain.Products) int { return item.ID },
			noErrEntityKey(func(item domain.Products) string { return fmt.Sprintf("product:%d", item.ID) }),
			j.target.PostProductsToGoods,
		)
	})
}

func (j *Job) runBaseData(ctx context.Context) error {
	return runEntitySync(j, ctx, store.EntityBaseData, func(ctx context.Context, pageNumber, pageSize, lastID int) ([]struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}, error) {
		baseData, err := j.source.FetchBaseData(ctx, pageNumber, pageSize, lastID)
		if err != nil {
			return nil, err
		}

		return baseData.PaymentTypes, nil
	}, func(item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}) int {
		return item.ID
	}, func(ctx context.Context, paymentTypes []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}) error {
		return syncOperation(j, ctx, "base_data", store.EntityBaseData, "deliver center sale select", store.OperationBaseDataDeliverCenter, paymentTypes,
			func(item struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}) int {
				return item.ID
			},
			noErrEntityKey(func(item struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}) string {
				return fmt.Sprintf("payment_type:%d", item.ID)
			}),
			func(ctx context.Context, paymentTypes []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}) error {
				return j.target.PostBaseDataToDeliverCenterSaleSelect(ctx, domain.BaseData{PaymentTypes: paymentTypes})
			},
		)
	})
}

func runEntitySync[T any](j *Job, ctx context.Context, entity store.Entity, fetch func(context.Context, int, int, int) ([]T, error), sourceCursorFn func(T) int, sync func(context.Context, []T) error) error {
	for {
		lastSourceID, err := j.store.GetSourceProgress(ctx, entity)
		if err != nil {
			return err
		}

		items, err := fetch(ctx, cursorPageNumber, j.pageSize, lastSourceID)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		batchLastID := sourceCursorFn(items[len(items)-1])
		if batchLastID <= lastSourceID {
			return fmt.Errorf("%s cursor did not advance: last_source_id=%d batch_last_id=%d", entity, lastSourceID, batchLastID)
		}

		j.telemetry.RecordFetched(ctx, string(entity), len(items))
		j.telemetry.RecordCheckpointLag(ctx, string(entity), batchLastID-lastSourceID)

		if err := sync(ctx, items); err != nil {
			return err
		}

		if err := j.store.AdvanceSourceProgress(ctx, entity, batchLastID); err != nil {
			return fmt.Errorf("advance %s source progress: %w", entity, err)
		}

		if len(items) < j.pageSize {
			return nil
		}
	}
}

func syncOperation[T any](j *Job, ctx context.Context, step string, entity store.Entity, name string, operation store.Operation, items []T, sourceCursorFn func(T) int, entityKeyFn func(T) (string, error), post func(context.Context, []T) error) error {
	prepared, err := prepareBatch(j, ctx, operation, items, sourceCursorFn, entityKeyFn)
	if err != nil {
		return fmt.Errorf("prepare %s batch: %w", name, err)
	}
	if prepared.lastSourceCursor == 0 {
		return nil
	}

	j.telemetry.RecordCheckpointLag(ctx, string(entity), prepared.lastSourceCursor-prepared.checkpointCursor)

	return deliverBatch(j, ctx, step, name, operation, prepared, post)
}

func prepareBatch[T any](j *Job, ctx context.Context, operation store.Operation, items []T, sourceCursorFn func(T) int, entityKeyFn func(T) (string, error)) (preparedBatch[T], error) {
	lastCheckpoint, err := j.store.GetOperationCheckpoint(ctx, operation)
	if err != nil {
		return preparedBatch[T]{}, fmt.Errorf("get operation checkpoint for %s: %w", operation, err)
	}

	trimmed := trimAfterCheckpoint(items, lastCheckpoint, sourceCursorFn)
	if len(trimmed) == 0 {
		return preparedBatch[T]{checkpointCursor: lastCheckpoint}, nil
	}

	candidates, entityKeys, err := dedupeCandidates(trimmed, sourceCursorFn, entityKeyFn)
	if err != nil {
		return preparedBatch[T]{}, err
	}

	deliveredKeys, err := j.store.GetDeliveredEntityKeys(ctx, operation, entityKeys)
	if err != nil {
		return preparedBatch[T]{}, fmt.Errorf("get delivered entity keys for %s: %w", operation, err)
	}

	pending := make([]deliveryCandidate[T], 0, len(candidates))
	for _, candidate := range candidates {
		if _, delivered := deliveredKeys[candidate.entityKey]; delivered {
			continue
		}
		pending = append(pending, candidate)
	}

	return preparedBatch[T]{
		pending:          pending,
		checkpointCursor: lastCheckpoint,
		lastSourceCursor: sourceCursorFn(trimmed[len(trimmed)-1]),
	}, nil
}

func deliverBatch[T any](j *Job, ctx context.Context, step, name string, operation store.Operation, prepared preparedBatch[T], post func(context.Context, []T) error) error {
	if len(prepared.pending) > 0 {
		observedCtx := withDeliveryObserver(ctx, j.store, j.log, operation, prepared.pending)
		if err := post(observedCtx, itemsFromCandidates(prepared.pending)); err != nil {
			return fmt.Errorf("post %s: %w", name, err)
		}

		j.telemetry.RecordPosted(ctx, string(operation), len(prepared.pending))

		if j.log != nil {
			j.log.Infow("completed target batch",
				"run_id", observability.RunIDFromContext(ctx),
				"step", step,
				"operation", string(operation),
				"batch_size", len(prepared.pending),
				"last_source_id", prepared.lastSourceCursor,
				"status_code", http.StatusOK,
			)
		}
	}

	if err := j.store.MarkBatchDelivered(ctx, operation, prepared.lastSourceCursor, buildDeliveredRecords(prepared.pending)); err != nil {
		return fmt.Errorf("save %s checkpoint: %w", name, err)
	}

	return nil
}

func withDeliveryObserver[T any](ctx context.Context, checkpointStore store.CheckpointStore, log *zap.SugaredLogger, operation store.Operation, candidates []deliveryCandidate[T]) context.Context {
	return observability.WithAttemptObserver(ctx, func(attempt observability.HTTPAttempt) {
		recordCtx := context.WithoutCancel(ctx)
		for _, candidate := range candidates {
			status := store.DeliveryStatusSucceeded
			errorMessage := ""
			if attempt.Error != nil {
				status = store.DeliveryStatusFailed
				errorMessage = attempt.Error.Error()
			}

			if err := checkpointStore.RecordDeliveryAttempt(recordCtx, store.DeliveryAttempt{
				Operation:    operation,
				SourceCursor: candidate.sourceCursor,
				EntityKey:    candidate.entityKey,
				Status:       status,
				HTTPStatus:   attempt.StatusCode,
				Error:        errorMessage,
				AttemptedAt:  time.Now().UTC(),
			}); err != nil && log != nil {
				log.Warnw("failed to persist delivery attempt",
					"run_id", observability.RunIDFromContext(ctx),
					"operation", string(operation),
					"source_cursor", candidate.sourceCursor,
					"entity_key", candidate.entityKey,
					"error", err.Error(),
				)
			}
		}
	})
}

func dedupeCandidates[T any](items []T, sourceCursorFn func(T) int, entityKeyFn func(T) (string, error)) ([]deliveryCandidate[T], []string, error) {
	seen := make(map[string]struct{}, len(items))
	candidates := make([]deliveryCandidate[T], 0, len(items))

	for index := len(items) - 1; index >= 0; index-- {
		entityKey, err := entityKeyFn(items[index])
		if err != nil {
			return nil, nil, err
		}
		if _, exists := seen[entityKey]; exists {
			continue
		}

		seen[entityKey] = struct{}{}
		candidates = append(candidates, deliveryCandidate[T]{
			item:         items[index],
			sourceCursor: sourceCursorFn(items[index]),
			entityKey:    entityKey,
		})
	}

	for left, right := 0, len(candidates)-1; left < right; left, right = left+1, right-1 {
		candidates[left], candidates[right] = candidates[right], candidates[left]
	}

	entityKeys := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		entityKeys = append(entityKeys, candidate.entityKey)
	}

	return candidates, entityKeys, nil
}

func itemsFromCandidates[T any](candidates []deliveryCandidate[T]) []T {
	items := make([]T, 0, len(candidates))
	for _, candidate := range candidates {
		items = append(items, candidate.item)
	}

	return items
}

func buildDeliveredRecords[T any](candidates []deliveryCandidate[T]) []store.DeliveredRecord {
	records := make([]store.DeliveredRecord, 0, len(candidates))
	now := time.Now().UTC()
	for _, candidate := range candidates {
		records = append(records, store.DeliveredRecord{
			SourceCursor: candidate.sourceCursor,
			EntityKey:    candidate.entityKey,
			HTTPStatus:   http.StatusOK,
			DeliveredAt:  now,
		})
	}
	return records
}

func noErrEntityKey[T any](fn func(T) string) func(T) (string, error) {
	return func(item T) (string, error) {
		return fn(item), nil
	}
}
