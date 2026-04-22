package fararavand

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"erp-job/internal/config"
	"erp-job/internal/observability"
	"erp-job/internal/retry"
)

func TestFetchInvoicesRetriesTransientFailure(t *testing.T) {
	t.Parallel()

	cfg := config.FararavandApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	attempts := 0
	client := newClient(cfg, testSourceTelemetry(t), nil, retry.Policy{
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
			return newResponse(http.StatusOK, `{"status":200,"new_invoice":[{"invoiceId":10}]}`), nil
		}),
	}

	invoices, err := client.FetchInvoices(context.Background(), 0, 100, 0)
	if err != nil {
		t.Fatalf("FetchInvoices returned error: %v", err)
	}

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
	if len(invoices) != 1 || invoices[0].InvoiceId != 10 {
		t.Fatalf("unexpected invoices: %#v", invoices)
	}
}

func TestFetchInvoicesDoesNotRetryPermanentFailure(t *testing.T) {
	t.Parallel()

	cfg := config.FararavandApp{
		BaseURL: serverURLPlaceholder,
		APIKey:  "test-api-key",
	}

	attempts := 0
	client := newClient(cfg, testSourceTelemetry(t), nil, retry.Policy{
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

	if _, err := client.FetchInvoices(context.Background(), 0, 100, 0); err == nil {
		t.Fatal("expected permanent failure")
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func testSourceTelemetry(t *testing.T) *observability.Telemetry {
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
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	return &http.Response{
		StatusCode: statusCode,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

const serverURLPlaceholder = "http://invalid.example/"
