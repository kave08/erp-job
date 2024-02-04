package fararavand

import "erp-job/models"

type FararavandInterface interface {
	GetCustomersForSaleCustomer() error
	GetBaseData() (*models.Fararavand, error)
	GetInvoicesForSaleFactor() ([]models.Invoices, error)
	GetInvoicesForSaleOrder() ([]models.Invoices, error)
	GetInvoicesForSalePayment() ([]models.Invoices, error)
	GetInvoicesForSalerSelect() ([]models.Invoices, error)
	GetInvoicesForSaleProforma() ([]models.Invoices, error)
	GetProductsToGoods() ([]models.Products, error)
	GetTreasuries() ([]models.Fararavand, error)
	GetInvoiceReturns() ([]models.Fararavand, error)
}
