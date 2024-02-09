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


// SyncBaseDataWithDeliverCenter retrieves all base data from the Fararavand ERP system and filters them based on the last processed payment type ID.
// It fetches the base data using the Fararavand API, then checks the database for the last payment type ID that was transferred to the Aryan system.
// If new payment types are found (payment types with an ID greater than the last processed ID), it sends them to the Aryan system using the PostBaseDataToDeliverCenterSaleSelect method.
// The function returns an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP  200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncCustomersWithSaleCustomer(customers []models.Customers) error {

	lastCustomerId := customers[len(customers)-1].ID
	lastSaleCustomerId, err := f.repos.Database.GetCustomerToSaleCustomer()
	if err != nil {
		return err
	}

	if lastCustomerId > lastSaleCustomerId {
		customers = customers[lastSaleCustomerId:]
		res, err := f.aryan.PostCustomerToSaleCustomer(customers)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertCustomerToSaleCustomer(lastCustomerId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncProductsWithGoods retrieves all product data from the Fararavand ERP system and filters them based on the last processed product ID.
// It fetches the products using the Fararavand API, then checks the database for the last product ID that was transferred to the Aryan system.
// If new products are found (products with an ID greater than the last processed ID), it sends them to the Aryan system using the PostProductsToGoods method.
// The function returns an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP  200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncProductsWithGoods(products []models.Products) error {

	lastProductId := products[len(products)-1].ID

	lastGoodsId, err := f.repos.Database.GetProductsToGoods()
	if err != nil {
		return err
	}

	if lastProductId > lastGoodsId {
		products = products[lastGoodsId:]
		res, err := f.aryan.PostProductsToGoods(products)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertProductsToGoods(lastProductId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncInvoicesWithSaleFactor retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInoviceToSaleFactor method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleFactor(invoices []models.Invoices) error {

	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSaleFactorId, err := f.repos.Database.GetInvoiceToSaleFactor()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSaleFactorId {
		invoices = invoices[lastSaleFactorId:]
		res, err := f.aryan.PostInoviceToSaleFactor(invoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleFactor(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncInvoicesWithSaleOrder retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSaleOrder method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleOrder(invoices []models.Invoices) error {

	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSaleOrderId, err := f.repos.Database.GetInvoiceToSaleOrder()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSaleOrderId {
		invoices = invoices[lastSaleOrderId:]
		res, err := f.aryan.PostInvoiceToSaleOrder(invoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleOrder(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncInvoicesWithSalePayment retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSalePayment method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSalePayment(invoices []models.Invoices) error {

	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSalePaymentId, err := f.repos.Database.GetInvoiceToSalePayment()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalePaymentId {
		invoices = invoices[lastSalePaymentId:]
		res, err := f.aryan.PostInvoiceToSalePayment(invoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSalePayment(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncInvoicesWithSalerSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSalerSelect method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSalerSelect(invoices []models.Invoices) error {

	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSalerSelectId, err := f.repos.Database.GetInvoiceToSalerSelect()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		invoices = invoices[lastSalerSelectId:]
		res, err := f.aryan.PostInvoiceToSalerSelect(invoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSalerSelect(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncInvoicesWithSaleProforma retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSaleProforma method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleProforma(invoices []models.Invoices) error {

	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSaleProformaId, err := f.repos.Database.GetInvoiceToSaleProforma()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSaleProformaId {
		invoices = invoices[lastSaleProformaId:]
		res, err := f.aryan.PostInvoiceToSaleProforma(invoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleProforma(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncInvoicesWithSaleCenter retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSaleProforma method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleCenter(invoices []models.Invoices) error {

	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSaleProformaId, err := f.repos.Database.GetInvoiceToSaleCenter()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSaleProformaId {
		invoices = invoices[lastSaleProformaId:]
		res, err := f.aryan.PostInvoiceToSaleCenter(invoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleCenter(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncInvoiceWithSaleTypeSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSalerSelect method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoiceWithSaleTypeSelect(invoices []models.Invoices) error {

	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSalerSelectId, err := f.repos.Database.GetInvoiceToSaleTypeSelect()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		invoices = invoices[lastSalerSelectId:]
		res, err := f.aryan.PostInvoiceToSaleTypeSelect(invoices)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertInvoiceToSaleTypeSelect(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// SyncBaseDataWithDeliverCenter retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInvoiceToSalerSelect method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncBaseDataWithDeliverCenter(baseData models.BaseData) error {

	paymentType := baseData.PaymentTypes

	lastInvoiceId := paymentType[len(paymentType)-1].ID

	lastSalerSelectId, err := f.repos.Database.GetBaseDataToDeliverCenter()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		paymentType = paymentType[lastSalerSelectId:]
		baseData := models.BaseData{
			PaymentTypes: paymentType,
		}
		res, err := f.aryan.PostBaseDataToDeliverCenterSaleSelect(baseData)
		if res.StatusCode() == http.StatusOK {
			err = f.repos.Database.InsertBaseDataToDeliverCenter(lastInvoiceId)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

// GetTreasuries get all treasuries data from the first ERP
func (f *Fararavand) SyncTreasuries(treasuries []models.Treasuries) error {

	lastInvoiceId := treasuries[len(treasuries)-1].InvoiceID

	lastSalerSelectId, err := f.repos.Database.GetTreasuries()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		treasuries = treasuries[lastSalerSelectId:]
		// res, err := f.aryan.PostTreasuries(treasuries)
		// if res.StatusCode() == http.StatusOK {
		// 	err = f.repos.Database.InsertTreasuries(lastInvoiceId)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		return err
	}

	return nil
}

// GetInvoiceReturns get all revert invoices data from the first ERP
func (f *Fararavand) SyncInvoiceReturns(invoiceReturn []models.InvoiceReturn) error {

	lastInvoiceId := invoiceReturn[len(invoiceReturn)-1].InvoiceID

	lastSalerSelectId, err := f.repos.Database.GetInvoiceReturn()
	if err != nil {
		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		invoiceReturn = invoiceReturn[lastSalerSelectId:]
		// res, err := f.aryan.PostinvoiceReturn(invoiceReturn)
		// if res.StatusCode() == http.StatusOK {
		// 	err = f.repos.Database.InsertinvoiceReturn(lastInvoiceId)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		return err
	}

	return nil
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