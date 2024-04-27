package fararavand

import "erp-job/models"

// FararavandInterface defines the contract for interacting with the Fararavand ERP system, including synchronizing customer, product, and invoice data.
//
// It abstracts the operations for data synchronization between Fararavand and Aryan systems, ensuring data consistency and integrity across both platforms.
type FararavandInterface interface {
	// SyncCustomersWithSaleCustomer synchronizes customer data from Fararavand to Aryan by filtering based on the last processed customer ID.
	//
	// It updates the database with the latest customer ID processed and logs any errors encountered during the process.
	SyncCustomersWithSaleCustomer(customers []models.Customers) error

	// SyncProductsWithGoods retrieves all product data from the Fararavand ERP system and filters them based on the last processed product ID.
	//
	// It sends new products to the Aryan system using the PostProductsToGoods method and updates the database with the last processed product ID.
	SyncProductsWithGoods(products []models.Products) error

	// SyncInvoicesWithSaleFactor retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
	//
	// It sends new invoices to the Aryan system using the PostInvoiceToSaleFactor method and updates the database with the last processed invoice ID.
	SyncInvoicesWithSaleFactor(invoices []models.Invoices) error

	// SyncInvoicesWithSaleOrder retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
	//
	// It sends new invoices to the Aryan system using the PostInvoiceToSaleOrder method and updates the database with the last processed invoice ID.
	SyncInvoicesWithSaleOrder(invoices []models.Invoices) error

	// SyncInvoicesWithSalePayment retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
	//
	// It sends new invoices to the Aryan system using the PostInvoiceToSalePayment method and updates the database with the last processed invoice ID.
	SyncInvoicesWithSalePayment(invoices []models.Invoices) error

	// SyncInvoicesWithSalerSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
	//
	// It sends new invoices to the Aryan system using the PostInvoiceToSalerSelect method and updates the database with the last processed invoice ID.
	SyncInvoicesWithSalerSelect(invoices []models.Invoices) error

	// SyncInvoicesWithSaleProforma retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
	//
	// It sends new invoices to the Aryan system using the PostInvoiceToSaleProforma method and updates the database with the last processed invoice ID.
	SyncInvoicesWithSaleProforma(invoices []models.Invoices) error

	// SyncInvoicesWithSaleCenter retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
	//
	// It sends new invoices to the Aryan system using the PostInvoiceToSaleCenter method and updates the database with the last processed invoice ID.
	SyncInvoicesWithSaleCenter(invoices []models.Invoices) error

	// SyncInvoiceWithSaleTypeSelect retrieves all invoices from the Fararavand ERP system and filters them based on the last processed invoice ID.
	//
	// It sends new invoices to the Aryan system using the PostInvoiceToSaleTypeSelect method and updates the database with the last processed invoice ID.
	SyncInvoiceWithSaleTypeSelect(invoices []models.Invoices) error

	// SyncBaseDataWithDeliverCenter retrieves all base data from the Fararavand ERP system and filters them based on the last processed base data ID.
	//
	// It sends new base data to the Aryan system using the PostBaseDataToDeliverCenterSaleSelect method and updates the database with the last processed base data ID.
	SyncBaseDataWithDeliverCenter(baseData models.BaseData) error

	// SyncTreasuries get all treasuries data from the first ERP
	SyncTreasuries(treasuries []models.Treasuries) error

	// SyncInvoiceReturns get all revert invoices data from the first ERP
	SyncInvoiceReturns(invoiceReturn []models.InvoiceReturn) error

	// SyncBaseData gets all base information from the first ERP
	SyncBaseData() error
}
