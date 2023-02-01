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
	GetProducts() (*fararavand.Products, error)
	GetCustomers() (*fararavand.Customers, error)
	GetInvoices() (*fararavand.Invoices, error)
	GetTreasuries() (*fararavand.Treasuries, error)
	GetInvoiceReturns() (*fararavand.Reverted, error)
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

func (f *Fararavand) GetProducts() (*fararavand.Products, error) {
	var newProducts = new(fararavand.Products)

	resp, err := f.restyClient.R().
		SetResult(newProducts).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.GetProducts),
		)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newProducts, nil
}

func (f *Fararavand) GetCustomers() (*fararavand.Customers, error) {
	var newCustomers = new(fararavand.Customers)

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

func (f *Fararavand) GetInvoices() (*fararavand.Invoices, error) {
	var newInvoices = new(fararavand.Invoices)

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

func (f *Fararavand) GetTreasuries() (*fararavand.Treasuries, error) {
	var newTreasuries = new(fararavand.Treasuries)

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

func (f *Fararavand) GetInvoiceReturns() (*fararavand.Reverted, error) {
	var newReverted = new(fararavand.Reverted)

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
