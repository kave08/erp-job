package fararavand

type FararavandInterface interface {
	SyncBaseData() error
	SyncCustomersWithSaleCustomer() error
	SyncProductsWithGoods() error
	SyncInvoicesWithSaleFactor() error
	SyncInvoicesWithSaleOrder() error
	SyncInvoicesWithSalePayment() error
	SyncInvoicesWithSalerSelect() error
	SyncInvoicesWithSaleProforma() error
	SyncTreasuries() error
	SyncInvoiceReturns() error
}
