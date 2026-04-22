package transfer

import (
	"context"
	"errors"
	"testing"

	"erp-job/internal/config"
	"erp-job/internal/domain"
	"erp-job/internal/observability"
	"erp-job/internal/store"
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

func TestJobUsesLastSourceIDProgressAndPersistsDeliveredRecords(t *testing.T) {
	t.Parallel()

	telemetry := testTransferTelemetry(t)
	checkpoints := newFakeStore()
	source := &fakeSource{
		invoices: map[int][]domain.Invoices{
			0: {
				{InvoiceId: 10, ProductID: 100, PaymentTypeID: 7, VisitorCode: "12", WareHouseID: 20, SNoePardakht: 30},
				{InvoiceId: 20, ProductID: 101, PaymentTypeID: 8, VisitorCode: "13", WareHouseID: 21, SNoePardakht: 31},
			},
			20: nil,
		},
	}
	target := newFakeTarget(true, "")

	job := NewJob(checkpoints, source, target, nil, telemetry)
	ctx := observability.WithRunID(context.Background(), "run-1")
	if err := job.Run(ctx); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if got := checkpoints.sourceProgress[store.EntityInvoice]; got != 20 {
		t.Fatalf("expected invoice source progress 20, got %d", got)
	}

	if got := checkpoints.operationCheckpoints[store.OperationInvoiceSaleFactor]; got != 20 {
		t.Fatalf("expected sale factor checkpoint 20, got %d", got)
	}

	if got := len(checkpoints.deliveryAttempts); got == 0 {
		t.Fatal("expected delivery attempts to be persisted")
	}
}

func TestJobSkipsRecordsAlreadyBehindOperationCheckpoint(t *testing.T) {
	t.Parallel()

	telemetry := testTransferTelemetry(t)
	checkpoints := newFakeStore()
	checkpoints.operationCheckpoints[store.OperationInvoiceSaleFactor] = 10
	source := &fakeSource{
		invoices: map[int][]domain.Invoices{
			0: {
				{InvoiceId: 10, ProductID: 100, PaymentTypeID: 7, VisitorCode: "12", WareHouseID: 20, SNoePardakht: 30},
				{InvoiceId: 20, ProductID: 101, PaymentTypeID: 8, VisitorCode: "13", WareHouseID: 21, SNoePardakht: 31},
			},
			20: nil,
		},
	}
	target := newFakeTarget(true, "")

	job := NewJob(checkpoints, source, target, nil, telemetry)
	if err := job.Run(observability.WithRunID(context.Background(), "run-2")); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	saleFactorBatch := target.invoiceOps[string(store.OperationInvoiceSaleFactor)]
	if len(saleFactorBatch) != 1 || saleFactorBatch[0][0].InvoiceId != 20 {
		t.Fatalf("expected trimmed invoice batch, got %#v", saleFactorBatch)
	}
}

func TestJobDoesNotAdvanceSourceProgressWhenOperationFails(t *testing.T) {
	t.Parallel()

	telemetry := testTransferTelemetry(t)
	checkpoints := newFakeStore()
	source := &fakeSource{
		invoices: map[int][]domain.Invoices{
			0: {
				{InvoiceId: 10, ProductID: 100, PaymentTypeID: 7, VisitorCode: "12", WareHouseID: 20, SNoePardakht: 30},
				{InvoiceId: 20, ProductID: 101, PaymentTypeID: 8, VisitorCode: "13", WareHouseID: 21, SNoePardakht: 31},
			},
		},
	}
	target := newFakeTarget(true, string(store.OperationInvoiceSaleOrder))

	job := NewJob(checkpoints, source, target, nil, telemetry)
	err := job.Run(observability.WithRunID(context.Background(), "run-3"))
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

type fakeStore struct {
	sourceProgress       map[store.Entity]int
	operationCheckpoints map[store.Operation]int
	deliveryAttempts     []store.DeliveryAttempt
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		sourceProgress:       make(map[store.Entity]int),
		operationCheckpoints: make(map[store.Operation]int),
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

func (f *fakeStore) RecordDeliveryAttempt(_ context.Context, attempt store.DeliveryAttempt) error {
	f.deliveryAttempts = append(f.deliveryAttempts, attempt)
	return nil
}

func (f *fakeStore) MarkBatchDelivered(_ context.Context, operation store.Operation, lastSourceID int, _ []store.DeliveredRecord) error {
	if lastSourceID > f.operationCheckpoints[operation] {
		f.operationCheckpoints[operation] = lastSourceID
	}
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
	invoiceOps   map[string][][]domain.Invoices
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

func (f *fakeTarget) PostProductsToGoods(context.Context, []domain.Products) error { return nil }

func (f *fakeTarget) PostCustomerToSaleCustomer(context.Context, []domain.Customers) error {
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

func (f *fakeTarget) PostBaseDataToDeliverCenterSaleSelect(context.Context, domain.BaseData) error {
	return nil
}

func (f *fakeTarget) recordInvoiceOp(ctx context.Context, operation string, invoices []domain.Invoices) error {
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

	copied := append([]domain.Invoices(nil), invoices...)
	f.invoiceOps[operation] = append(f.invoiceOps[operation], copied)
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
