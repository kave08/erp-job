package transfer

import (
	"context"
	"errors"
	"testing"

	"erp-job/internal/config"
	"erp-job/internal/domain"
	"erp-job/internal/observability"
	"erp-job/internal/store"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestTrimAfterCheckpointReturnsFirstRecordAfterCheckpoint(t *testing.T) {
	t.Parallel()

	ids := []int{2, 4, 6, 9}
	trimmed := trimAfterCheckpoint(ids, 4, func(item int) int {
		return item
	})

	if len(trimmed) != 2 || trimmed[0] != 6 || trimmed[1] != 9 {
		t.Fatalf("unexpected trimmed result: %#v", trimmed)
	}
}

func TestTrimAfterCheckpointReturnsNilWhenBatchAlreadySynced(t *testing.T) {
	t.Parallel()

	ids := []int{2, 4, 6, 9}
	trimmed := trimAfterCheckpoint(ids, 9, func(item int) int {
		return item
	})

	if trimmed != nil {
		t.Fatalf("expected nil for fully synced batch, got %#v", trimmed)
	}
}

func TestJobRunsMasterDataBeforeInvoicesAndInvoiceReferencesBeforeDocuments(t *testing.T) {
	t.Parallel()

	job, _, target := newTestJob(t)
	job.source = &fakeSource{
		invoices: map[int][]domain.Invoices{
			0: {
				{InvoiceId: 10, CustomerID: 1, ProductID: 100, VisitorCode: "12", WareHouseID: 20, SNoePardakht: 30},
			},
			10: nil,
		},
		customers: map[int][]domain.Customers{
			0: {{ID: 1, Code: 11}},
			1: nil,
		},
		products: map[int][]domain.Products{
			0:   {{ID: 100, Code: "P-100"}},
			100: nil,
		},
		baseData: map[int]domain.BaseData{
			0: {
				PaymentTypes: []struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{{ID: 30, Name: "cash"}},
			},
			30: {},
		},
	}

	if err := job.Run(observability.WithRunID(context.Background(), "run-order")); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	expected := []string{
		string(store.OperationBaseDataDeliverCenter),
		string(store.OperationCustomerSaleCustomer),
		string(store.OperationProductsGoods),
		string(store.OperationInvoiceSalePayment),
		string(store.OperationInvoiceSaleCenter),
		string(store.OperationInvoiceSalerSelect),
		string(store.OperationInvoiceSaleTypeSelect),
		string(store.OperationInvoiceSaleFactor),
		string(store.OperationInvoiceSaleOrder),
		string(store.OperationInvoiceSaleProforma),
	}
	if len(target.calls) != len(expected) {
		t.Fatalf("unexpected call count: got %d want %d (%#v)", len(target.calls), len(expected), target.calls)
	}
	for index, call := range expected {
		if target.calls[index] != call {
			t.Fatalf("unexpected call order at %d: got %q want %q", index, target.calls[index], call)
		}
	}
}

func TestInvoiceOperationsDeduplicateByEntityKeyAndSkipDeliveredRows(t *testing.T) {
	t.Parallel()

	job, checkpoints, target := newTestJob(t)
	checkpoints.deliveredKeys[store.OperationInvoiceSalePayment] = map[string]struct{}{
		"payment:30": {},
	}
	job.source = &fakeSource{
		invoices: map[int][]domain.Invoices{
			0: {
				{InvoiceId: 10, ProductID: 100, PaymentTypeID: 7, VisitorCode: "12", WareHouseID: 20, SNoePardakht: 30},
				{InvoiceId: 20, ProductID: 101, PaymentTypeID: 8, VisitorCode: "13", WareHouseID: 21, SNoePardakht: 40},
				{InvoiceId: 30, ProductID: 102, PaymentTypeID: 9, VisitorCode: "14", WareHouseID: 22, SNoePardakht: 40},
			},
			30: nil,
		},
	}

	if err := job.Run(observability.WithRunID(context.Background(), "run-dedupe")); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	batches := target.invoiceOps[string(store.OperationInvoiceSalePayment)]
	if len(batches) != 1 {
		t.Fatalf("expected one sale payment batch, got %d", len(batches))
	}
	if len(batches[0]) != 1 || batches[0][0].InvoiceId != 30 {
		t.Fatalf("expected only the last undelivered payment row to be posted, got %#v", batches[0])
	}
	if got := checkpoints.operationCheckpoints[store.OperationInvoiceSalePayment]; got != 30 {
		t.Fatalf("expected sale payment checkpoint 30, got %d", got)
	}
}

func TestSyncOperationAdvancesCheckpointWhenAllEntityKeysAlreadyDelivered(t *testing.T) {
	t.Parallel()

	job, checkpoints, target := newTestJob(t)
	checkpoints.deliveredKeys[store.OperationProductsGoods] = map[string]struct{}{
		"product:10": {},
		"product:20": {},
	}
	job.source = &fakeSource{
		products: map[int][]domain.Products{
			0: {
				{ID: 10, Code: "P-10"},
				{ID: 20, Code: "P-20"},
			},
			20: nil,
		},
	}

	if err := job.Run(observability.WithRunID(context.Background(), "run-products")); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := len(target.productOps); got != 0 {
		t.Fatalf("expected skipped product post, got %d product batches", got)
	}
	if got := checkpoints.operationCheckpoints[store.OperationProductsGoods]; got != 20 {
		t.Fatalf("expected products checkpoint 20, got %d", got)
	}
	if got := checkpoints.sourceProgress[store.EntityProduct]; got != 20 {
		t.Fatalf("expected products source progress 20, got %d", got)
	}
}

func TestJobDoesNotAdvanceSourceProgressWhenOperationFails(t *testing.T) {
	t.Parallel()

	job, checkpoints, _ := newTestJob(t)
	job.source = &fakeSource{
		invoices: map[int][]domain.Invoices{
			0: {
				{InvoiceId: 10, ProductID: 100, PaymentTypeID: 7, VisitorCode: "12", WareHouseID: 20, SNoePardakht: 30},
				{InvoiceId: 20, ProductID: 101, PaymentTypeID: 8, VisitorCode: "13", WareHouseID: 21, SNoePardakht: 31},
			},
		},
	}
	job.target.(*fakeTarget).failOp = string(store.OperationInvoiceSaleOrder)

	err := job.Run(observability.WithRunID(context.Background(), "run-failure"))
	if err == nil {
		t.Fatal("expected run failure")
	}

	if got := checkpoints.sourceProgress[store.EntityInvoice]; got != 0 {
		t.Fatalf("expected invoice source progress to stay at 0, got %d", got)
	}
	if got := checkpoints.operationCheckpoints[store.OperationInvoiceSaleFactor]; got != 20 {
		t.Fatalf("expected completed sale factor checkpoint 20, got %d", got)
	}
}

func TestSyncInvoicesWarnsOnConflictingSaleTypeDescriptions(t *testing.T) {
	t.Parallel()

	core, logs := observer.New(zap.WarnLevel)
	logger := zap.New(core).Sugar()
	t.Cleanup(func() {
		_ = logger.Sync()
	})

	telemetry := testTransferTelemetry(t)
	checkpoints := newFakeStore()
	target := newFakeTarget(true, "")
	job := NewJob(checkpoints, &fakeSource{}, target, logger, telemetry)

	err := job.syncInvoices(observability.WithRunID(context.Background(), "run-sale-type-warn"), []domain.Invoices{
		{InvoiceId: 10, ProductID: 100, VisitorCode: "12", WareHouseID: 20, SNoePardakht: 30, TxtNoePardakht: "cash"},
		{InvoiceId: 20, ProductID: 101, VisitorCode: "13", WareHouseID: 21, SNoePardakht: 30, TxtNoePardakht: "credit"},
	})
	if err != nil {
		t.Fatalf("syncInvoices returned error: %v", err)
	}

	entries := logs.FilterMessage("conflicting sale type descriptions in invoice batch").All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 conflict warning, got %d", len(entries))
	}
}

