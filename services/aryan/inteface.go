package aryan

import (
	"erp-job/models"

	"github.com/go-resty/resty/v2"
)

type AryanInterface interface {
	PostInoviceToSaleFactor(fp []models.Invoices) (*resty.Response, error)
	PostProductsToGoods(fp []models.Products) error
	PostCustomerToSaleCustomer(fc []models.Customers) (*resty.Response, error)
	PostInvoiceToSaleOrder(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleCenter(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSalePayment(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSalerSelect(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleProforma(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleTypeSelect(fp []models.Invoices) (*resty.Response, error)
	PostBaseDataToSaleCenterSelect(baseData models.BaseData) (*resty.Response, error)
	PostBaseDataToDeliverCenterSaleSelect(baseData models.BaseData) (*resty.Response, error)
	PostBaseDataToSaleSellerVisitor(baseData models.BaseData) (*resty.Response, error)
}
