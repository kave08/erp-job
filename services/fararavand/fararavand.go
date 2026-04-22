package fararavand

import (
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/utility/logger"
	"errors"

	"go.uber.org/zap"
)

var (
	errTreasuriesSyncNotImplemented     = errors.New("treasury sync is not implemented: missing Aryan target contract")
	errInvoiceReturnsSyncNotImplemented = errors.New("invoice return sync is not implemented: missing Aryan target contract")
	errDirectBaseDataSyncNotImplemented = errors.New("direct base-data sync entrypoint is not implemented")
)

// Fararavand represents the service for interacting with the Fararavand ERP system, encapsulating logging, repository, and Aryan interface for data synchronization.
type Fararavand struct {
	log   *zap.SugaredLogger
	repos *repository.Repository
	aryan aryan.AryanInterface
}

// NewFararavand initializes and returns a new Fararavand service instance.
func NewFararavand(repos *repository.Repository, aryan aryan.AryanInterface) Interface {
	return &Fararavand{
		log:   logger.Logger(),
		repos: repos,
		aryan: aryan,
	}
}

func firstUnsyncedIndex(length int, lastSyncedID int, idAt func(int) int) int {
	if length == 0 || idAt(length-1) <= lastSyncedID {
		return -1
	}

	for index := 0; index < length; index++ {
		if idAt(index) > lastSyncedID {
			return index
		}
	}

	return -1
}

