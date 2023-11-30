package logics

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type AryanInterface interface {
	PostProductsToSaleFactor(fp []models.Products) (*resty.Response, error)
	PostSaleCustomer(f models.Fararavand) (*resty.Response, error)
}

type Aryan struct {
	restyClient *resty.Client
	baseUrl     string
	repos       *repository.Repository
}

func NewAryan(repos *repository.Repository) AryanInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.AryanApp.APIKey).SetBaseURL(config.Cfg.AryanApp.BaseURL)

	return &Aryan{
		restyClient: c,
		baseUrl:     config.Cfg.AryanApp.BaseURL,
		repos:       repos,
	}
}

// PostSalesOrder Post all sale order data to the secound ERP
func (a *Aryan) PostProductsToSaleFactor(fp []models.Products) (*resty.Response, error) {
	var newSaleFactor []models.SaleFactor

	for _, item := range fp {
		newSaleFactor = append(newSaleFactor, models.SaleFactor{
			CustomerId:     item.Customers.CustomerId,
			ServiceGoodsID: item.Codekala, // ok
			Quantity:       float64(item.Invoices.ProductCount),
			Fee:            float64(item.ProductFee),
			VoucherDesc:    "ETL-Form Fararavand",
			SecondNumber:   strconv.Itoa(item.Invoices.InvoiceId),
			// VoucherDate:      strconv.Itoa(item.Invoices.InvoiceDate),
			StockID:          10000006,
			SaleTypeId:       10000001,
			DeliveryCenterID: 10000002,
			SaleCenterID:     10000001,
			PaymentWayID:     10000001,
			SellerID:         10000002,
			SaleManID:        10000001,
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleFactor).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostSalesOrder Post all sale order data to the secound ERP
func (a *Aryan) PostCustomersToSaleFactor(fp []models.Customers) (*resty.Response, error) {
	var newSaleFactor []models.SaleFactor

	for _, item := range fp {
		newSaleFactor = append(newSaleFactor, models.SaleFactor{
			CustomerId:     item.CustomerId,
			ServiceGoodsID: 0,
			Quantity:       0,
			Fee:            0,
			VoucherDesc:    "ETL-Form Fararavand",
			SecondNumber:   strconv.Itoa(item.Invoices.InvoiceId),
			// VoucherDate:      strconv.Itoa(item.Invoices.InvoiceDate),
			StockID:          10000006,
			SaleTypeId:       10000001,
			DeliveryCenterID: 10000002,
			SaleCenterID:     10000001,
			PaymentWayID:     10000001,
			SellerID:         10000002,
			SaleManID:        10000001,
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleFactor).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostSaleCustomer Post all sale customer data to the secound ERP
func (a *Aryan) PostSaleCustomer(f models.Fararavand) (*resty.Response, error) {
	var newSaleCustomer []models.SaleCustomer

	for _, item := range fp {
		newSaleCustomer = append(newSaleCustomer, models.SaleCustomer{
			CustomerID:   item.Customers.CustomerId,
			CustomerCode: strconv.Itoa(item.Customers.CustomerCodePosti),
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleCustomer).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}
