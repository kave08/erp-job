package aryan

import (
	"bytes"
	"encoding/json"
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Aryan represents the Aryan ERP system's data synchronization service.
type Aryan struct {
	httpClient *http.Client
	baseURL    string
	repos      *repository.Repository
}

// NewAryan initializes and returns a new Aryan service instance.
func NewAryan(repos *repository.Repository) AryanInterface {
	return &Aryan{
		baseURL: config.Cfg.AryanApp.BaseURL,
		repos:   repos,
		httpClient: &http.Client{
			Timeout: config.Cfg.AryanApp.Timeout,
		},
	}
}

// Req defines the structure for a request to the Aryan ERP system.
type Req struct {
	ID     string  `json:"id"`
	Params []Param `json:"Params"`
}

// Param represents a parameter with a name and an optional value.
type Param struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value,omitempty"`
}

// PostInvoiceToSaleFactor synchronizes invoices by converting them into SaleFactors for the Aryan ERP system.
//
// It maps invoice fields to SaleFactor fields and sends a POST request to the sale factor service.
func (a *Aryan) PostInvoiceToSaleFactor(fp []models.Invoices) error {
	var newSaleFactor []models.SaleFactor
	var req = []Req{}

	for _, item := range fp {
		req = append(req, Req{
			ID: "SaleFactor",
			Params: []Param{
				{
					Name:  "CustomerId",
					Value: item.CustomerID,
				},
				{
					Name:  "VoucherDate",
					Value: item.InvoiceDate,
				},
				{
					Name:  "StockID",
					Value: item.WareHouseID,
				},
				{
					Name:  "VoucherDesc",
					Value: "ETL-Form Fararavand",
				},
				{
					Name:  "SaleTypeId",
					Value: 10000001,
				},
				{
					Name:  "DeliveryCenterID",
					Value: 10000002,
				},
				{
					Name:  "SaleCenterID",
					Value: item.CodeMahal,
				},
				{
					Name:  "PaymentWayID",
					Value: item.SNoePardakht,
				},
				{
					Name:  "SellerID",
					Value: item.CCForoshandeh,
				},
				{
					Name:  "SaleManID",
					Value: item.CodeForoshandeh,
				},
				{
					Name:  "DistributerId",
					Value: 0,
				},
				{
					Name:  "CustomerId",
					Value: item.CustomerID,
				},
			},
		})

		//TODO:
		// ServiceGoodsID:   item.ProductID, // ok
		// 	Quantity:         float64(item.ProductCount),
		// 	Fee:              float64(item.ProductFee),
		// 	VoucherDesc:      "ETL-Form Fararavand",
		// 	SecondNumber:     strconv.Itoa(item.InvoiceId),

		// newSaleFactor = append(newSaleFactor, models.SaleFactor{
		// CustomerId:       item.CustomerID,

		// VoucherDate:      item.InvoiceDate,
		// StockID:          item.WareHouseID,
		// SaleTypeId:       10000001,
		// DeliveryCenterID: 10000002,
		// SaleCenterID:     item.CodeMahal,
		// PaymentWayID:     item.SNoePardakht,
		// SellerID:         item.CCForoshandeh,
		// SaleManID:        item.CodeForoshandeh,
		// })s
	}

	body, err := json.Marshal(newSaleFactor)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleFactor, bytes.NewReader(body))
	if err != nil {
		return err
	}

	request.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostProductsToGoods takes a slice of Products and posts them to the goods service.
