package fararavand

import (
	"erp-job/config"
	"erp-job/logics"
	"erp-job/models/fararavand"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type FararavandInterface interface {
	GetBaseData() (*fararavand.BaseData, error)
	GetProducts() ([]fararavand.Products, error)
	GetCustomers() ([]fararavand.Customers, error)
	GetInvoices() ([]fararavand.Invoices, error)
	GetTreasuries() ([]fararavand.Treasuries, error)
	GetInvoiceReturns() ([]fararavand.Reverted, error)
}

type Fararavand struct {
	restyClient *resty.Client
	baseUrl     string
}

func NewFararavand() FararavandInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.ApiKey)

	return &Fararavand{
		restyClient: c,
		baseUrl:     config.Cfg.BaseURL,
	}
}

// GetBaseData gets all base information from the first ERP
func (f *Fararavand) GetBaseData() (*fararavand.BaseData, error) {
	var newData = new(fararavand.BaseData)

	resp, err := f.restyClient.R().
		SetResult(newData).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.GetBaseData),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newData, nil
}

// GetProducts gets all products data from the first ERP
func (f *Fararavand) GetProducts() ([]fararavand.Products, error) {
	var newProducts []fararavand.Products

	resp, err := f.restyClient.R().
		SetResult(newProducts).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.GetProducts),
		)
	if err != nil {
		return nil, err
	}
	// get last product id and insert to db
	// check if product id already exits from db and fetch product id bigger than last product id
	//lastId := newProducts[len(newProducts)-1].ID

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newProducts, nil
}

// GetCustomers get all customers' data from the first ERP
func (f *Fararavand) GetCustomers() ([]fararavand.Customers, error) {
	var newCustomers []fararavand.Customers

	resp, err := f.restyClient.R().
		SetResult(newCustomers).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.GetCustomers),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newCustomers, nil
}

// GetInvoices get all invoices' data from the first ERP
func (f *Fararavand) GetInvoices() ([]fararavand.Invoices, error) {
	var newInvoices []fararavand.Invoices

	resp, err := f.restyClient.R().
		SetResult(newInvoices).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.GetInvoices),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newInvoices, nil
}

// GetTreasuries get all treasuries data from the first ERP
func (f *Fararavand) GetTreasuries() ([]fararavand.Treasuries, error) {
	var newTreasuries []fararavand.Treasuries

	resp, err := f.restyClient.R().
		SetResult(newTreasuries).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.GetTreasuries),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newTreasuries, nil
}

// GetInvoiceReturns get all revert invoices data from the first ERP
func (f *Fararavand) GetInvoiceReturns() ([]fararavand.Reverted, error) {
	var newReverted []fararavand.Reverted

	resp, err := f.restyClient.R().
		SetResult(newReverted).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.GetInvoiceReturns),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newReverted, nil
}
