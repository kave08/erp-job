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
	GetBaseData() (*models.Fararavand, error)
	GetInvoicesForSaleFactor() ([]models.Invoices, error)
	GetProducts() ([]models.Products, error)
	GetCustomers() ([]models.Customers, error)
	GetTreasuries() ([]models.Fararavand, error)
	GetInvoiceReturns() ([]models.Fararavand, error)
}

type Fararavand struct {
	restyClient *resty.Client
	baseUrl     string
	repos       *repository.Repository
	aryan       AryanInterface
}

func NewFararavand(repos *repository.Repository, aryan AryanInterface) FararavandInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &Fararavand{
		restyClient: c,
		baseUrl:     config.Cfg.FararavandApp.BaseURL,
		repos:       repos,
		aryan:       aryan,
	}
}

// GetBaseData gets all base information from the first ERP
func (f *Fararavand) GetBaseData() (*models.Fararavand, error) {
	var newData = new(models.Fararavand)

	resp, err := f.restyClient.R().
		SetResult(newData).
		Get(FGetBaseData)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newData, nil
}

// GetCustomers get all customers' data from the first ERP
func (f *Fararavand) GetCustomers() ([]models.Customers, error) {
	var newCustomers []models.Customers

	resp, err := f.restyClient.R().SetResult(newCustomers).Get(FGetCustomers)
	if err != nil {
		return nil, err
	}

	lastId := newCustomers[len(newCustomers)-1].ID
	cId, err := f.repos.Database.GetCustomerToSaleCustomer()
	if err != nil {
		return nil, err
	}

	if lastId > cId {
		newCustomers = newCustomers[cId:]
		res, err := f.aryan.PostCustomerToSaleCustomer(newCustomers)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertCustomerToSaleCustomer(lastId)
			if err != nil {
				return nil, err
			}
		}
		return newCustomers, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newCustomers, nil
}

// GetProducts gets all products data from the first ERP
func (f *Fararavand) GetProducts() ([]models.Products, error) {
	var newProducts []models.Products

	resp, err := f.restyClient.R().SetResult(&newProducts).Get(FGetProducts)
	if err != nil {
		return nil, err
	}

	lastId := newProducts[len(newProducts)-1].ID
	pId, err := f.repos.Database.GetProductsToGoods()
	if err != nil {
		return nil, err
	}

	if lastId > pId {
		newProducts = newProducts[pId:]
		res, err := f.aryan.PostProductsToGoods(newProducts)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertProductsToGoods(lastId)
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

// GetInvoices get all invoices' data from the first ERP
func (f *Fararavand) GetInvoicesForSaleFactor() ([]models.Invoices, error) {
	var newInvoices []models.Invoices

	resp, err := f.restyClient.R().SetResult(newInvoices).Get(FGetInvoices)
	if err != nil {
		return nil, err
	}

	lastInvoiceId := newInvoices[len(newInvoices)-1].InvoiceId

	iToSaleFactorId, err := f.repos.Database.GetInvoiceToSaleFactor()
	if err != nil {
		return nil, err
	}

	if lastInvoiceId > iToSaleFactorId {
		newInvoices = newInvoices[iToSaleFactorId:]
		res, err := f.aryan.PostInoviceToSaleFactor(newInvoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleFactor(lastInvoiceId)
			if err != nil {
				return nil, err
			}
		}
		return newInvoices, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newInvoices, nil
}

// GetTreasuries get all treasuries data from the first ERP
func (f *Fararavand) GetTreasuries() ([]models.Fararavand, error) {
	var newTreasuries []models.Fararavand

	resp, err := f.restyClient.R().SetResult(newTreasuries).Get(FGetTreasuries)
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
func (f *Fararavand) GetInvoiceReturns() ([]models.Fararavand, error) {
	var newReverted []models.Fararavand

	resp, err := f.restyClient.R().SetResult(newReverted).Get(FGetInvoiceReturns)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(ErrNotOk)
	}

	return newReverted, nil
}
