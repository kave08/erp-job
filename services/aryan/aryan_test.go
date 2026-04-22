package aryan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"erp-job/config"
	"erp-job/models"
)

func TestPostInvoiceToSaleFactorSendsMappedPayload(t *testing.T) {
	t.Parallel()

	config.Cfg.AryanApp.APIKey = "test-api-key"

	var received []models.SaleFactor
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/SaleFactor" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if got := r.Header.Get("ApiKey"); got != config.Cfg.AryanApp.APIKey {
			t.Fatalf("unexpected api key header: %s", got)
		}

		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := &Aryan{
		baseURL:    server.URL + "/",
		httpClient: server.Client(),
	}

	invoice := models.Invoices{
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

	if err := service.PostInvoiceToSaleFactor([]models.Invoices{invoice}); err != nil {
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
	if got.SaleCenterID != invoice.CodeMahal {
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
