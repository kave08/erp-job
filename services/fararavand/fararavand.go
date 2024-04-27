package fararavand

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/utility/logger"

	"go.uber.org/zap"
)

// Fararavand represents the service for interacting with the Fararavand ERP system, encapsulating logging, base URL, repository, and Aryan interface for data synchronization.
type Fararavand struct {
	log     *zap.SugaredLogger
	baseURL string
	repos   *repository.Repository
	aryan   aryan.AryanInterface
}

// NewFararavand initializes and returns a new Fararavand service instance.
func NewFararavand(repos *repository.Repository, aryan aryan.AryanInterface) Interface {
	return &Fararavand{
		log:     logger.Logger(),
		baseURL: config.Cfg.FararavandApp.BaseURL,
		repos:   repos,
		aryan:   aryan,
	}
}

// SyncCustomersWithSaleCustomer synchronizes customer data from Fararavand to Aryan by filtering based on the last processed customer ID.
//
// It updates the database with the latest customer ID processed and logs any errors encountered during the process.
func (f *Fararavand) SyncCustomersWithSaleCustomer(customers []models.Customers) error {

	lastCustomerID := customers[len(customers)-1].ID

	lastSaleCustomerID, err := f.repos.Database.GetCustomerToSaleCustomer()
	if err != nil {
		f.log.Errorw("GetCustomerToSaleCustomer encountered an error: ",
			"error", err,
			"last_sale_customer_id", lastSaleCustomerID,
		)

		return err
	}

	if lastCustomerID > lastSaleCustomerID {
		for index, customer := range customers {
			if customer.ID > lastSaleCustomerID {
				customers = customers[index:]
				break
			}
		}
	}

	err = f.aryan.PostCustomerToSaleCustomer(customers)
	if err != nil {
		f.log.Errorw("PostCustomerToSaleCustomer encountered an error: ",
			"error", err,
		)

		return err
	}

	err = f.repos.Database.InsertCustomerToSaleCustomer(lastCustomerID)
	if err != nil {
		f.log.Errorw("InsertCustomerToSaleCustomer encountered an error: ",
			"error", err,
			"last_customer_id", lastCustomerID,
		)

		return err
	}

	return err
}

// SyncProductsWithGoods retrieves all product data from the Fararavand ERP system and filters them based on the last processed product ID.
//
// It sends new products to the Aryan system using the PostProductsToGoods method and updates the database with the last processed product ID.
func (f *Fararavand) SyncProductsWithGoods(products []models.Products) error {

	lastProductID := products[len(products)-1].ID

	lastGoodsID, err := f.repos.Database.GetProductsToGoods()
	if err != nil {
		f.log.Errorw("GetProductsToGoods encountered an error: ",
			"error", err,
			"last_goods_id", lastGoodsID,
		)

		return err
	}

	if lastProductID > lastGoodsID {
		for index, product := range products {
			if product.ID > lastGoodsID {
				products = products[index:]
				break
			}
		}
	}

	err = f.aryan.PostProductsToGoods(products)
	if err != nil {
		f.log.Errorw("PostProductsToGoods encountered an error: ",
			"error", err,
		)

		return err
	}

	err = f.repos.Database.InsertProductsToGoods(lastProductID)
	if err != nil {
		f.log.Errorw("InsertProductsToGoods encountered an error: ",
			"error", err,
			"last_product_id", lastProductID,
		)

		return err
	}

	return err
}

// SyncInvoicesWithSaleFactor retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
//
// It sends new invoices to the Aryan system using the PostInvoiceToSaleFactor method and updates the database with the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleFactor(invoices []models.Invoices) error {

	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	lastSaleFactorID, err := f.repos.Database.GetInvoiceToSaleFactor()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_factor_id", lastSaleFactorID,
		)

		return err
	}

	if lastInvoiceID > lastSaleFactorID {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceID {
				invoices = invoices[index:]
				break
			}
		}
	}

	err = f.aryan.PostInvoiceToSaleFactor(invoices)
	if err != nil {
		f.log.Errorw("PostInvoiceToSaleFactor encountered an error: ",
			"error", err,
		)

		return err
	}

	err = f.repos.Database.InsertInvoiceToSaleFactor(lastInvoiceID)
	if err != nil {
		f.log.Errorw("InsertInvoiceToSaleFactor to encountered an error: ",
			"type", "database",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)
		return err
	}

	return nil
}

