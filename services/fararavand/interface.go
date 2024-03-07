package fararavand

import "erp-job/models"

type FararavandInterface interface {
	SyncBaseData() error
	SyncCustomersWithSaleCustomer(customers []models.Customers) (int, error)
	SyncProductsWithGoods(products []models.Products) (int, error)
	SyncInvoicesWithSaleFactor(invoices []models.Invoices) (int, error)
	SyncInvoicesWithSaleOrder(invoices []models.Invoices) (int, error)
	SyncInvoicesWithSalePayment(invoices []models.Invoices) (int, error)
	SyncInvoicesWithSalerSelect(invoices []models.Invoices) (int, error)
	SyncInvoicesWithSaleProforma(invoices []models.Invoices) error
	SyncInvoicesWithSaleCenter(invoices []models.Invoices) error
	SyncInvoiceWithSaleTypeSelect(invoices []models.Invoices) error
	SyncTreasuries(treasuries []models.Treasuries) error
	SyncInvoiceReturns(invoiceReturn []models.InvoiceReturn) error
	SyncBaseDataWithDeliverCenter(baseData models.BaseData) error
}
