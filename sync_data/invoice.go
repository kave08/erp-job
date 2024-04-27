package syncdata

import (
	"database/sql"
	"encoding/json"
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"erp-job/utility/logger"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// InvoiceResponse is the response for the invoice
type InvoiceResponse struct {
	Status      int               `json:"status"`
	NewInvoices []models.Invoices `json:"new_invoice"`
}

type Invoice struct {
	log        *zap.SugaredLogger
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.Interface
}

func NewInvoice(repos *repository.Repository, fr fararavand.Interface, ar aryan.AryanInterface) *Invoice {
	return &Invoice{
		log:        logger.Logger(),
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
		lastInvoiceId, lastPageNumber, err := i.repos.Database.GetInvoiceProgress()
		if err == sql.ErrNoRows {
			lastId = 0
			pageNumber = 0
		} else {
			lastId = lastInvoiceId
			pageNumber = lastPageNumber + pageSize + 1
		}

		req, err := http.NewRequest(http.MethodGet, i.baseURL+
			fmt.Sprintf(GetInvoices, pageNumber, pageSize, lastId), nil)
		if err != nil {
			i.log.Errorw("get invoice request encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

		res, err := i.httpClient.Do(req)
		if err != nil {
			i.log.Errorw("get invoice response encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if res.StatusCode != http.StatusOK {
			i.log.Errorw("get invoice http request failed.",
				"error", err,
				"status:", res.StatusCode,
				"response", res.Body,
			)

			return fmt.Errorf("get invoice http request failed. status: %d, response: %v", res.StatusCode, res.Body)
		}

		response := new(InvoiceResponse)
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			i.log.Errorw("get invoice decode response encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if len(response.NewInvoices) == 0 {
			break
		}

		err = i.fararavand.SyncInvoicesWithSaleFactor(response.NewInvoices)
		if err != nil {
			i.log.Errorw("load SyncInvoicesWithSaleFactor encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		err = i.fararavand.SyncInvoicesWithSaleOrder(response.NewInvoices)
		if err != nil {
			i.log.Errorw("load SyncInvoicesWithSaleOrder encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		err = i.fararavand.SyncInvoicesWithSalePayment(response.NewInvoices)
		if err != nil {
			i.log.Errorw("load SyncInvoicesWithSalePayment encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		err = i.fararavand.SyncInvoicesWithSalerSelect(response.NewInvoices)
		if err != nil {
			i.log.Errorw("load SyncInvoicesWithSalerSelect encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		err = i.fararavand.SyncInvoicesWithSaleProforma(response.NewInvoices)
		if err != nil {
			i.log.Errorw("load SyncInvoicesWithSaleProforma encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		err = i.fararavand.SyncInvoicesWithSaleCenter(response.NewInvoices)
		if err != nil {
			i.log.Errorw("load SyncInvoicesWithSaleCenter encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		err = i.fararavand.SyncInvoiceWithSaleTypeSelect(response.NewInvoices)
		if err != nil {
			i.log.Errorw("load SyncInvoiceWithSaleTypeSelect encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		lastId = pageNumber + len(response.NewInvoices)

		err = i.repos.Database.InsertInvoiceProgress(lastId, pageNumber)
		if err != nil {
			i.log.Errorw("InsertInvoiceProgress encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if len(response.NewInvoices) < pageSize {
			break
		}

	}

	return nil
}
