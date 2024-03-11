package aryan

import (
	"erp-job/models"
)

type AryanInterface interface {
	PostInvoiceToSaleFactor(fp []models.Invoices) error
	PostProductsToGoods(fp []models.Products) error
	PostCustomerToSaleCustomer(fc []models.Customers) error
	PostInvoiceToSaleOrder(fp []models.Invoices) error
	PostInvoiceToSalePayment(fp []models.Invoices) error
	PostInvoiceToSaleCenter(fp []models.Invoices) error
	PostInvoiceToSalerSelect(fp []models.Invoices) error
	PostInvoiceToSaleProforma(fp []models.Invoices) error
	PostInvoiceToSaleTypeSelect(fp []models.Invoices) error
	PostBaseDataToSaleCenterSelect(baseData models.BaseData) error
	PostBaseDataToDeliverCenterSaleSelect(baseData models.BaseData) error
	PostBaseDataToSaleSellerVisitor(baseData models.BaseData) error
}
