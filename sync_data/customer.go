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

// CustomerResponse is the response for the customer
type CustomerResponse struct {
	Status       int                `json:"status"`
	NewCustomers []models.Customers `json:"new_customer"`
}

type Customer struct {
	log        *zap.SugaredLogger
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.Interface
}

func NewCustomer(repos *repository.Repository, fr fararavand.Interface, ar aryan.AryanInterface) *Customer {
	return &Customer{
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

func (c Customer) Customers() error {
	var lastId int
	var pageNumber int
	var pageSize int = 1000

	for {
		lastCustomerId, lastPageNumber, err := c.repos.Database.GetCustomerProgress()
		if err == sql.ErrNoRows {
			lastId = 0
			pageNumber = 0
		} else {
			lastId = lastCustomerId
			pageNumber = lastPageNumber + pageSize + 1
		}

		req, err := http.NewRequest(http.MethodGet, c.baseURL+
			fmt.Sprintf(GetCustomers, pageNumber, pageSize, lastId), nil)
		if err != nil {
			c.log.Errorw("get customer request encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			c.log.Errorw("get customer http request failed.",
				"error", err,
				"status:", res.StatusCode,
				"response", res.Body,
			)

			return fmt.Errorf("get customer http request failed. status: %d, response: %v", res.StatusCode, res.Body)
		}

		response := new(CustomerResponse)
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			c.log.Errorw("get customer decode response encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)
			return err
		}

		if len(response.NewCustomers) == 0 {
			break
		}

		err = c.fararavand.SyncCustomersWithSaleCustomer(response.NewCustomers)
		if err != nil {
			c.log.Errorw("load SyncCustomersWithSaleCustomer encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		lastId = pageNumber + len(response.NewCustomers)

		err = c.repos.Database.InsertCustomerProgress(lastId, pageNumber)
		if err != nil {
			c.log.Errorw("load InsertCustomerProgress encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if len(response.NewCustomers) < pageSize {
			break
		}
	}

	return nil
}