// It converts each Product into a Goods by copying the ID, Code, and Name fields,
// and setting the GroupId to the FirstProductGroupID. All other fields are set to  0.
// The function then creates a new HTTP POST request with the JSON-encoded slice of Goods
// as the request body and sends it to the goods service endpoint.
// If the request is successful and the response status code is OK, the function returns nil.
func (a *Aryan) PostProductsToGoods(fp []models.Products) error {
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

	body, err := json.Marshal(newGoods)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		Goods, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostCustomerToSaleCustomer takes a slice of Customers and posts them to the sale customer service.
// It converts each Customer into a SaleCustomer by copying the ID and converting the Code to a string.
// The function then makes a POST request to the sale customer service endpoint with the slice of SaleCustomers as the request body.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostCustomerToSaleCustomer(fc []models.Customers) error {
	var newSaleCustomer []models.SaleCustomer

	for _, item := range fc {
		newSaleCustomer = append(newSaleCustomer, models.SaleCustomer{
			CustomerID:   item.ID,
			CustomerCode: strconv.Itoa(item.Code),
		})
	}

	body, err := json.Marshal(newSaleCustomer)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleCustomer, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostInvoiceToSaleOrder takes a slice of Invoices and posts them to the sale order service.
// It converts each Invoice into a SaleOrder by mapping its fields to the corresponding SaleOrder fields.
// The function then sends a POST request with the slice of SaleOrders as the request body to the sale order service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSaleOrder(fp []models.Invoices) error {
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

	body, err := json.Marshal(newSaleOrder)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleOrder, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostInvoiceToSalePayment takes a slice of Invoices and posts them to the sale payment select service.
// It converts each Invoice into a SalePaymentSelect by mapping the PaymentTypeID and TxtNoePardakht fields.
// The function then sends a POST request with the slice of SalePaymentSelect as the request body to the sale payment select service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSalePayment(fp []models.Invoices) error {
	var newSalePaymentSelect []models.SalePaymentSelect

	for _, item := range fp {
		newSalePaymentSelect = append(newSalePaymentSelect, models.SalePaymentSelect{
			PaymentWayID:   item.PaymentTypeID,
			PaymentwayDesc: item.TxtNoePardakht,
		})
	}

	body, err := json.Marshal(newSalePaymentSelect)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SalePaymentSelect, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostInvoiceToSaleCenter takes a slice of Invoices and posts them to the sale center service.
// It converts each Invoice into a SaleCenter4SaleSelect by mapping relevant fields such as StockID and StockDesc.
// The function then makes a POST request to the sale center service endpoint with the slice of SaleCenter4SaleSelect as the request body.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSaleCenter(fp []models.Invoices) error {
	var newSaleCenter []models.SaleCenter4SaleSelect

	for _, item := range fp {
		newSaleCenter = append(newSaleCenter, models.SaleCenter4SaleSelect{
			StockID:   item.WareHouseID,
			StockCode: strconv.Itoa(item.WareHouseID),
			StockDesc: item.NameAnbar,
		})
	}

	body, err := json.Marshal(newSaleCenter)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleCenter4SaleSelect, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostInvoiceToSalerSelect takes a slice of Invoices and posts them to the saler select service.
// It converts each Invoice into a SalerSelect by mapping the VisitorCode and VisitorName fields.
// The function then sends a POST request with the slice of SalerSelect as the request body to the saler select service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSalerSelect(fp []models.Invoices) error {
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

	body, err := json.Marshal(newSalerSelect)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SalerSelect, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostInvoiceToSaleProforma takes a slice of Invoices and posts them to the sale proforma service.
// It converts each Invoice into a SaleProforma by mapping its fields to the corresponding SaleProforma fields.
// The function then sends a POST request with the slice of SaleProforma as the request body to the sale proforma service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSaleProforma(fp []models.Invoices) error {
	var newSaleProforma []models.SaleProforma

	for _, item := range fp {
		visitorCode, err := strconv.Atoi(item.VisitorCode)
		if err != nil {

			return err
		}
		newSaleProforma = append(newSaleProforma, models.SaleProforma{
			CustomerId:       item.CustomerID,
			VoucherDate:      item.Date,
			SecondNumber:     "",
			VoucherDesc:      "",
			StockID:          item.WareHouseID,
			SaleTypeId:       item.SNoePardakht,
			DeliveryCenterID: 0,
			SaleCenterID:     0,
			PaymentWayID:     item.SNoePardakht,
			SellerVisitorID:  visitorCode,
			ServiceGoodsID:   item.ProductID,
			Quantity:         float64(item.ProductCount),
			Fee:              float64(item.ProductFee),
			DetailDesc:       item.TozihatFaktor,
		})
	}

	body, err := json.Marshal(newSaleProforma)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleProforma, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostInvoiceToSaleTypeSelect takes a slice of Invoices and posts them to the sale type select service.
// It converts each Invoice into a SaleTypeSelect by mapping its fields to the corresponding SaleTypeSelect fields.
// The function then sends a POST request with the slice of SaleTypeSelect as the request body to the sale proforma service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostInvoiceToSaleTypeSelect(fp []models.Invoices) error {
	var newSaleTypeSelect []models.SaleTypeSelect

	for _, item := range fp {
		newSaleTypeSelect = append(newSaleTypeSelect, models.SaleTypeSelect{
			BuySaleTypeID:   item.SNoePardakht,
			BuySaleTypeCode: item.Codekala,
			BuySaleTypeDesc: item.TxtNoePardakht,
		})
	}

	body, err := json.Marshal(newSaleTypeSelect)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleTypeSelect, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostBaseDataToSaleCenterSelect takes a slice of payment types of BaseData and posts them to the sale center select service.
// It converts each base data into a SaleCenterSelect by mapping its fields to the corresponding SaleCenterSelect fields.
// The function then sends a POST request with the slice of SaleCenterSelect as the request body to the sale proforma service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostBaseDataToSaleCenterSelect(baseData models.BaseData) error {
	var newSaleCenterSelect []models.SaleCenterSelect

	for _, item := range baseData.PaymentTypes {
		newSaleCenterSelect = append(newSaleCenterSelect, models.SaleCenterSelect{
			CentersID:   item.ID,
			CentersCode: strconv.Itoa(item.ID),
			CenterDesc:  item.Name,
		})
	}

	body, err := json.Marshal(newSaleCenterSelect)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleCenterSelect, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostBaseDataToDeliverCenterSaleSelect takes a slice of payment types of BaseData and posts them to the deliver center sale select service.
// It converts each PaymentTypes into a DeliverCenterSaleSelect by mapping its fields to the corresponding DeliverCenterSaleSelect fields.
// The function then sends a POST request with the slice of DeliverCenterSaleSelect as the request body to the deliver center sale select service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostBaseDataToDeliverCenterSaleSelect(baseData models.BaseData) error {
	var newADeliverCenterSaleSelect []models.DeliverCenter_SaleSelect

	for _, item := range baseData.PaymentTypes {
		newADeliverCenterSaleSelect = append(newADeliverCenterSaleSelect, models.DeliverCenter_SaleSelect{
			CentersID:   item.ID,
			CentersCode: strconv.Itoa(item.ID),
		})
	}

	body, err := json.Marshal(newADeliverCenterSaleSelect)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		DeliverCenterSaleSelect, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}

