package syncdata

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"erp-job/utility"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Customer struct {
	restyClient *resty.Client
	baseURL     string
	httpClient  *http.Client
	repos       *repository.Repository
	aryan       aryan.AryanInterface
	fararavand  fararavand.FararavandInterface
}

func NewCustomer(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Customer {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &Customer{
		restyClient: c,
		baseURL:     config.Cfg.FararavandApp.BaseURL,
		repos:       repos,
		aryan:       ar,
		fararavand:  fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (c Customer) Customers() error {

	var newCustomers []models.Customers

	resp, err := c.restyClient.R().SetResult(newCustomers).Get(utility.FGetCustomers)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	lastId, err = c.fararavand.SyncCustomersWithSaleCustomer(newCustomers)
	if err != nil {
		fmt.Println("Load SyncCustomersWithSaleCustomer encountered an error", err.Error())
		return err
	}

	return nil
}
