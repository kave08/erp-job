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
	"erp-job/internal/observability"
	"erp-job/internal/retry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	saleOrderPath               = "SaleOrder"
	saleCustomerPath            = "SaleCustomer"
	saleTypeSelectPath          = "SaleTypeSelect"
	saleCenter4SaleSelectPath   = "SaleCenter4SaleSelect"
	salePaymentSelectPath       = "SalePaymentSelect"
	deliverCenterSaleSelectPath = "DeliverCenter_SaleSelect"
	salerSelectPath             = "SalerSelect"
	goodsPath                   = "Goods"
	saleProformaPath            = "SaleProforma"
	saleFactorPath              = "SaleFactor"
)

type Client struct {
	baseURL     string
	apiKey      string
	httpClient  *http.Client
	retryPolicy retry.Policy
	telemetry   *observability.Telemetry
	log         *zap.SugaredLogger
}

func NewClient(cfg config.AryanApp, telemetry *observability.Telemetry, log *zap.SugaredLogger) *Client {
	return newClient(cfg, telemetry, log, retry.DefaultPolicy())
}

func newClient(cfg config.AryanApp, telemetry *observability.Telemetry, log *zap.SugaredLogger, policy retry.Policy) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/") + "/",
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		retryPolicy: policy,
		telemetry:   telemetry,
		log:         log,
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
			SaleTypeID:       item.SNoePardakht,
			DeliveryCenterID: item.SNoePardakht,
			SaleCenterID:     item.WareHouseID,
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
		visitorID, err := domain.ParseVisitorCode(item.VisitorCode)
		if err != nil {
			return err
		}
		payload = append(payload, domain.SaleOrder{
			CustomerId:       item.CustomerID,
			VoucherDate:      item.InvoiceDate,
			SecondNumber:     item.InvoiceId,
			VoucherDesc:      "",
			StockID:          item.WareHouseID,
			SaleTypeId:       item.SNoePardakht,
			DeliveryCenterID: item.SNoePardakht,
			SaleCenterID:     item.WareHouseID,
			PaymentWayID:     item.SNoePardakht,
			SellerVisitorID:  visitorID,
			ServiceGoodsID:   item.ProductID,
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
			PaymentWayID:   item.SNoePardakht,
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
		visitorID, err := domain.ParseVisitorCode(item.VisitorCode)
		if err != nil {
			return err
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
		visitorCode, err := domain.ParseVisitorCode(item.VisitorCode)
		if err != nil {
			return err
		}
		payload = append(payload, domain.SaleProforma{
			CustomerId:       item.CustomerID,
			VoucherDate:      item.Date,
			SecondNumber:     strconv.Itoa(item.InvoiceId),
			VoucherDesc:      "",
			StockID:          item.WareHouseID,
			SaleTypeId:       item.SNoePardakht,
			DeliveryCenterID: item.SNoePardakht,
			SaleCenterID:     item.WareHouseID,
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
	descriptions := make(map[int]string, len(invoices))
	for _, item := range invoices {
		if item.TxtNoePardakht != "" {
			if previous, exists := descriptions[item.SNoePardakht]; exists && previous != "" && previous != item.TxtNoePardakht && c.log != nil {
				c.log.Warnw("conflicting sale type descriptions in invoice batch",
					"sale_type_id", item.SNoePardakht,
					"previous_desc", previous,
					"current_desc", item.TxtNoePardakht,
				)
			}
			descriptions[item.SNoePardakht] = item.TxtNoePardakht
		}

		payload = append(payload, domain.SaleTypeSelect{
			BuySaleTypeID:   item.SNoePardakht,
			BuySaleTypeCode: strconv.Itoa(item.SNoePardakht),
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
	ctx, span := c.telemetry.Tracer("erp-job/target/aryan").Start(ctx, "aryan.post",
		trace.WithAttributes(
			attribute.String("component", "aryan"),
			attribute.String("endpoint_group", endpoint),
			attribute.String("http.method", http.MethodPost),
		),
	)
	defer span.End()

	attemptObserver := observability.AttemptObserverFromContext(ctx)
	observer := func(attempt retry.Attempt) {
		observability.LogHTTPAttempt(c.log, ctx, "aryan", endpoint, attempt)
		if attempt.WillRetry {
			c.telemetry.RecordRetry(ctx, "aryan", endpoint)
		}
		if attempt.Error != nil && !attempt.WillRetry {
			c.telemetry.RecordFailure(ctx, "aryan", endpoint, attempt.StatusCode, observability.ClassifyHTTPError(attempt.StatusCode, attempt.Error))
		}
		if attemptObserver != nil {
			attemptObserver(observability.HTTPAttempt{
				Endpoint:   endpoint,
				Attempt:    attempt.Attempt,
				StatusCode: attempt.StatusCode,
				Error:      attempt.Error,
				WillRetry:  attempt.WillRetry,
				Duration:   attempt.Duration,
			})
		}
	}

	result, err := retry.Do(ctx, c.retryPolicy, observer, func(ctx context.Context) (int, error) {
		return c.postJSONOnce(ctx, endpoint, payload)
	})
	if err != nil {
		span.SetAttributes(attribute.Int("http.status_code", result.StatusCode))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetAttributes(attribute.Int("http.status_code", result.StatusCode))
	return nil
}

func (c *Client) postJSONOnce(ctx context.Context, endpoint string, payload interface{}) (int, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}

	req.Header.Set("ApiKey", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		return res.StatusCode, fmt.Errorf("aryan request failed. status: %d, response: %s", res.StatusCode, strings.TrimSpace(string(resBody)))
	}

	return res.StatusCode, nil
}
