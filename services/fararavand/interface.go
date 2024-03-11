package fararavand

import "erp-job/models"

type FararavandInterface interface {
	SyncBaseData() error
	SyncCustomersWithSaleCustomer(customers []models.Customers) error
	SyncProductsWithGoods(products []models.Products) error
	SyncInvoicesWithSaleFactor(invoices []models.Invoices) error
	SyncInvoicesWithSaleOrder(invoices []models.Invoices) error
	SyncInvoicesWithSalePayment(invoices []models.Invoices) error
	SyncInvoicesWithSalerSelect(invoices []models.Invoices) error
	SyncInvoicesWithSaleProforma(invoices []models.Invoices) error
	SyncInvoicesWithSaleCenter(invoices []models.Invoices) error
	SyncInvoiceWithSaleTypeSelect(invoices []models.Invoices) error
	SyncTreasuries(treasuries []models.Treasuries) error
	SyncInvoiceReturns(invoiceReturn []models.InvoiceReturn) error
	SyncBaseDataWithDeliverCenter(baseData models.BaseData) error
}
