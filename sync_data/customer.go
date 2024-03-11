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
)

// CustomerResponse is the response for the customer
type CustomerResponse struct {
	Status       int                `json:"status"`
	NewCustomers []models.Customers `json:"new_customer"`
}

type Customer struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewCustomer(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Customer {
	return &Customer{
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

		//TODO: select invoice_progress_info last,page_number
		//TODO: if sql.NoRows {
		// 	lastId = 0
		// 	pageNumber = 0
		// }else {
		// 	lastId = last_id
		// 	pageNumber = page_number + page_size + 1
		// }

		req, err := http.NewRequest(http.MethodGet, c.baseURL+
			fmt.Sprintf(GetCustomers, pageNumber, pageSize, lastId), nil)
		if err != nil {
			return err
		}

		req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("get invoice http request failed. status: %d, response: %v", res.StatusCode, res.Body)
		}

		response := new(CustomerResponse)
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			return err
		}

		if len(response.NewCustomers) == 0 {
			break
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("status code: %d", res.StatusCode)
			return fmt.Errorf(utility.ErrNotOk)
		}

		lastId, err = c.fararavand.SyncCustomersWithSaleCustomer(response.NewCustomers)
		if err != nil {
			fmt.Println("Load SyncCustomersWithSaleCustomer encountered an error", err.Error())
			return err
		}

		lastId = pageNumber + len(response.NewCustomers)
		//TODO: insert in data base last_id, old_page_number is store in invoice_progress_info

		if len(response.NewCustomers) < pageSize {
			break
		}
	}

	return nil
}
