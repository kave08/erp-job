package aryan

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"erp-job/internal/config"
	"erp-job/internal/domain"
	"erp-job/internal/observability"
	"erp-job/internal/retry"
)

func TestPostInvoiceToSaleFactorSendsMappedPayload(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []domain.SaleFactor
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy())
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/SaleFactor" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if got := r.Header.Get("ApiKey"); got != cfg.APIKey {
				t.Fatalf("unexpected api key header: %s", got)
			}
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("decode request body: %v", err)
			}

			return newResponse(http.StatusOK, ""), nil
		}),
	}

	invoice := domain.Invoices{
		CustomerID:      11,
		InvoiceDate:     "14030101",
		WareHouseID:     22,
		CodeMahal:       33,
		SNoePardakht:    44,
		CCForoshandeh:   55,
		CodeForoshandeh: 66,
		InvoiceId:       77,
		ProductID:       88,
		ProductCount:    99,
		ProductFee:      111,
		TozihatFaktor:   "invoice detail",
	}

	if err := client.PostInvoiceToSaleFactor(context.Background(), []domain.Invoices{invoice}); err != nil {
		t.Fatalf("PostInvoiceToSaleFactor returned error: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 sale factor payload, got %d", len(received))
	}

	got := received[0]
	if got.CustomerID != invoice.CustomerID {
		t.Fatalf("unexpected customer id: %d", got.CustomerID)
	}
	if got.StockID != invoice.WareHouseID {
		t.Fatalf("unexpected stock id: %d", got.StockID)
	}
	if got.SaleTypeID != invoice.SNoePardakht {
		t.Fatalf("unexpected sale type id: %d", got.SaleTypeID)
	}
	if got.DeliveryCenterID != invoice.SNoePardakht {
		t.Fatalf("unexpected delivery center id: %d", got.DeliveryCenterID)
	}
	if got.SaleCenterID != invoice.WareHouseID {
		t.Fatalf("unexpected sale center id: %d", got.SaleCenterID)
	}
	if got.PaymentWayID != invoice.SNoePardakht {
		t.Fatalf("unexpected payment way id: %d", got.PaymentWayID)
	}
	if got.SellerID != invoice.CCForoshandeh {
		t.Fatalf("unexpected seller id: %d", got.SellerID)
	}
	if got.SaleManID != invoice.CodeForoshandeh {
		t.Fatalf("unexpected sale man id: %d", got.SaleManID)
	}
	if got.SecondNumber != "77" {
		t.Fatalf("unexpected second number: %s", got.SecondNumber)
	}
	if got.ServiceGoodsID != invoice.ProductID {
		t.Fatalf("unexpected service goods id: %d", got.ServiceGoodsID)
	}
	if got.Quantity != float64(invoice.ProductCount) {
		t.Fatalf("unexpected quantity: %f", got.Quantity)
	}
	if got.Fee != float64(invoice.ProductFee) {
		t.Fatalf("unexpected fee: %f", got.Fee)
	}
	if got.DetailDesc != invoice.TozihatFaktor {
		t.Fatalf("unexpected detail description: %s", got.DetailDesc)
	}
}

func TestPostInvoiceToSaleOrderUsesInvoiceIDAsSecondNumber(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []domain.SaleOrder
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy())
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostInvoiceToSaleOrder(context.Background(), []domain.Invoices{{
		InvoiceId:     91,
		CustomerID:    1,
		WareHouseID:   22,
		SNoePardakht:  44,
		ProductID:     55,
		ProductCount:  2,
		ProductFee:    3,
		VisitorCode:   "66",
		TozihatFaktor: "detail",
	}}); err != nil {
		t.Fatalf("PostInvoiceToSaleOrder returned error: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 sale order payload, got %#v", received)
	}

	if received[0].SecondNumber != 91 {
		t.Fatalf("expected sale order second number 91, got %#v", received)
	}
	if received[0].SaleTypeId != 44 {
		t.Fatalf("unexpected sale type id: %d", received[0].SaleTypeId)
	}
	if received[0].DeliveryCenterID != 44 {
		t.Fatalf("unexpected delivery center id: %d", received[0].DeliveryCenterID)
	}
	if received[0].SaleCenterID != 22 {
		t.Fatalf("unexpected sale center id: %d", received[0].SaleCenterID)
	}
	if received[0].SellerVisitorID != 66 {
		t.Fatalf("unexpected seller visitor id: %d", received[0].SellerVisitorID)
	}
}

