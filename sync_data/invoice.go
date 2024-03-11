package syncdata

import (
	"encoding/json"
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
	"log"
	"net/http"
)

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
	var lastId int
	var pageNumber int
	var pageSize int = 1000

	for {

		//TODO: select invoice_progress_info last,page_number
		//TODO: if sql.NoRows {
		// 	lastId = 0
		// 	pageNumber = 0
		// }else {
		// 	lastId = last_id
		// 	pageNumber = page_number + page_size + 1
		// }

		req, err := http.NewRequest(http.MethodGet, i.baseURL+
			fmt.Sprintf(GetInvoices, pageNumber, pageSize, lastId), nil)
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

		if len(response.NewInvoices) == 0 {
			break
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("status code: %d", res.StatusCode)
			return fmt.Errorf(ErrNotOk)
		}

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

		lastId = pageNumber + len(response.NewInvoices)
		//TODO: insert in data base last_id, old_page_number is store in invoice_progress_info

		if len(response.NewInvoices) < pageSize {
			break
		}

	}

	return nil
}
