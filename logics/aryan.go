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
	PostProductsToGoods(fp []models.Products) (*resty.Response, error)
	PostCustomerToSaleCustomer(fc []models.Customers) (*resty.Response, error)
	PostInvoiceToSaleOrder(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleCenter(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSalePaymentSelect(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSalerSelect(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleProforma(fp []models.Invoices) (*resty.Response, error)
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
			StockID:          item.WareHouseID,
			SaleTypeId:       10000001,
			DeliveryCenterID: 10000002,
			SaleCenterID:     item.CodeMahal,
			PaymentWayID:     item.SNoePardakht,
			SellerID:         item.CCForoshandeh,
			SaleManID:        item.CodeForoshandeh,
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
func (a *Aryan) PostProductsToGoods(fp []models.Products) (*resty.Response, error) {
	var newGoods []models.Goods

	for _, item := range fp {
		newGoods = append(newGoods, models.Goods{
			ServiceGoodsID:   item.ID,
			ServiceGoodsCode: item.Code,
			ServiceGoodsDesc: item.Name,
			GroupId:          item.FirstProductGroupID,
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
			CustomerID:   item.ID,
			CustomerCode: strconv.Itoa(item.Code),
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
			VoucherDate:      item.InvoiceDate,
			SecondNumber:     0,
			VoucherDesc:      "",
			StockID:          item.WareHouseID,
			SaleTypeId:       10000001,
			DeliveryCenterID: 10000002,
			SaleCenterID:     item.CodeMahal,
			PaymentWayID:     item.SNoePardakht,
			SellerVisitorID:  item.CCForoshandeh,
			ServiceGoodsID:   0,
			Quantity:         float64(item.ProductCount),
			Fee:              float64(item.ProductFee),
			DetailDesc:       item.TozihatFaktor,
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

func (a *Aryan) PostInvoiceToSaleCenter(fp []models.Invoices) (*resty.Response, error) {
	var newSaleCenter []models.SaleCenter4SaleSelect

	for _, item := range fp {
		newSaleCenter = append(newSaleCenter, models.SaleCenter4SaleSelect{
			StockID:   item.WareHouseID,
			StockCode: strconv.Itoa(item.WareHouseID),
			StockDesc: item.NameAnbar,
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
			PaymentwayDesc: item.TxtNoePardakht,
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
			SaleVisitorID:   strconv.Atoi(item.VisitorCode), //TODO: convert it
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
			VoucherDate:      item.Date,
			SecondNumber:     "",
			VoucherDesc:      "",
			StockID:          item.WareHouseID,
			SaleTypeId:       0,
			DeliveryCenterID: 0,
			SaleCenterID:     0,
			PaymentWayID:     0,
			SellerVisitorID:  0,
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
