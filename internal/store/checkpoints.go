package store

type Progress struct {
	LastID     int
	PageNumber int
}

type CheckpointStore interface {
	GetInvoiceProgress() (Progress, error)
	SaveInvoiceProgress(progress Progress) error
	GetCustomerProgress() (Progress, error)
	SaveCustomerProgress(progress Progress) error
	GetProductProgress() (Progress, error)
	SaveProductProgress(progress Progress) error
	GetBaseDataProgress() (Progress, error)
	SaveBaseDataProgress(progress Progress) error

	GetProductsToGoods() (int, error)
	SaveProductsToGoods(id int) error
	GetCustomerToSaleCustomer() (int, error)
	SaveCustomerToSaleCustomer(id int) error
	GetInvoiceToSaleFactor() (int, error)
	SaveInvoiceToSaleFactor(id int) error
	GetInvoiceToSaleOrder() (int, error)
	SaveInvoiceToSaleOrder(id int) error
	GetInvoiceToSalePayment() (int, error)
	SaveInvoiceToSalePayment(id int) error
	GetInvoiceToSalerSelect() (int, error)
	SaveInvoiceToSalerSelect(id int) error
	GetInvoiceToSaleProforma() (int, error)
	SaveInvoiceToSaleProforma(id int) error
	GetInvoiceToSaleCenter() (int, error)
	SaveInvoiceToSaleCenter(id int) error
	GetInvoiceToSaleTypeSelect() (int, error)
	SaveInvoiceToSaleTypeSelect(id int) error
	GetBaseDataToDeliverCenter() (int, error)
	SaveBaseDataToDeliverCenter(id int) error
}
