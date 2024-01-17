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
	PostInoviceToSaleFactor(fp []models.Invoices) (*resty.Response, error)
	PostCustomerToSaleCustomer(fc []models.Customers) (*resty.Response, error)
	PostInvoiceToGoods(fp []models.Invoices) (*resty.Response, error)
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
func (a *Aryan) PostInoviceToSaleFactor(fp []models.Invoices) (*resty.Response, error) {
	var newSaleFactor []models.SaleFactor

	for _, item := range fp {
		newSaleFactor = append(newSaleFactor, models.SaleFactor{
			CustomerId:       item.CustomerID,
			ServiceGoodsID:   item.ProductID, // ok
			Quantity:         float64(item.ProductCount),
			Fee:              float64(item.ProductFee),
			VoucherDesc:      "ETL-Form Fararavand",
			SecondNumber:     strconv.Itoa(item.InvoiceId),
			VoucherDate:      item.InvoiceDate,
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
func (a *Aryan) PostInvoiceToGoods(fp []models.Invoices) (*resty.Response, error) {
	var newGoods []models.Goods

	for _, item := range fp {
		newGoods = append(newGoods, models.Goods{
			ServiceGoodsID:   item.ProductID,
			ServiceGoodsCode: "",
			ServiceGoodsDesc: item.NameKalaFaktor,
			GroupId:          0,
			TypeID:           0,
			SecUnitType:      0,
			Level1:           0,
			Level2:           0,
			Level3:           0,
		})
	}

	res, err := a.restyClient.R().SetBody(newGoods).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

func (a *Aryan) PostCustomerToSaleCustomer(fc []models.Customers) (*resty.Response, error) {
	var newSaleCustomer []models.SaleCustomer

	for _, item := range fc {
		newSaleCustomer = append(newSaleCustomer, models.SaleCustomer{
			CustomerID:   item.CustomerId,
			CustomerCode: strconv.Itoa(item.CustomerCodePosti),
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

func (a *Aryan) PostInvoiceToSaleOrder(fp []models.Invoices) (*resty.Response, error) {
	var newSaleOrder []models.SaleOrder

	for _, item := range fp {
		newSaleOrder = append(newSaleOrder, models.SaleOrder{
			CustomerId:       item.CustomerID,
			VoucherDate:      "",
			SecondNumber:     0,
			VoucherDesc:      "",
			StockID:          item.WareHouseID,
			SaleTypeId:       0,
			DeliveryCenterID: 0,
			SaleCenterID:     0,
			PaymentWayID:     0,
			SellerVisitorID:  item.VisitorCode,
			ServiceGoodsID:   item.Codekala,
			Quantity:         float64(item.ProductCount),
			Fee:              float64(item.ProductFee),
			DetailDesc:       "",
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleOrder).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

func (a *Aryan) PostProductsToSaleCenter(fp []models.Products) (*resty.Response, error) {
	var newSaleCenter []models.SaleCenter4SaleSelect

	for _, item := range fp {
		newSaleCenter = append(newSaleCenter, models.SaleCenter4SaleSelect{
			StockID:   item.ProductId,
			StockCode: strconv.Itoa(item.Codekala),
			StockDesc: item.Name,
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleCenter).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

func (a *Aryan) PostInvoiceToSalePaymentSelect(fp []models.Invoices) (*resty.Response, error) {
	var newSalePaymentSelect []models.SalePaymentSelect

	for _, item := range fp {
		newSalePaymentSelect = append(newSalePaymentSelect, models.SalePaymentSelect{
			PaymentWayID:   item.PaymentTypeID,
			PaymentwayDesc: "",
		})
	}

	res, err := a.restyClient.R().SetBody(newSalePaymentSelect).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

func (a *Aryan) PostInvoiceToSalerSelect(fp []models.Invoices) (*resty.Response, error) {
	var newSalerSelect []models.SalerSelect

	for _, item := range fp {
		newSalerSelect = append(newSalerSelect, models.SalerSelect{
			SaleVisitorID:   item.VisitorCode,
			SaleVisitorDesc: item.VisitorName,
		})
	}

	res, err := a.restyClient.R().SetBody(newSalerSelect).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

func (a *Aryan) PostInvoiceToSaleProforma(fp []models.Invoices) (*resty.Response, error) {
	var newSaleProforma []models.SaleProforma

	for _, item := range fp {
		newSaleProforma = append(newSaleProforma, models.SaleProforma{
			CustomerId:       item.CustomerID,
			VoucherDate:      0,
			SecondNumber:     "",
			VoucherDesc:      "",
			StockID:          item.WareHouseID,
			SaleTypeId:       0,
			DeliveryCenterID: 0,
			SaleCenterID:     0,
			PaymentWayID:     0,
			SellerVisitorID:  item.VisitorCode,
			ServiceGoodsID:   item.ProductID,
			Quantity:         float64(item.ProductCount),
			Fee:              float64(item.ProductFee),
			DetailDesc:       "",
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleProforma).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}