// PostBaseDataToSaleSellerVisitor take a struct of payment types of BaseData and posts them to the sale seller visitor service.
// It converts each PaymentTypes into a SaleSellerVisitor by mapping its fields to the corresponding SaleSellerVisitor fields.
// The function then sends a POST request with the slice of SaleSellerVisitor as the request body to the sale seller visitor service endpoint.
// The function returns the server response and an error if the request fails.
func (a *Aryan) PostBaseDataToSaleSellerVisitor(baseData models.BaseData) error {
	var newSaleSellerVisitor []models.SaleSellerVisitor

	for _, item := range baseData.PaymentTypes {
		newSaleSellerVisitor = append(newSaleSellerVisitor, models.SaleSellerVisitor{
			CentersID:   item.ID,
			CentersCode: strconv.Itoa(item.ID),
		})
	}

	body, err := json.Marshal(newSaleSellerVisitor)
	if err != nil {

		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+
		SaleSellerVisitor, bytes.NewReader(body))
	if err != nil {

		return err
	}

	req.Header.Set("ApiKey", config.Cfg.AryanApp.APIKey)

	res, err := a.httpClient.Do(req)
	if err != nil {

		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)

		return fmt.Errorf("http request failed. status: %d, response: %s", res.StatusCode, resBody)
	}

	return nil
}
