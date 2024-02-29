package aryan

import (
	"erp-job/models"

	"github.com/go-resty/resty/v2"
)

type AryanInterface interface {
	PostInoviceToSaleFactor(fp []models.Invoices) error
	PostProductsToGoods(fp []models.Products) error
	PostCustomerToSaleCustomer(fc []models.Customers) error
	PostInvoiceToSaleOrder(fp []models.Invoices) error
	PostInvoiceToSalePayment(fp []models.Invoices) error
	PostInvoiceToSaleCenter(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSalerSelect(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleProforma(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleTypeSelect(fp []models.Invoices) (*resty.Response, error)
	PostBaseDataToSaleCenterSelect(baseData models.BaseData) (*resty.Response, error)
	PostBaseDataToDeliverCenterSaleSelect(baseData models.BaseData) (*resty.Response, error)
	PostBaseDataToSaleSellerVisitor(baseData models.BaseData) (*resty.Response, error)
}