func TestPostInvoiceToSaleProformaAlignsReferenceIDs(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []domain.SaleProforma
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy())
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostInvoiceToSaleProforma(context.Background(), []domain.Invoices{{
		InvoiceId:     91,
		CustomerID:    1,
		WareHouseID:   22,
		SNoePardakht:  44,
		ProductID:     55,
		ProductCount:  2,
		ProductFee:    3,
		VisitorCode:   "66",
		TozihatFaktor: "detail",
	}}); err != nil {
		t.Fatalf("PostInvoiceToSaleProforma returned error: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 sale proforma payload, got %#v", received)
	}

	if received[0].SaleTypeId != 44 {
		t.Fatalf("unexpected sale type id: %d", received[0].SaleTypeId)
	}
	if received[0].DeliveryCenterID != 44 {
		t.Fatalf("unexpected delivery center id: %d", received[0].DeliveryCenterID)
	}
	if received[0].SaleCenterID != 22 {
		t.Fatalf("unexpected sale center id: %d", received[0].SaleCenterID)
	}
}

func TestPostInvoiceToSalePaymentUsesSNoePardakht(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []domain.SalePaymentSelect
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy())
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostInvoiceToSalePayment(context.Background(), []domain.Invoices{{
		PaymentTypeID:  7,
		SNoePardakht:   44,
		TxtNoePardakht: "cash",
	}}); err != nil {
		t.Fatalf("PostInvoiceToSalePayment returned error: %v", err)
	}

	if len(received) != 1 || received[0].PaymentWayID != 44 {
		t.Fatalf("expected payment way id 44 from SNoePardakht, got %#v", received)
	}
}

func TestPostInvoiceToSalerSelectRejectsInvalidVisitorCodeBeforeHTTP(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	requests := 0
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy())
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			requests++
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostInvoiceToSalerSelect(context.Background(), []domain.Invoices{{VisitorCode: "bad-code"}}); err == nil {
		t.Fatal("expected invalid visitor code error")
	}

	if requests != 0 {
		t.Fatalf("expected no HTTP request on invalid visitor code, got %d", requests)
	}
}

func TestPostInvoiceToSaleTypeSelectUsesStableCodeFromSNoePardakht(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []domain.SaleTypeSelect
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy())
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostInvoiceToSaleTypeSelect(context.Background(), []domain.Invoices{{
		SNoePardakht:   44,
		Codekala:       "should-not-be-used",
		TxtNoePardakht: "cash",
	}}); err != nil {
		t.Fatalf("PostInvoiceToSaleTypeSelect returned error: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 sale type payload, got %#v", received)
	}
	if received[0].BuySaleTypeCode != "44" {
		t.Fatalf("unexpected sale type code: %q", received[0].BuySaleTypeCode)
	}
	if received[0].BuySaleTypeDesc != "cash" {
		t.Fatalf("unexpected sale type desc: %q", received[0].BuySaleTypeDesc)
	}
}

func TestPostInvoiceToSaleFactorRetriesTransientFailure(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	attempts := 0
	client := newClient(cfg, testTelemetry(t), nil, retry.Policy{
		MaxAttempts:    3,
		InitialBackoff: 0,
		MaxBackoff:     0,
		Multiplier:     1,
		Jitter:         0,
	})
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				return newResponse(http.StatusServiceUnavailable, "temporary failure\n"), nil
			}
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostInvoiceToSaleFactor(context.Background(), []domain.Invoices{{InvoiceId: 10}}); err != nil {
		t.Fatalf("PostInvoiceToSaleFactor returned error: %v", err)
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestPostInvoiceToSaleFactorDoesNotRetryPermanentFailure(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	attempts := 0
	client := newClient(cfg, testTelemetry(t), nil, retry.Policy{
		MaxAttempts:    3,
		InitialBackoff: 0,
		MaxBackoff:     0,
		Multiplier:     1,
		Jitter:         0,
	})
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			attempts++
			return newResponse(http.StatusBadRequest, "bad request\n"), nil
		}),
	}

	if err := client.PostInvoiceToSaleFactor(context.Background(), []domain.Invoices{{InvoiceId: 10}}); err == nil {
		t.Fatal("expected permanent failure")
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt for permanent failure, got %d", attempts)
	}
}

func testTelemetry(t *testing.T) *observability.Telemetry {
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func newResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

const serverURLPlaceholder = "http://invalid.example/"