// SyncInvoicesWithSaleOrder retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
//
// It sends new invoices to the Aryan system using the PostInvoiceToSaleOrder method and updates the database with the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleOrder(invoices []models.Invoices) error {

	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	lastSaleOrderID, err := f.repos.Database.GetInvoiceToSaleOrder()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleOrder encountered an error: ",
			"error", err,
			"last_sale_order_id", lastSaleOrderID,
		)

		return err
	}

	if lastInvoiceID > lastSaleOrderID {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceID {
				invoices = invoices[index:]
				break
			}
		}
	}

	err = f.aryan.PostInvoiceToSaleOrder(invoices)
	if err != nil {
		f.log.Errorw("PostInvoiceToSaleOrder encountered an error: ",
			"error", err,
		)

		return err

	}

	err = f.repos.Database.InsertInvoiceToSaleOrder(lastInvoiceID)
	if err != nil {
		f.log.Errorw("InsertInvoiceToSaleOrder encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncInvoicesWithSalePayment retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
//
// It sends new invoices to the Aryan system using the PostInvoiceToSalePayment method and updates the database with the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSalePayment(invoices []models.Invoices) error {

	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	lastSalePaymentID, err := f.repos.Database.GetInvoiceToSalePayment()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_payment_id", lastSalePaymentID,
		)

		return err
	}

	if lastInvoiceID > lastSalePaymentID {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceID {
				invoices = invoices[index:]
				break
			}
		}

		err = f.aryan.PostInvoiceToSalePayment(invoices)
		if err != nil {
			f.log.Errorw("PostInvoiceToSalePayment encountered an error: ",
				"error", err,
			)

			return err
		}

		err = f.repos.Database.InsertInvoiceToSalePayment(lastInvoiceID)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSalePayment encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceID,
			)

			return err
		}
	}

	return nil
}

// SyncInvoicesWithSalerSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
//
// It sends new invoices to the Aryan system using the PostInvoiceToSalerSelect method and updates the database with the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSalerSelect(invoices []models.Invoices) error {

	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	lastSalerSelectID, err := f.repos.Database.GetInvoiceToSalerSelect()
	if err != nil {
		f.log.Errorw("GetInvoiceToSalerSelect encountered an error: ",
			"error", err,
			"last_saler_select_id", lastSalerSelectID,
		)

		return err
	}

	if lastInvoiceID > lastSalerSelectID {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceID {
				invoices = invoices[index:]
				break
			}
		}

		err := f.aryan.PostInvoiceToSalerSelect(invoices)
		if err != nil {
			f.log.Errorw("PostInvoiceToSalerSelect encountered an error: ",
				"error", err,
			)

			return err
		}

		err = f.repos.Database.InsertInvoiceToSalerSelect(lastInvoiceID)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSalerSelect encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceID,
			)

			return err
		}

		return err
	}

	return nil
}

// SyncInvoicesWithSaleProforma retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
//
// It sends new invoices to the Aryan system using the PostInvoiceToSaleProforma method and updates the database with the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleProforma(invoices []models.Invoices) error {

	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	lastSaleProformaID, err := f.repos.Database.GetInvoiceToSaleProforma()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_proforma_id", lastSaleProformaID,
		)

		return err
	}

	if lastInvoiceID > lastSaleProformaID {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceID {
				invoices = invoices[index:]
				break
			}
		}

		err := f.aryan.PostInvoiceToSaleProforma(invoices)
		if err != nil {
			f.log.Errorw("PostInvoiceToSaleProforma encountered an error: ",
				"error", err,
			)

			return err
		}

		err = f.repos.Database.InsertInvoiceToSaleProforma(lastInvoiceID)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSaleProforma encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceID,
			)

			return err
		}

	}

	return nil
}

