package aryan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"erp-job/internal/config"
	"erp-job/internal/domain"
)

const (
	saleOrderPath               = "SaleOrder"
	saleCustomerPath            = "SaleCustomer"
	saleTypeSelectPath          = "SaleTypeSelect"
	saleCenter4SaleSelectPath   = "SaleCenter4SaleSelect"
	salePaymentSelectPath       = "SalePaymentSelect"
	saleCenterSelectPath        = "SaleCenterSelect"
	deliverCenterSaleSelectPath = "DeliverCenter_SaleSelect"
	salerSelectPath             = "SalerSelect"
	saleSellerVisitorPath       = "SaleSellerVisitor"
	goodsPath                   = "Goods"
	saleProformaPath            = "SaleProforma"
	saleFactorPath              = "SaleFactor"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(cfg config.AryanApp) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/") + "/",
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *Client) PostInvoiceToSaleFactor(ctx context.Context, invoices []domain.Invoices) error {
	payload := make([]domain.SaleFactor, 0, len(invoices))
	for _, item := range invoices {
		payload = append(payload, domain.SaleFactor{
			CustomerID:       item.CustomerID,
			VoucherDate:      item.InvoiceDate,
			StockID:          item.WareHouseID,
			VoucherDesc:      "ETL-Form Fararavand",
			SaleTypeID:       10000001,
			DeliveryCenterID: 10000002,
			SaleCenterID:     item.CodeMahal,
			PaymentWayID:     item.SNoePardakht,
			SellerID:         item.CCForoshandeh,
			SaleManID:        item.CodeForoshandeh,
			DistributerID:    0,
			SecondNumber:     strconv.Itoa(item.InvoiceId),
			ServiceGoodsID:   item.ProductID,
			Quantity:         float64(item.ProductCount),
			Fee:              float64(item.ProductFee),
			DetailDesc:       item.TozihatFaktor,
		})
	}

	return c.postJSON(ctx, saleFactorPath, payload)
}

func (c *Client) PostProductsToGoods(ctx context.Context, products []domain.Products) error {
	payload := make([]domain.Goods, 0, len(products))
	for _, item := range products {
		payload = append(payload, domain.Goods{
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

	return c.postJSON(ctx, goodsPath, payload)
}

func (c *Client) PostCustomerToSaleCustomer(ctx context.Context, customers []domain.Customers) error {
	payload := make([]domain.SaleCustomer, 0, len(customers))
	for _, item := range customers {
		payload = append(payload, domain.SaleCustomer{
			CustomerID:   item.ID,
			CustomerCode: strconv.Itoa(item.Code),
		})
	}

	return c.postJSON(ctx, saleCustomerPath, payload)
}

func (c *Client) PostInvoiceToSaleOrder(ctx context.Context, invoices []domain.Invoices) error {
	payload := make([]domain.SaleOrder, 0, len(invoices))
	for _, item := range invoices {
		payload = append(payload, domain.SaleOrder{
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

	return c.postJSON(ctx, saleOrderPath, payload)
}

func (c *Client) PostInvoiceToSalePayment(ctx context.Context, invoices []domain.Invoices) error {
	payload := make([]domain.SalePaymentSelect, 0, len(invoices))
	for _, item := range invoices {
		payload = append(payload, domain.SalePaymentSelect{
			PaymentWayID:   item.PaymentTypeID,
			PaymentwayDesc: item.TxtNoePardakht,
		})
	}

	return c.postJSON(ctx, salePaymentSelectPath, payload)
}

func (c *Client) PostInvoiceToSaleCenter(ctx context.Context, invoices []domain.Invoices) error {
	payload := make([]domain.SaleCenter4SaleSelect, 0, len(invoices))
	for _, item := range invoices {
		payload = append(payload, domain.SaleCenter4SaleSelect{
			StockID:   item.WareHouseID,
			StockCode: strconv.Itoa(item.WareHouseID),
			StockDesc: item.NameAnbar,
		})
	}

	return c.postJSON(ctx, saleCenter4SaleSelectPath, payload)
}

func (c *Client) PostInvoiceToSalerSelect(ctx context.Context, invoices []domain.Invoices) error {
	payload := make([]domain.SalerSelect, 0, len(invoices))
	for _, item := range invoices {
		visitorID, err := strconv.Atoi(item.VisitorCode)
		if err != nil {
			return fmt.Errorf("parse visitor code %q: %w", item.VisitorCode, err)
		}
		payload = append(payload, domain.SalerSelect{
			SaleVisitorID:   visitorID,
			SaleVisitorDesc: item.VisitorName,
		})
	}

	return c.postJSON(ctx, salerSelectPath, payload)
}

func (c *Client) PostInvoiceToSaleProforma(ctx context.Context, invoices []domain.Invoices) error {
	payload := make([]domain.SaleProforma, 0, len(invoices))
	for _, item := range invoices {
		visitorCode, err := strconv.Atoi(item.VisitorCode)
		if err != nil {
			return fmt.Errorf("parse visitor code %q: %w", item.VisitorCode, err)
		}
		payload = append(payload, domain.SaleProforma{
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

	return c.postJSON(ctx, saleProformaPath, payload)
}

func (c *Client) PostInvoiceToSaleTypeSelect(ctx context.Context, invoices []domain.Invoices) error {
	payload := make([]domain.SaleTypeSelect, 0, len(invoices))
	for _, item := range invoices {
		payload = append(payload, domain.SaleTypeSelect{
			BuySaleTypeID:   item.SNoePardakht,
			BuySaleTypeCode: item.Codekala,
			BuySaleTypeDesc: item.TxtNoePardakht,
		})
	}

	return c.postJSON(ctx, saleTypeSelectPath, payload)
}

func (c *Client) PostBaseDataToDeliverCenterSaleSelect(ctx context.Context, baseData domain.BaseData) error {
	payload := make([]domain.DeliverCenterSaleSelect, 0, len(baseData.PaymentTypes))
	for _, item := range baseData.PaymentTypes {
		payload = append(payload, domain.DeliverCenterSaleSelect{
			CentersID:   item.ID,
			CentersCode: strconv.Itoa(item.ID),
		})
	}

	return c.postJSON(ctx, deliverCenterSaleSelectPath, payload)
}

func (c *Client) postJSON(ctx context.Context, endpoint string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("ApiKey", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("aryan request failed. status: %d, response: %s", res.StatusCode, strings.TrimSpace(string(resBody)))
	}

	return nil
}
