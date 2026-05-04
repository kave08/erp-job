package aryan

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"erp-job/internal/circuitbreaker"
	"erp-job/internal/config"
	"erp-job/internal/domain"
	"erp-job/internal/observability"
	"erp-job/internal/retry"
)

func paramByName(params []ParamEntry, name string) interface{} {
	for _, p := range params {
		if p.Name == name {
			if p.ArrayValue != nil {
				return p.ArrayValue
			}
			return p.Value
		}
	}
	return nil
}

func TestPostInvoiceToSaleFactorSendsMappedPayload(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []ParamsPayload
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
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
		CustomerID:    11,
		InvoiceDate:   "14030101",
		WareHouseID:   22,
		SNoePardakht:  44,
		CCForoshandeh: 55,
		InvoiceId:     77,
		ProductID:     88,
		ProductCount:  99,
		ProductFee:    111,
		TozihatFaktor: "invoice detail",
	}

	if err := client.PostInvoiceToSaleFactor(context.Background(), []domain.Invoices{invoice}); err != nil {
		t.Fatalf("PostInvoiceToSaleFactor returned error: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 sale factor payload, got %d", len(received))
	}

	got := received[0]
	if got.ID != "SaleFactor" {
		t.Fatalf("unexpected payload id: %s", got.ID)
	}

	params := got.Params

	if v := paramByName(params, "CustomerId"); v != float64(11) {
		t.Fatalf("unexpected customer id: %v", v)
	}
	if v := paramByName(params, "StockId"); v != float64(22) {
		t.Fatalf("unexpected stock id: %v", v)
	}
	if v := paramByName(params, "SaleTypeId"); v != float64(44) {
		t.Fatalf("unexpected sale type id: %v", v)
	}
	if v := paramByName(params, "DeliveryCenterID"); v != float64(44) {
		t.Fatalf("unexpected delivery center id: %v", v)
	}
	if v := paramByName(params, "SaleCenterID"); v != float64(22) {
		t.Fatalf("unexpected sale center id: %v", v)
	}
	if v := paramByName(params, "PaymentWayID"); v != float64(44) {
		t.Fatalf("unexpected payment way id: %v", v)
	}
	if v := paramByName(params, "SellerVisitorID"); v != float64(55) {
		t.Fatalf("unexpected seller visitor id: %v", v)
	}
	if v := paramByName(params, "SecondNumber"); v != float64(77) {
		t.Fatalf("unexpected second number: %v", v)
	}
	if v := paramByName(params, "VoucherDesc"); v != "ETL-Form Fararavand" {
		t.Fatalf("unexpected voucher desc: %v", v)
	}
	if v := paramByName(params, "isrow"); v != "1" {
		t.Fatalf("unexpected isrow: %v", v)
	}

	inserted, ok := paramByName(params, "[Inserted]").([]interface{})
	if !ok || len(inserted) != 5 {
		t.Fatalf("unexpected [Inserted]: %v", paramByName(params, "[Inserted]"))
	}
	if inserted[1] != float64(88) {
		t.Fatalf("unexpected product id in [Inserted]: %v", inserted[1])
	}
	if inserted[2] != float64(99) {
		t.Fatalf("unexpected quantity in [Inserted]: %v", inserted[2])
	}
	if inserted[3] != float64(111) {
		t.Fatalf("unexpected fee in [Inserted]: %v", inserted[3])
	}

	elInserted, ok := paramByName(params, "[el_Inserted]").([]interface{})
	if !ok || len(elInserted) != 3 {
		t.Fatalf("unexpected [el_Inserted]: %v", paramByName(params, "[el_Inserted]"))
	}
	if elInserted[0] != float64(55) {
		t.Fatalf("unexpected seller id in [el_Inserted]: %v", elInserted[0])
	}
}