// SyncInvoicesWithSaleCenter retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
//
// It sends new invoices to the Aryan system using the PostInvoiceToSaleCenter method and updates the database with the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleCenter(invoices []models.Invoices) error {

	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	lastSaleProformaID, err := f.repos.Database.GetInvoiceToSaleCenter()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_proforma_id", lastSaleProformaID,
		)

		return err
	}

	if lastInvoiceID > lastSaleProformaID {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceID {
				invoices = invoices[index:]
				break
			}
		}

		err := f.aryan.PostInvoiceToSaleCenter(invoices)
		if err != nil {
			f.log.Errorw("PostInvoiceToSaleCenter encountered an error: ",
				"error", err,
			)

			return err
		}

		err = f.repos.Database.InsertInvoiceToSaleCenter(lastInvoiceID)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSaleCenter encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceID,
			)

			return err
		}
	}

	return nil
}

// SyncInvoiceWithSaleTypeSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
//
// It sends new invoices to the Aryan system using the PostInvoiceToSaleTypeSelect method and updates the database with the last processed invoice ID.
func (f *Fararavand) SyncInvoiceWithSaleTypeSelect(invoices []models.Invoices) error {

	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	lastSalerSelectID, err := f.repos.Database.GetInvoiceToSaleTypeSelect()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleTypeSelect encountered an error: ",
			"error", err,
			"last_saler_select_id", lastSalerSelectID,
		)

		return err
	}

	if lastInvoiceID > lastSalerSelectID {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceID {
				invoices = invoices[index:]
				break
			}
		}

		err := f.aryan.PostInvoiceToSaleTypeSelect(invoices)
		if err != nil {
			f.log.Errorw("PostInvoiceToSaleTypeSelect encountered an error: ",
				"error", err,
			)

			return err
		}

		err = f.repos.Database.InsertInvoiceToSaleTypeSelect(lastInvoiceID)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSaleTypeSelect encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceID,
			)

			return err
		}

	}

	return nil
}

// SyncBaseDataWithDeliverCenter retrieves all base data from the Fararavand ERP system and filters them based on the last processed base data ID.
//
// It sends new base data to the Aryan system using the PostBaseDataToDeliverCenterSaleSelect method and updates the database with the last processed base data ID.
func (f *Fararavand) SyncBaseDataWithDeliverCenter(baseData models.BaseData) error {

	paymentType := baseData.PaymentTypes
	lastInvoiceID := paymentType[len(paymentType)-1].ID

	lastSalerSelectID, err := f.repos.Database.GetBaseDataToDeliverCenter()
	if err != nil {
		f.log.Errorw("GetBaseDataToDeliverCenter encountered an error: ",
			"error", err,
			"last_saler_select_id", lastSalerSelectID,
		)
		return err
	}

	if lastInvoiceID > lastSalerSelectID {
		for index, invoice := range paymentType {
			if invoice.ID > lastInvoiceID {
				paymentType = paymentType[index:]
				break
			}
		}
		baseData := models.BaseData{
			PaymentTypes: paymentType,
		}

		err := f.aryan.PostBaseDataToDeliverCenterSaleSelect(baseData)
		if err != nil {
			f.log.Errorw("PostBaseDataToDeliverCenterSaleSelect encountered an error: ",
				"error", err,
			)

			return err
		}

		err = f.repos.Database.InsertBaseDataToDeliverCenter(lastInvoiceID)
		if err != nil {
			f.log.Errorw("InsertBaseDataToDeliverCenter encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceID,
			)

			return err
		}

	}

	return nil
}

// SyncTreasuries get all treasuries data from the first ERP
func (f *Fararavand) SyncTreasuries(treasuries []models.Treasuries) error {

	return nil
}

// SyncInvoiceReturns get all revert invoices data from the first ERP
func (f *Fararavand) SyncInvoiceReturns(invoiceReturn []models.InvoiceReturn) error {

	return nil
}

// SyncBaseData gets all base information from the first ERP
func (f *Fararavand) SyncBaseData() error {

	return nil
}
