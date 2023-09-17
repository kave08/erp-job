package logics

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type FararavandInterface interface {
	GetBaseData() (*models.FararavandBaseData, error)
	GetProducts() ([]models.FararavandProducts, error)
	GetCustomers() ([]models.FararavandCustomers, error)
	GetInvoices() ([]models.FararavandInvoices, error)
	GetTreasuries() ([]models.FararavandTreasuries, error)
	GetInvoiceReturns() ([]models.FararavandReverted, error)
}

type Fararavand struct {
	restyClient *resty.Client
	baseUrl     string
	repos       *repository.Repository
	aryan       AryanInterface
}

func NewFararavand(repos *repository.Repository, aryan AryanInterface) FararavandInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.ApiKey)

	return &Fararavand{
		restyClient: c,
		baseUrl:     config.Cfg.BaseURL,
		repos:       repos,
		aryan:       aryan,
	}
}

// GetBaseData gets all base information from the first ERP
func (f *Fararavand) GetBaseData() (*models.FararavandBaseData, error) {
	var newData = new(models.FararavandBaseData)

	resp, err := f.restyClient.R().
		SetResult(newData).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, FGetBaseData),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newData, nil
}

// GetProducts gets all products data from the first ERP
func (f *Fararavand) GetProducts() ([]models.FararavandProducts, error) {
	var newProducts []models.FararavandProducts

	resp, err := f.restyClient.R().
		SetResult(newProducts).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, FGetProducts),
		)
	if err != nil {
		return nil, err
	}

	// get last product id from response --100
	lastId := newProducts[len(newProducts)-1].ID
	// get last product id from data --80
	pId, err := f.repos.Database.GetProduct()
	if err != nil {
		return nil, err
	}
	// fetch new product id
	if lastId > pId {
		newProducts = newProducts[pId:]
		//insert new product id into db
		res, err := f.aryan.PostGoods(newProducts)

		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertProduct(lastId)
			if err != nil {
				return nil, err
			}
		}
		return newProducts, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newProducts, nil
}

// GetCustomers get all customers' data from the first ERP
func (f *Fararavand) GetCustomers() ([]models.FararavandCustomers, error) {
	var newCustomers []models.FararavandCustomers

	resp, err := f.restyClient.R().
		SetResult(newCustomers).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, FGetCustomers),
		)
	if err != nil {
		return nil, err
	}

	// get last customer id from response
	lastId := newCustomers[len(newCustomers)-1].ID
	// get last customer id from data
	pId, err := f.repos.Database.GetCustomer()
	if err != nil {
		return nil, err
	}
	// if customer id is empty insert in db
	if pId == 0 {
		err = f.repos.Database.InsertCustomer(lastId)
		if err != nil {
			return nil, err
		}
	}
	// fetch new customer id
	if lastId > pId {
		newCustomers = newCustomers[pId:]
		//insert new customer id into db
		err = f.repos.Database.InsertCustomer(lastId)
		if err != nil {
			return nil, err
		}
		return newCustomers, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newCustomers, nil
}

// GetInvoices get all invoices' data from the first ERP
func (f *Fararavand) GetInvoices() ([]models.FararavandInvoices, error) {
	var newInvoices []models.FararavandInvoices

	resp, err := f.restyClient.R().
		SetResult(newInvoices).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, FGetInvoices),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newInvoices, nil
}

// GetTreasuries get all treasuries data from the first ERP
func (f *Fararavand) GetTreasuries() ([]models.FararavandTreasuries, error) {
	var newTreasuries []models.FararavandTreasuries

	resp, err := f.restyClient.R().
		SetResult(newTreasuries).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, FGetTreasuries),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newTreasuries, nil
}

// GetInvoiceReturns get all revert invoices data from the first ERP
func (f *Fararavand) GetInvoiceReturns() ([]models.FararavandReverted, error) {
	var newReverted []models.FararavandReverted

	resp, err := f.restyClient.R().
		SetResult(newReverted).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, FGetInvoiceReturns),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newReverted, nil
}
