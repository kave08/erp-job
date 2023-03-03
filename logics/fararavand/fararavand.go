package fararavand

import (
	"database/sql"
	"erp-job/config"
	"erp-job/logics"
	"erp-job/models/fararavand"
	"erp-job/repository"
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
	repos       *repository.Repository
}

func NewFararavand() FararavandInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.ApiKey)

	return &Fararavand{
		restyClient: c,
		baseUrl:     config.Cfg.BaseURL,
	}
}

func NewLogics(repos *repository.Repository) *Fararavand {
	return &Fararavand{
		repos: repos,
	}
}

// GetBaseData gets all base information from the first ERP
func (f *Fararavand) GetBaseData() (*fararavand.BaseData, error) {
	var newData = new(fararavand.BaseData)

	resp, err := f.restyClient.R().
		SetResult(newData).
		Get(
			fmt.Sprintf("%s/%s", f.baseUrl, logics.FararavandGetBaseData),
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
			fmt.Sprintf("%s/%s", f.baseUrl, logics.FararavandGetProducts),
		)
	if err != nil {
		return nil, err
	}

	lastId := newProducts[len(newProducts)-1].ID
	// get last product id
	pId, err := f.repos.Database.GetProduct(lastId)
	if err != nil {
		return nil, err
	}
	//check if product is empty and insert
	if pId == 0 {
		err = f.repos.Database.InsertProduct(lastId)
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// check if product id already exits from db and fetch product id bigger than last product id
	if lastId > pId {
		newProducts = newProducts[:lastId-pId]
		//afeter get last product id and send data to db, insert new product id into db
		err = f.repos.Database.InsertProduct(lastId)
		if err != sql.ErrNoRows {
			return nil, err
		}
		return newProducts, err
	}

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
			fmt.Sprintf("%s/%s", f.baseUrl, logics.FararavandGetCustomers),
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
			fmt.Sprintf("%s/%s", f.baseUrl, logics.FararavandGetInvoices),
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
			fmt.Sprintf("%s/%s", f.baseUrl, logics.FararavandGetTreasuries),
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
			fmt.Sprintf("%s/%s", f.baseUrl, logics.FararavandGetInvoiceReturns),
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
