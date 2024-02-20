package syncdata

import (
	"encoding/json"
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"erp-job/utility"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var lastId int
var pageNumber int
var pageSize int = 1000
var mutex *sync.Mutex

// InvoiceResponse is the response for the invoice
type InvoiceResponse struct {
	Status      int               `json:"status"`
	NewInvoices []models.Invoices `json:"new_invoice"`
}

type Invoice struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewInvoice(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Invoice {
	return &Invoice{
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (i *Invoice) Invoices() error {

	req, err := http.NewRequest(http.MethodGet, i.baseURL+
		fmt.Sprintf("/GetInvoices?PageNumeber=%d&PageSize=%d&LastId=%d/", pageNumber, pageSize, lastId), nil)
	if err != nil {
		return err
	}

	req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

	res, err := i.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("get invoice http request failed. status: %d, response: %v", res.StatusCode, res.Body)
	}

	response := new(InvoiceResponse)
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return err
	}

	if res.StatusCode != response.Status {
		return fmt.Errorf("get invoice http request failed(body). status: %d, response: %v", response.Status, res.Body)
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("status code: %d", res.StatusCode)
		return fmt.Errorf(utility.ErrNotOk)
	}

	// if lastId <= 0 {
	// 	return fmt.Errorf("validation.required %d", http.StatusBadRequest)
	// }

	err = i.fararavand.SyncInvoicesWithSaleFactor(response.NewInvoices)
	if err != nil {
		fmt.Println("load SyncInvoicesWithSaleFactor encountered an error: %w", err)
		return err
	}

	err = i.fararavand.SyncInvoicesWithSaleOrder(response.NewInvoices)
	if err != nil {
		fmt.Println("load SyncInvoicesWithSaleOrder encountered an error: %w", err)
		return err
	}

	err = i.fararavand.SyncInvoicesWithSalePayment(response.NewInvoices)
	if err != nil {
		fmt.Println("load SyncInvoicesWithSalePayment encountered an error: %w", err)
		return err
	}

	err = i.fararavand.SyncInvoicesWithSalerSelect(response.NewInvoices)
	if err != nil {
		fmt.Println("load SyncInvoicesWithSalerSelect encountered an error: %w", err)
		return err
	}

	err = i.fararavand.SyncInvoicesWithSaleProforma(response.NewInvoices)
	if err != nil {
		fmt.Println("load SyncInvoicesWithSaleProforma encountered an error: %w", err)
		return err
	}

	err = i.fararavand.SyncInvoicesWithSaleCenter(response.NewInvoices)
	if err != nil {
		fmt.Println("load SyncInvoicesWithSaleCenter encountered an error: %w", err)
		return err
	}

	err = i.fararavand.SyncInvoiceWithSaleTypeSelect(response.NewInvoices)
	if err != nil {
		fmt.Println("load SyncInvoicesWithSaleCenter encountered an error: %w", err)
		return err
	}

	return nil
}