// SyncCustomersWithSaleCustomer synchronizes customer data from Fararavand to Aryan by filtering based on the last processed customer ID.
func (f *Fararavand) SyncCustomersWithSaleCustomer(customers []models.Customers) error {
	lastSaleCustomerID, err := f.repos.Database.GetCustomerToSaleCustomer()
	if err != nil {
		f.log.Errorw("GetCustomerToSaleCustomer encountered an error: ",
			"error", err,
			"last_sale_customer_id", lastSaleCustomerID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(customers), lastSaleCustomerID, func(i int) int {
		return customers[i].ID
	})
	if index < 0 {
		return nil
	}

	customers = customers[index:]
	lastCustomerID := customers[len(customers)-1].ID

	if err := f.aryan.PostCustomerToSaleCustomer(customers); err != nil {
		f.log.Errorw("PostCustomerToSaleCustomer encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertCustomerToSaleCustomer(lastCustomerID); err != nil {
		f.log.Errorw("InsertCustomerToSaleCustomer encountered an error: ",
			"error", err,
			"last_customer_id", lastCustomerID,
		)

		return err
	}

	return nil
}

// SyncProductsWithGoods retrieves all product data from the Fararavand ERP system and filters them based on the last processed product ID.
func (f *Fararavand) SyncProductsWithGoods(products []models.Products) error {
	lastGoodsID, err := f.repos.Database.GetProductsToGoods()
	if err != nil {
		f.log.Errorw("GetProductsToGoods encountered an error: ",
			"error", err,
			"last_goods_id", lastGoodsID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(products), lastGoodsID, func(i int) int {
		return products[i].ID
	})
	if index < 0 {
		return nil
	}

	products = products[index:]
	lastProductID := products[len(products)-1].ID

	if err := f.aryan.PostProductsToGoods(products); err != nil {
		f.log.Errorw("PostProductsToGoods encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertProductsToGoods(lastProductID); err != nil {
		f.log.Errorw("InsertProductsToGoods encountered an error: ",
			"error", err,
			"last_product_id", lastProductID,
		)

		return err
	}

	return nil
}

// SyncInvoicesWithSaleFactor retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleFactor(invoices []models.Invoices) error {
	lastSaleFactorID, err := f.repos.Database.GetInvoiceToSaleFactor()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleFactor encountered an error: ",
			"error", err,
			"last_sale_factor_id", lastSaleFactorID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(invoices), lastSaleFactorID, func(i int) int {
		return invoices[i].InvoiceId
	})
	if index < 0 {
		return nil
	}

	invoices = invoices[index:]
	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	if err := f.aryan.PostInvoiceToSaleFactor(invoices); err != nil {
		f.log.Errorw("PostInvoiceToSaleFactor encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertInvoiceToSaleFactor(lastInvoiceID); err != nil {
		f.log.Errorw("InsertInvoiceToSaleFactor encountered an error: ",
			"type", "database",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncInvoicesWithSaleOrder retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleOrder(invoices []models.Invoices) error {
	lastSaleOrderID, err := f.repos.Database.GetInvoiceToSaleOrder()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleOrder encountered an error: ",
			"error", err,
			"last_sale_order_id", lastSaleOrderID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(invoices), lastSaleOrderID, func(i int) int {
		return invoices[i].InvoiceId
	})
	if index < 0 {
		return nil
	}

	invoices = invoices[index:]
	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	if err := f.aryan.PostInvoiceToSaleOrder(invoices); err != nil {
		f.log.Errorw("PostInvoiceToSaleOrder encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertInvoiceToSaleOrder(lastInvoiceID); err != nil {
		f.log.Errorw("InsertInvoiceToSaleOrder encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncInvoicesWithSalePayment retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSalePayment(invoices []models.Invoices) error {
	lastSalePaymentID, err := f.repos.Database.GetInvoiceToSalePayment()
	if err != nil {
		f.log.Errorw("GetInvoiceToSalePayment encountered an error: ",
			"error", err,
			"last_sale_payment_id", lastSalePaymentID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(invoices), lastSalePaymentID, func(i int) int {
		return invoices[i].InvoiceId
	})
	if index < 0 {
		return nil
	}

	invoices = invoices[index:]
	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	if err := f.aryan.PostInvoiceToSalePayment(invoices); err != nil {
		f.log.Errorw("PostInvoiceToSalePayment encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertInvoiceToSalePayment(lastInvoiceID); err != nil {
		f.log.Errorw("InsertInvoiceToSalePayment encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncInvoicesWithSalerSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSalerSelect(invoices []models.Invoices) error {
	lastSalerSelectID, err := f.repos.Database.GetInvoiceToSalerSelect()
	if err != nil {
		f.log.Errorw("GetInvoiceToSalerSelect encountered an error: ",
			"error", err,
			"last_saler_select_id", lastSalerSelectID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(invoices), lastSalerSelectID, func(i int) int {
		return invoices[i].InvoiceId
	})
	if index < 0 {
		return nil
	}

	invoices = invoices[index:]
	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	if err := f.aryan.PostInvoiceToSalerSelect(invoices); err != nil {
		f.log.Errorw("PostInvoiceToSalerSelect encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertInvoiceToSalerSelect(lastInvoiceID); err != nil {
		f.log.Errorw("InsertInvoiceToSalerSelect encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncInvoicesWithSaleProforma retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleProforma(invoices []models.Invoices) error {
	lastSaleProformaID, err := f.repos.Database.GetInvoiceToSaleProforma()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleProforma encountered an error: ",
			"error", err,
			"last_sale_proforma_id", lastSaleProformaID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(invoices), lastSaleProformaID, func(i int) int {
		return invoices[i].InvoiceId
	})
	if index < 0 {
		return nil
	}

	invoices = invoices[index:]
	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	if err := f.aryan.PostInvoiceToSaleProforma(invoices); err != nil {
		f.log.Errorw("PostInvoiceToSaleProforma encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertInvoiceToSaleProforma(lastInvoiceID); err != nil {
		f.log.Errorw("InsertInvoiceToSaleProforma encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncInvoicesWithSaleCenter retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
func (f *Fararavand) SyncInvoicesWithSaleCenter(invoices []models.Invoices) error {
	lastSaleCenterID, err := f.repos.Database.GetInvoiceToSaleCenter()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleCenter encountered an error: ",
			"error", err,
			"last_sale_center_id", lastSaleCenterID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(invoices), lastSaleCenterID, func(i int) int {
		return invoices[i].InvoiceId
	})
	if index < 0 {
		return nil
	}

	invoices = invoices[index:]
	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	if err := f.aryan.PostInvoiceToSaleCenter(invoices); err != nil {
		f.log.Errorw("PostInvoiceToSaleCenter encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertInvoiceToSaleCenter(lastInvoiceID); err != nil {
		f.log.Errorw("InsertInvoiceToSaleCenter encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncInvoiceWithSaleTypeSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
func (f *Fararavand) SyncInvoiceWithSaleTypeSelect(invoices []models.Invoices) error {
	lastSaleTypeSelectID, err := f.repos.Database.GetInvoiceToSaleTypeSelect()
	if err != nil {
		f.log.Errorw("GetInvoiceToSaleTypeSelect encountered an error: ",
			"error", err,
			"last_sale_type_select_id", lastSaleTypeSelectID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(invoices), lastSaleTypeSelectID, func(i int) int {
		return invoices[i].InvoiceId
	})
	if index < 0 {
		return nil
	}

	invoices = invoices[index:]
	lastInvoiceID := invoices[len(invoices)-1].InvoiceId

	if err := f.aryan.PostInvoiceToSaleTypeSelect(invoices); err != nil {
		f.log.Errorw("PostInvoiceToSaleTypeSelect encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertInvoiceToSaleTypeSelect(lastInvoiceID); err != nil {
		f.log.Errorw("InsertInvoiceToSaleTypeSelect encountered an error: ",
			"error", err,
			"last_invoice_id", lastInvoiceID,
		)

		return err
	}

	return nil
}

// SyncBaseDataWithDeliverCenter retrieves all base data from the Fararavand ERP system and filters them based on the last processed base data ID.
func (f *Fararavand) SyncBaseDataWithDeliverCenter(baseData models.BaseData) error {
	paymentTypes := baseData.PaymentTypes
	lastDeliverCenterID, err := f.repos.Database.GetBaseDataToDeliverCenter()
	if err != nil {
		f.log.Errorw("GetBaseDataToDeliverCenter encountered an error: ",
			"error", err,
			"last_deliver_center_id", lastDeliverCenterID,
		)

		return err
	}

	index := firstUnsyncedIndex(len(paymentTypes), lastDeliverCenterID, func(i int) int {
		return paymentTypes[i].ID
	})
	if index < 0 {
		return nil
	}

	paymentTypes = paymentTypes[index:]
	lastPaymentTypeID := paymentTypes[len(paymentTypes)-1].ID

	if err := f.aryan.PostBaseDataToDeliverCenterSaleSelect(models.BaseData{
		PaymentTypes: paymentTypes,
	}); err != nil {
		f.log.Errorw("PostBaseDataToDeliverCenterSaleSelect encountered an error: ",
			"error", err,
		)

		return err
	}

	if err := f.repos.Database.InsertBaseDataToDeliverCenter(lastPaymentTypeID); err != nil {
		f.log.Errorw("InsertBaseDataToDeliverCenter encountered an error: ",
			"error", err,
			"last_payment_type_id", lastPaymentTypeID,
		)

		return err
	}

	return nil
}

// SyncTreasuries get all treasuries data from the first ERP
func (f *Fararavand) SyncTreasuries(treasuries []models.Treasuries) error {
	return errTreasuriesSyncNotImplemented
}

// SyncInvoiceReturns get all revert invoices data from the first ERP
func (f *Fararavand) SyncInvoiceReturns(invoiceReturn []models.InvoiceReturn) error {
	return errInvoiceReturnsSyncNotImplemented
}

// SyncBaseData gets all base information from the first ERP
func (f *Fararavand) SyncBaseData() error {
	return errDirectBaseDataSyncNotImplemented
}
