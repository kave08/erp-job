package fararavand

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/utility"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Fararavand struct {
	restyClient *resty.Client
	baseUrl     string
	repos       *repository.Repository
	aryan       aryan.AryanInterface
}

func NewFararavand(repos *repository.Repository, aryan aryan.AryanInterface) FararavandInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.FararavandApp.APIKey).SetBaseURL(config.Cfg.FararavandApp.BaseURL)

	return &Fararavand{
		restyClient: c,
		baseUrl:     config.Cfg.FararavandApp.BaseURL,
		repos:       repos,
		aryan:       aryan,
	}
}

// SyncBaseData gets all base information from the first ERP
func (f *Fararavand) SyncBaseData() error {
	var newData = new(models.Fararavand)

	resp, err := f.restyClient.R().
		SetResult(newData).
		Get(utility.FGetBaseData)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// SyncCustomersWithSaleCustomer retrieves all customer data from the Fararavand ERP system and filters them based on the last processed customer ID.
// It fetches the customers using the Fararavand API, then checks the database for the last customer ID that was transferred to the Aryan system.
// If new customers are found (customer with an ID greater than the last processed ID), it sends them to the Aryan system using the PostCustomerToSaleCustomer method.
// The function returns a slice of new customers and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncCustomersWithSaleCustomer() error {
	var newCustomers []models.Customers

	resp, err := f.restyClient.R().SetResult(newCustomers).Get(utility.FGetCustomers)
	if err != nil {
		return err
	}

	lastCustomerId := newCustomers[len(newCustomers)-1].ID
	lastSaleCustomerId, err := f.repos.Database.GetCustomerToSaleCustomer()
	if err != nil {
		return err
	}

	if lastCustomerId > lastSaleCustomerId {
		newCustomers = newCustomers[lastSaleCustomerId:]
		res, err := f.aryan.PostCustomerToSaleCustomer(newCustomers)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertCustomerToSaleCustomer(lastCustomerId)
			if err != nil {
				return err
			}
		}
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// SyncProductsWithGoods retrieves all product data from the Fararavand ERP system and filters them based on the last processed product ID.
// It fetches the products using the Fararavand API, then checks the database for the last product ID that was transferred to the Aryan system.
// If new products are found (products with an ID greater than the last processed ID), it sends them to the Aryan system using the PostProductsToGoods method.
// The function returns a slice of new products and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncProductsWithGoods() error {
	var newProducts []models.Products

	resp, err := f.restyClient.R().SetResult(&newProducts).Get(utility.FGetProducts)
	if err != nil {
		return err
	}

	lastProductId := newProducts[len(newProducts)-1].ID

	lastGoodsId, err := f.repos.Database.GetProductsToGoods()
	if err != nil {
		return err
	}

	if lastProductId > lastGoodsId {
		newProducts = newProducts[lastGoodsId:]
		res, err := f.aryan.PostProductsToGoods(newProducts)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertProductsToGoods(lastProductId)
			if err != nil {
				return err
			}
		}
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// SyncInvoicesWithSaleFactor retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInoviceToSaleFactor method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleFactor(invoices []models.Invoices) error {
	var newInvoices []models.Invoices

	resp, err := f.restyClient.R().SetResult(newInvoices).Get(utility.FGetInvoices)
	if err != nil {
		return err
	}

	lastInvoiceId := newInvoices[len(newInvoices)-1].InvoiceId

	lastSaleFactorId, err := f.repos.Database.GetInvoiceToSaleFactor()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSaleFactorId {
		newInvoices = newInvoices[lastSaleFactorId:]
		res, err := f.aryan.PostInoviceToSaleFactor(newInvoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleFactor(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// SyncInvoicesWithSaleOrder retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSaleOrder method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleOrder(invoices []models.Invoices) error {
	var newInvoices []models.Invoices

	resp, err := f.restyClient.R().SetResult(newInvoices).Get(utility.FGetInvoices)
	if err != nil {
		return err
	}

	lastInvoiceId := newInvoices[len(newInvoices)-1].InvoiceId

	lastSaleOrderId, err := f.repos.Database.GetInvoiceToSaleOrder()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSaleOrderId {
		newInvoices = newInvoices[lastSaleOrderId:]
		res, err := f.aryan.PostInvoiceToSaleOrder(newInvoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleOrder(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// SyncInvoicesWithSalePayment retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSalePayment method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSalePayment(invoices []models.Invoices) error {
	var newInvoices []models.Invoices

	resp, err := f.restyClient.R().SetResult(newInvoices).Get(utility.FGetInvoices)
	if err != nil {
		return err
	}

	lastInvoiceId := newInvoices[len(newInvoices)-1].InvoiceId

	lastSalePaymentId, err := f.repos.Database.GetInvoiceToSalePayment()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalePaymentId {
		newInvoices = newInvoices[lastSalePaymentId:]
		res, err := f.aryan.PostInvoiceToSalePayment(newInvoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSalePayment(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// SyncInvoicesWithSalerSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSalerSelect method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSalerSelect(invoices []models.Invoices) error {
	var newInvoices []models.Invoices

	resp, err := f.restyClient.R().SetResult(newInvoices).Get(utility.FGetInvoices)
	if err != nil {
		return err
	}

	lastInvoiceId := newInvoices[len(newInvoices)-1].InvoiceId

	lastSalerSelectId, err := f.repos.Database.GetInvoiceToSalerSelect()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		newInvoices = newInvoices[lastSalerSelectId:]
		res, err := f.aryan.PostInvoiceToSalerSelect(newInvoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSalerSelect(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// SyncInvoicesWithSaleProforma retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSaleProforma method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleProforma(invoices []models.Invoices) error {
	var newInvoices []models.Invoices

	resp, err := f.restyClient.R().SetResult(newInvoices).Get(utility.FGetInvoices)
	if err != nil {
		return err
	}

	lastInvoiceId := newInvoices[len(newInvoices)-1].InvoiceId

	lastSaleProformaId, err := f.repos.Database.GetInvoiceToSaleProforma()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSaleProformaId {
		newInvoices = newInvoices[lastSaleProformaId:]
		res, err := f.aryan.PostInvoiceToSaleProforma(newInvoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleProforma(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// GetTreasuries get all treasuries data from the first ERP
func (f *Fararavand) SyncTreasuries() error {
	var newTreasuries []models.Fararavand

	resp, err := f.restyClient.R().SetResult(newTreasuries).Get(utility.FGetTreasuries)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}

// GetInvoiceReturns get all revert invoices data from the first ERP
func (f *Fararavand) SyncInvoiceReturns() error {
	var newReverted []models.Fararavand

	resp, err := f.restyClient.R().SetResult(newReverted).Get(utility.FGetInvoiceReturns)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return fmt.Errorf(utility.ErrNotOk)
	}

	return nil
}