func newTestJob(t *testing.T) (*Job, *fakeStore, *fakeTarget) {
	t.Helper()

	telemetry := testTransferTelemetry(t)
	checkpoints := newFakeStore()
	target := newFakeTarget(true, "")
	job := NewJob(checkpoints, &fakeSource{}, target, nil, telemetry)
	return job, checkpoints, target
}

type fakeStore struct {
	sourceProgress       map[store.Entity]int
	operationCheckpoints map[store.Operation]int
	deliveredKeys        map[store.Operation]map[string]struct{}
	deliveryAttempts     []store.DeliveryAttempt
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		sourceProgress:       make(map[store.Entity]int),
		operationCheckpoints: make(map[store.Operation]int),
		deliveredKeys:        make(map[store.Operation]map[string]struct{}),
	}
}

func (f *fakeStore) GetSourceProgress(_ context.Context, entity store.Entity) (int, error) {
	return f.sourceProgress[entity], nil
}

func (f *fakeStore) AdvanceSourceProgress(_ context.Context, entity store.Entity, lastSourceID int) error {
	if lastSourceID > f.sourceProgress[entity] {
		f.sourceProgress[entity] = lastSourceID
	}
	return nil
}

func (f *fakeStore) GetOperationCheckpoint(_ context.Context, operation store.Operation) (int, error) {
	return f.operationCheckpoints[operation], nil
}

