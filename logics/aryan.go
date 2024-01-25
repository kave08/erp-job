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
	PostInvoiceToSaleTypeSelect(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToSaleCenterSelect(fp []models.Invoices) (*resty.Response, error)
	PostInvoiceToDeliverCenterSaleSelect(fp []models.Invoices) (*resty.Response, error)
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

// PostInoviceToSaleFactor takes a slice of Invoices and converts them into SaleFactors.
// Each Invoice is transformed into a SaleFactor by mapping its fields to the corresponding SaleFactor fields.
// The function then sends a POST request with the slice of SaleFactors as the request body to the sale factor service.
// The function returns the server response and an error if the request fails.
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

	res, err := a.restyClient.R().SetBody(newSaleFactor).Post(ASaleFactor)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostProductsToGoods takes a slice of Products and posts them to the goods service.
// It converts each Product into a Goods structure by mapping relevant fields.
// The function then makes a POST request to the goods service endpoint with the slice of Goods as the request body.
// The function returns the server response and an error if the request fails.
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

	res, err := a.restyClient.R().SetBody(newGoods).Post(AGoods)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostCustomerToSaleCustomer takes a slice of Customers and posts them to the sale customer service.
// It converts each Customer into a SaleCustomer by copying the ID and converting the Code to a string.
// The function then makes a POST request to the sale customer service endpoint with the slice of SaleCustomers as the request body.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostCustomerToSaleCustomer(fc []models.Customers) (*resty.Response, error) {
	var newSaleCustomer []models.SaleCustomer

	for _, item := range fc {
		newSaleCustomer = append(newSaleCustomer, models.SaleCustomer{
			CustomerID:   item.ID,
			CustomerCode: strconv.Itoa(item.Code),
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleCustomer).Post(ASaleCustomer)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostInvoiceToSaleOrder takes a slice of Invoices and posts them to the sale order service.
// It converts each Invoice into a SaleOrder by mapping its fields to the corresponding SaleOrder fields.
// The function then sends a POST request with the slice of SaleOrders as the request body to the sale order service endpoint.
// The function returns the server response and an error if the request fails.
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

	res, err := a.restyClient.R().SetBody(newSaleOrder).Post(ASaleOrder)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostInvoiceToSaleCenter takes a slice of Invoices and posts them to the sale center service.
// It converts each Invoice into a SaleCenter4SaleSelect by mapping relevant fields such as StockID and StockDesc.
// The function then makes a POST request to the sale center service endpoint with the slice of SaleCenter4SaleSelect as the request body.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSaleCenter(fp []models.Invoices) (*resty.Response, error) {
	var newSaleCenter []models.SaleCenter4SaleSelect

	for _, item := range fp {
		newSaleCenter = append(newSaleCenter, models.SaleCenter4SaleSelect{
			StockID:   item.WareHouseID,
			StockCode: strconv.Itoa(item.WareHouseID),
			StockDesc: item.NameAnbar,
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleCenter).Post(ASaleCenter4SaleSelect)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostInvoiceToSalePaymentSelect takes a slice of Invoices and posts them to the sale payment select service.
// It converts each Invoice into a SalePaymentSelect by mapping the PaymentTypeID and TxtNoePardakht fields.
// The function then sends a POST request with the slice of SalePaymentSelect as the request body to the sale payment select service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSalePaymentSelect(fp []models.Invoices) (*resty.Response, error) {
	var newSalePaymentSelect []models.SalePaymentSelect

	for _, item := range fp {
		newSalePaymentSelect = append(newSalePaymentSelect, models.SalePaymentSelect{
			PaymentWayID:   item.PaymentTypeID,
			PaymentwayDesc: item.TxtNoePardakht,
		})
	}

	res, err := a.restyClient.R().SetBody(newSalePaymentSelect).Post(ASalePaymentSelect)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostInvoiceToSalerSelect takes a slice of Invoices and posts them to the saler select service.
// It converts each Invoice into a SalerSelect by mapping the VisitorCode and VisitorName fields.
// The function then sends a POST request with the slice of SalerSelect as the request body to the saler select service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSalerSelect(fp []models.Invoices) (*resty.Response, error) {
	var newSalerSelect []models.SalerSelect

	for _, item := range fp {
		visitorID, err := strconv.Atoi(item.VisitorCode)
		if err != nil {
			fmt.Println("Error converting VisitorCode to int:", err)
			continue
		}
		newSalerSelect = append(newSalerSelect, models.SalerSelect{
			SaleVisitorID:   visitorID,
			SaleVisitorDesc: item.VisitorName,
		})
	}

	res, err := a.restyClient.R().SetBody(newSalerSelect).Post(ASalerSelect)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// TODO: update feilds
// PostInvoiceToSaleProforma takes a slice of Invoices and posts them to the sale proforma service.
// It converts each Invoice into a SaleProforma by mapping its fields to the corresponding SaleProforma fields.
// The function then sends a POST request with the slice of SaleProforma as the request body to the sale proforma service endpoint.
// The function returns the server response and an error if the request fails.
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

	res, err := a.restyClient.R().SetBody(newSaleProforma).Post(ASaleProforma)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// TODO: update feilds
// PostInvoiceToSaleTypeSelect takes a slice of Invoices and posts them to the sale type select service.
// It converts each Invoice into a SaleTypeSelect by mapping its fields to the corresponding SaleTypeSelect fields.
// The function then sends a POST request with the slice of SaleTypeSelect as the request body to the sale proforma service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSaleTypeSelect(fp []models.Invoices) (*resty.Response, error) {
	var newSaleTypeSelect []models.SaleTypeSelect

	for _, item := range fp {
		newSaleTypeSelect = append(newSaleTypeSelect, models.SaleTypeSelect{
			BuySaleTypeID:   item.BranchID,
			BuySaleTypeCode: "",
			BuySaleTypeDesc: "",
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleTypeSelect).Post(ASaleTypeSelect)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// TODO: update feilds
// PostInvoiceToSaleCenterSelect takes a slice of Invoices and posts them to the sale type select service.
// It converts each Invoice into a SaleCenterSelect by mapping its fields to the corresponding SaleCenterSelect fields.
// The function then sends a POST request with the slice of SaleCenterSelect as the request body to the sale proforma service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSaleCenterSelect(fp []models.Invoices) (*resty.Response, error) {
	var newSaleCenterSelect []models.SaleCenterSelect

	for _, item := range fp {
		newSaleCenterSelect = append(newSaleCenterSelect, models.SaleCenterSelect{
			CentersID:   item.InvoiceId,
			CentersCode: "",
			CenterDesc:  "",
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleCenterSelect).Post(ASaleCenterSelect)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// TODO: update feilds
// PostInvoiceToDeliverCenterSaleSelect takes a slice of Invoices and posts them to the sale type select service.
// It converts each Invoice into a ADeliverCenterSaleSelect by mapping its fields to the corresponding ADeliverCenterSaleSelect fields.
// The function then sends a POST request with the slice of ADeliverCenterSaleSelect as the request body to the sale proforma service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToDeliverCenterSaleSelect(fp []models.Invoices) (*resty.Response, error) {
	var newADeliverCenterSaleSelect []models.DeliverCenter_SaleSelect

	for _, item := range fp {
		newADeliverCenterSaleSelect = append(newADeliverCenterSaleSelect, models.DeliverCenter_SaleSelect{
			CentersID:   item.BranchID,
			CentersCode: "",
		})
	}

	res, err := a.restyClient.R().SetBody(newADeliverCenterSaleSelect).Post(ADeliverCenterSaleSelect)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}
