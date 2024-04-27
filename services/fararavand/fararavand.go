package fararavand

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/utility/logger"

	"go.uber.org/zap"
)

type Fararavand struct {
	log     *zap.SugaredLogger
	baseUrl string
	repos   *repository.Repository
	aryan   aryan.AryanInterface
}

func NewFararavand(repos *repository.Repository, aryan aryan.AryanInterface) FararavandInterface {
	return &Fararavand{
		log:     logger.Logger(),
		baseUrl: config.Cfg.FararavandApp.BaseURL,
		repos:   repos,
		aryan:   aryan,
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
		f.log.Errorw("GetCustomerToSaleCustomer encountered an error: ",
			"error", err,
			"last_sale_customer_id", lastSaleCustomerId,
		)

		return err
	}

	if lastCustomerId > lastSaleCustomerId {
		for index, customer := range customers {
			if customer.ID > lastSaleCustomerId {
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

	err = f.repos.Database.InsertCustomerToSaleCustomer(lastCustomerId)
	if err != nil {
		f.log.Errorw("InsertCustomerToSaleCustomer encountered an error: ",
			"error", err,
			"last_customer_id", lastCustomerId,
		)

		return err
	}

	return err
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
		f.log.Errorw("GetProductsToGoods encountered an error: ",
			"error", err,
			"last_goods_id", lastGoodsId,
		)

		return err
	}

	if lastProductId > lastGoodsId {
		for index, product := range products {
			if product.ID > lastGoodsId {
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

	err = f.repos.Database.InsertProductsToGoods(lastProductId)
	if err != nil {
		f.log.Errorw("InsertProductsToGoods encountered an error: ",
			"error", err,
			"last_product_id", lastProductId,
		)

		return err
	}

	return err
}

// SyncInvoicesWithSaleFactor retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
// It fetches the invoices using the Fararavand API, then checks the database for the last invoice ID that was transferred to the Aryan system.
// If new invoices are found (invoices with an ID greater than the last processed ID), it sends them to the Aryan system using the PostInoviceToSaleFactor method.
// The function returns a slice of new invoices and an error if any occurs during the process.
// If the response status code from the Fararavand API is not HTTP 200 OK, it logs the status code and returns an error.
func (f *Fararavand) SyncInvoicesWithSaleFactor(invoices []models.Invoices) error {
	//TODO:check max id
	lastInvoiceId := invoices[len(invoices)-1].InvoiceId

	lastSaleFactorId, err := f.repos.Database.GetInvoiceToSaleFactor()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_factor_id", lastSaleFactorId,
		)

		return err
	}

	if lastInvoiceId > lastSaleFactorId {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceId {
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

	err = f.repos.Database.InsertInvoiceToSaleFactor(lastInvoiceId)
	if err != nil {
		f.log.Errorw("InsertInvoiceToSaleFactor to encountered an error: ",
		"type", "database",
		"error", err,
		"last_invoice_id", lastInvoiceId,
	)
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
		f.log.Errorw("GetInvoiceToSaleOrder encountered an error: ",
			"error", err,
			"last_sale_order_id", lastSaleOrderId,
		)

		return err
	}

	if lastInvoiceId > lastSaleOrderId {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceId {
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

	err = f.repos.Database.InsertInvoiceToSaleOrder(lastInvoiceId)
	if err != nil {
		f.log.Errorw("InsertInvoiceToSaleOrder encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceId,
		)

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
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_payment_id", lastSalePaymentId,
		)

		return err
	}

	if lastInvoiceId > lastSalePaymentId {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceId {
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

		err = f.repos.Database.InsertInvoiceToSalePayment(lastInvoiceId)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSalePayment encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceId,
			)

			return err
		}
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
		f.log.Errorw("GetInvoiceToSalerSelect encountered an error: ",
			"error", err,
			"last_saler_select_id", lastSalerSelectId,
		)

		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceId {
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

		err = f.repos.Database.InsertInvoiceToSalerSelect(lastInvoiceId)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSalerSelect encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceId,
			)

			return err
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
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_proforma_id", lastSaleProformaId,
		)

		return err
	}

	if lastInvoiceId > lastSaleProformaId {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceId {
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

		err = f.repos.Database.InsertInvoiceToSaleProforma(lastInvoiceId)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSaleProforma encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceId,
			)

			return err
		}

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
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_proforma_id", lastSaleProformaId,
		)

		return err
	}

	if lastInvoiceId > lastSaleProformaId {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceId {
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

		err = f.repos.Database.InsertInvoiceToSaleCenter(lastInvoiceId)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSaleCenter encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceId,
			)

			return err
		}
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
		f.log.Errorw("GetInvoiceToSaleTypeSelect encountered an error: ",
			"error", err,
			"last_saler_select_id", lastSalerSelectId,
		)

		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		for index, invoice := range invoices {
			if invoice.InvoiceId > lastInvoiceId {
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

		err = f.repos.Database.InsertInvoiceToSaleTypeSelect(lastInvoiceId)
		if err != nil {
			f.log.Errorw("InsertInvoiceToSaleTypeSelect encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceId,
			)

			return err
		}

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
		f.log.Errorw("GetBaseDataToDeliverCenter encountered an error: ",
			"error", err,
			"last_saler_select_id", lastSalerSelectId,
		)
		return err
	}

	if lastInvoiceId > lastSalerSelectId {
		for index, invoice := range paymentType {
			if invoice.ID > lastInvoiceId {
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

		err = f.repos.Database.InsertBaseDataToDeliverCenter(lastInvoiceId)
		if err != nil {
			f.log.Errorw("InsertBaseDataToDeliverCenter encountered an error: ",
				"error", err,
				"last_invoice_id", lastInvoiceId,
			)

			return err
		}

	}

	return nil
}

// GetTreasuries get all treasuries data from the first ERP
func (f *Fararavand) SyncTreasuries(treasuries []models.Treasuries) error {

	return nil
}

// GetInvoiceReturns get all revert invoices data from the first ERP
func (f *Fararavand) SyncInvoiceReturns(invoiceReturn []models.InvoiceReturn) error {

	return nil
}

// SyncBaseData gets all base information from the first ERP
func (f *Fararavand) SyncBaseData() error {

	return nil
}