func (f *fakeStore) GetDeliveredEntityKeys(_ context.Context, operation store.Operation, entityKeys []string) (map[string]struct{}, error) {
	delivered := make(map[string]struct{}, len(entityKeys))
	for _, entityKey := range entityKeys {
		if _, exists := f.deliveredKeys[operation][entityKey]; exists {
			delivered[entityKey] = struct{}{}
		}
	}
	return delivered, nil
}

func (f *fakeStore) RecordDeliveryAttempt(_ context.Context, attempt store.DeliveryAttempt) error {
	f.deliveryAttempts = append(f.deliveryAttempts, attempt)
	return nil
}

func (f *fakeStore) GetAttemptCounts(_ context.Context, _ store.Operation, entityKeys []string) (map[string]int, error) {
	counts := make(map[string]int, len(entityKeys))
	return counts, nil
}

func (f *fakeStore) MarkBatchDelivered(_ context.Context, operation store.Operation, lastSourceID int, records []store.DeliveredRecord) error {
	if lastSourceID > f.operationCheckpoints[operation] {
		f.operationCheckpoints[operation] = lastSourceID
	}
	if f.deliveredKeys[operation] == nil {
		f.deliveredKeys[operation] = make(map[string]struct{}, len(records))
	}
	for _, record := range records {
		f.deliveredKeys[operation][record.EntityKey] = struct{}{}
	}
	return nil
}

func (f *fakeStore) MarkPermanentFailures(_ context.Context, _ store.Operation, _ []string) error {
	return nil
}

type fakeSource struct {
	invoices  map[int][]domain.Invoices
	customers map[int][]domain.Customers
	products  map[int][]domain.Products
	baseData  map[int]domain.BaseData
}

func (f *fakeSource) FetchInvoices(_ context.Context, _, _ int, lastID int) ([]domain.Invoices, error) {
	return f.invoices[lastID], nil
}

func (f *fakeSource) FetchCustomers(_ context.Context, _, _ int, lastID int) ([]domain.Customers, error) {
	return f.customers[lastID], nil
}

func (f *fakeSource) FetchProducts(_ context.Context, _, _ int, lastID int) ([]domain.Products, error) {
	return f.products[lastID], nil
}

func (f *fakeSource) FetchBaseData(_ context.Context, _, _ int, lastID int) (domain.BaseData, error) {
	return f.baseData[lastID], nil
}

type fakeTarget struct {
	emitAttempts bool
	failOp       string
	calls        []string
	invoiceOps   map[string][][]domain.Invoices
	customerOps  [][]domain.Customers
	productOps   [][]domain.Products
	baseDataOps  []domain.BaseData
}

func newFakeTarget(emitAttempts bool, failOp string) *fakeTarget {
	return &fakeTarget{
		emitAttempts: emitAttempts,
		failOp:       failOp,
		invoiceOps:   make(map[string][][]domain.Invoices),
	}
}

func (f *fakeTarget) PostInvoiceToSaleFactor(ctx context.Context, invoices []domain.Invoices) error {
	return f.recordInvoiceOp(ctx, string(store.OperationInvoiceSaleFactor), invoices)
}