func TestPostInvoiceToSaleOrderUsesInvoiceIDAsSecondNumber(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []ParamsPayload
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
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
		t.Fatalf("expected 1 sale order payload, got %d", len(received))
	}

	params := received[0].Params
	if v := paramByName(params, "SecondNumber"); v != float64(91) {
		t.Fatalf("expected second number 91, got %v", v)
	}
	if v := paramByName(params, "SaleTypeId"); v != float64(44) {
		t.Fatalf("unexpected sale type id: %v", v)
	}
	if v := paramByName(params, "DeliveryCenterID"); v != float64(44) {
		t.Fatalf("unexpected delivery center id: %v", v)
	}
	if v := paramByName(params, "SaleCenterID"); v != float64(22) {
		t.Fatalf("unexpected sale center id: %v", v)
	}
	if v := paramByName(params, "SellerVisitorID"); v != float64(66) {
		t.Fatalf("unexpected seller visitor id: %v", v)
	}
}

func TestPostInvoiceToSaleProformaAlignsReferenceIDs(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []ParamsPayload
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
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
		t.Fatalf("expected 1 sale proforma payload, got %d", len(received))
	}

	params := received[0].Params
	if v := paramByName(params, "SaleTypeId"); v != float64(44) {
		t.Fatalf("unexpected sale type id: %v", v)
	}
	if v := paramByName(params, "DeliveryCenterID"); v != float64(44) {
		t.Fatalf("unexpected delivery center id: %v", v)
	}
	if v := paramByName(params, "SaleCenterID"); v != float64(22) {
		t.Fatalf("unexpected sale center id: %v", v)
	}
	if v := paramByName(params, "SellerVisitorID"); v != float64(66) {
		t.Fatalf("unexpected seller visitor id: %v", v)
	}
}

func TestPostInvoiceToSalePaymentUsesSNoePardakht(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received []ParamsPayload
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostInvoiceToSalePayment(context.Background(), []domain.Invoices{{
		SNoePardakht:   44,
		TxtNoePardakht: "cash",
	}}); err != nil {
		t.Fatalf("PostInvoiceToSalePayment returned error: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 payment payload, got %d", len(received))
	}
	if v := paramByName(received[0].Params, "PaymentWayID"); v != float64(44) {
		t.Fatalf("expected payment way id 44, got %v", v)
	}
	if v := paramByName(received[0].Params, "PaymentwayDesc"); v != "cash" {
		t.Fatalf("expected payment way desc 'cash', got %v", v)
	}
}

func TestPostInvoiceToSalerSelectRejectsInvalidVisitorCodeBeforeHTTP(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	requests := 0
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
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

	var received []ParamsPayload
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
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
		TxtNoePardakht: "cash",
	}}); err != nil {
		t.Fatalf("PostInvoiceToSaleTypeSelect returned error: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 sale type payload, got %d", len(received))
	}
	if v := paramByName(received[0].Params, "BuySaleTypeCode"); v != "44" {
		t.Fatalf("unexpected sale type code: %v", v)
	}
	if v := paramByName(received[0].Params, "BuySaleTypeDesc"); v != "cash" {
		t.Fatalf("unexpected sale type desc: %v", v)
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
	}, circuitbreaker.New(5, 30*time.Second))
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
	}, circuitbreaker.New(5, 30*time.Second))
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

func TestPostAddSubElementSelectSendsParamsFormat(t *testing.T) {
	t.Parallel()

	cfg := config.AryanApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	var received ParamsPayload
	client := newClient(cfg, testTelemetry(t), nil, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/AddSubElementSelect" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			return newResponse(http.StatusOK, ""), nil
		}),
	}

	if err := client.PostAddSubElementSelect(context.Background()); err != nil {
		t.Fatalf("PostAddSubElementSelect returned error: %v", err)
	}

	if received.ID != "AddSubElementSelect" {
		t.Fatalf("unexpected payload id: %s", received.ID)
	}
	if v := paramByName(received.Params, "IsStruct"); v != float64(1) {
		t.Fatalf("unexpected IsStruct value: %v", v)
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
