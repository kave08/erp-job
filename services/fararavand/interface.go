package fararavand

import "erp-job/models"

type FararavandInterface interface {
	SyncBaseData() error
	SyncCustomersWithSaleCustomer(customers []models.Customers) error
	SyncProductsWithGoods(newProducts []models.Products) error
	SyncInvoicesWithSaleFactor(invoices []models.Invoices) error
	SyncInvoicesWithSaleOrder(invoices []models.Invoices) error
	SyncInvoicesWithSalePayment(invoices []models.Invoices) error
	SyncInvoicesWithSalerSelect(invoices []models.Invoices) error
	SyncInvoicesWithSaleProforma(invoices []models.Invoices) error
	SyncTreasuries() error
	SyncInvoiceReturns() error
}