func (f *fakeTarget) PostProductsToGoods(ctx context.Context, products []domain.Products) error {
	if err := f.emitResult(ctx, string(store.OperationProductsGoods)); err != nil {
		return err
	}

	f.calls = append(f.calls, string(store.OperationProductsGoods))
	f.productOps = append(f.productOps, append([]domain.Products(nil), products...))
	return nil
}

func (f *fakeTarget) PostCustomerToSaleCustomer(ctx context.Context, customers []domain.Customers) error {
	if err := f.emitResult(ctx, string(store.OperationCustomerSaleCustomer)); err != nil {
		return err
	}

	f.calls = append(f.calls, string(store.OperationCustomerSaleCustomer))
	f.customerOps = append(f.customerOps, append([]domain.Customers(nil), customers...))
	return nil
}

func (f *fakeTarget) PostInvoiceToSaleOrder(ctx context.Context, invoices []domain.Invoices) error {
	return f.recordInvoiceOp(ctx, string(store.OperationInvoiceSaleOrder), invoices)
}

func (f *fakeTarget) PostInvoiceToSalePayment(ctx context.Context, invoices []domain.Invoices) error {
	return f.recordInvoiceOp(ctx, string(store.OperationInvoiceSalePayment), invoices)
}

func (f *fakeTarget) PostInvoiceToSaleCenter(ctx context.Context, invoices []domain.Invoices) error {
	return f.recordInvoiceOp(ctx, string(store.OperationInvoiceSaleCenter), invoices)
}

func (f *fakeTarget) PostInvoiceToSalerSelect(ctx context.Context, invoices []domain.Invoices) error {
	return f.recordInvoiceOp(ctx, string(store.OperationInvoiceSalerSelect), invoices)
}

func (f *fakeTarget) PostInvoiceToSaleProforma(ctx context.Context, invoices []domain.Invoices) error {
	return f.recordInvoiceOp(ctx, string(store.OperationInvoiceSaleProforma), invoices)
}

func (f *fakeTarget) PostInvoiceToSaleTypeSelect(ctx context.Context, invoices []domain.Invoices) error {
	return f.recordInvoiceOp(ctx, string(store.OperationInvoiceSaleTypeSelect), invoices)
}

func (f *fakeTarget) PostBaseDataToDeliverCenterSaleSelect(ctx context.Context, baseData domain.BaseData) error {
	if err := f.emitResult(ctx, string(store.OperationBaseDataDeliverCenter)); err != nil {
		return err
	}

	f.calls = append(f.calls, string(store.OperationBaseDataDeliverCenter))
	f.baseDataOps = append(f.baseDataOps, baseData)
	return nil
}

func (f *fakeTarget) recordInvoiceOp(ctx context.Context, operation string, invoices []domain.Invoices) error {
	if err := f.emitResult(ctx, operation); err != nil {
		return err
	}

	f.calls = append(f.calls, operation)
	copied := append([]domain.Invoices(nil), invoices...)
	f.invoiceOps[operation] = append(f.invoiceOps[operation], copied)
	return nil
}

func (f *fakeTarget) emitResult(ctx context.Context, operation string) error {
	if operation == f.failOp {
		if f.emitAttempts {
			if observer := observability.AttemptObserverFromContext(ctx); observer != nil {
				observer(observability.HTTPAttempt{Endpoint: operation, Attempt: 1, StatusCode: 500, Error: errors.New("target operation failed")})
			}
		}
		return errors.New("target operation failed")
	}

	if f.emitAttempts {
		if observer := observability.AttemptObserverFromContext(ctx); observer != nil {
			observer(observability.HTTPAttempt{Endpoint: operation, Attempt: 1, StatusCode: 200})
		}
	}

	return nil
}

func testTransferTelemetry(t *testing.T) *observability.Telemetry {
	t.Helper()

	telemetry, shutdown, err := observability.New(context.Background(), config.OTel{})
	if err != nil {
		t.Fatalf("create telemetry: %v", err)
	}
	t.Cleanup(func() {
		_ = shutdown(context.Background())
	})

	return telemetry
}
