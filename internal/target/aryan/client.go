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
	"time"

	"erp-job/internal/circuitbreaker"
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
	saleCenterSelectPath        = "SaleCenterSelect"
	deliverCenterSaleSelectPath = "DeliverCenter_SaleSelect"
	salerSelectPath             = "SalerSelect"
	saleSellerVisitorPath       = "SaleSellerVisitor"
	goodsPath                   = "Goods"
	addSubElementSelectPath     = "AddSubElementSelect"
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
	breaker     *circuitbreaker.Breaker
}

func NewClient(cfg config.AryanApp, telemetry *observability.Telemetry, log *zap.SugaredLogger) *Client {
	return newClient(cfg, telemetry, log, retry.DefaultPolicy(), circuitbreaker.New(5, 30*time.Second))
}

func newClient(cfg config.AryanApp, telemetry *observability.Telemetry, log *zap.SugaredLogger, policy retry.Policy, breaker *circuitbreaker.Breaker) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/") + "/",
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		retryPolicy: policy,
		telemetry:   telemetry,
		log:         log,
		breaker:     breaker,
	}
}

func (c *Client) PostInvoiceToSaleFactor(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
	for i, item := range invoices {
		payloads = append(payloads, ParamsPayload{
			ID: saleFactorPath,
			Params: []ParamEntry{
				{Name: "isrow", Value: "1"},
				{Name: "CustomerId", Value: item.CustomerID},
				{Name: "VoucherDate", Value: item.InvoiceDate},
				{Name: "SecondNumber", Value: item.InvoiceId},
				{Name: "VoucherDesc", Value: "ETL-Form Fararavand"},
				{Name: "StockId", Value: item.WareHouseID},
				{Name: "SaleTypeId", Value: item.SNoePardakht},
				{Name: "DeliveryCenterID", Value: item.SNoePardakht},
				{Name: "SaleCenterID", Value: item.WareHouseID},
				{Name: "PaymentWayID", Value: item.SNoePardakht},
				{Name: "SellerVisitorID", Value: item.CCForoshandeh},
				{Name: "[Inserted]", ArrayValue: []interface{}{i + 1, item.ProductID, item.ProductCount, item.ProductFee, item.TozihatFaktor}},
				{Name: "[el_Inserted]", ArrayValue: []interface{}{item.CCForoshandeh, 1000, 0}},
			},
		})
	}

	return c.postJSON(ctx, saleFactorPath, payloads)
}

func (c *Client) PostProductsToGoods(ctx context.Context, products []domain.Products) error {
	payloads := make([]ParamsPayload, 0, len(products))
	for _, item := range products {
		payloads = append(payloads, ParamsPayload{
			ID: goodsPath,
			Params: []ParamEntry{
				{Name: "ServiceGoodsID", Value: item.ID},
				{Name: "ServiceGoodsCode", Value: item.Code},
				{Name: "ServiceGoodsDesc", Value: item.Name},
				{Name: "GroupId", Value: item.FirstProductGroupID},
				{Name: "TypeID", Value: 0},
				{Name: "SecUnitType", Value: 0},
				{Name: "Level1", Value: 0},
				{Name: "Level2", Value: 0},
				{Name: "Level3", Value: 0},
			},
		})
	}

	return c.postJSON(ctx, goodsPath, payloads)
}

func (c *Client) PostCustomerToSaleCustomer(ctx context.Context, customers []domain.Customers) error {
	payloads := make([]ParamsPayload, 0, len(customers))
	for _, item := range customers {
		payloads = append(payloads, ParamsPayload{
			ID: saleCustomerPath,
			Params: []ParamEntry{
				{Name: "CustomerID", Value: item.ID},
				{Name: "CustomerCode", Value: strconv.Itoa(item.Code)},
			},
		})
	}

	return c.postJSON(ctx, saleCustomerPath, payloads)
}

func (c *Client) PostInvoiceToSaleOrder(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
	for i, item := range invoices {
		visitorID, err := domain.ParseVisitorCode(item.VisitorCode)
		if err != nil {
			return err
		}
		payloads = append(payloads, ParamsPayload{
			ID: saleOrderPath,
			Params: []ParamEntry{
				{Name: "isrow", Value: "1"},
				{Name: "CustomerId", Value: item.CustomerID},
				{Name: "VoucherDate", Value: item.InvoiceDate},
				{Name: "SecondNumber", Value: item.InvoiceId},
				{Name: "VoucherDesc", Value: "ETL-Form Fararavand"},
				{Name: "StockId", Value: item.WareHouseID},
				{Name: "SaleTypeId", Value: item.SNoePardakht},
				{Name: "DeliveryCenterID", Value: item.SNoePardakht},
				{Name: "SaleCenterID", Value: item.WareHouseID},
				{Name: "PaymentWayID", Value: item.SNoePardakht},
				{Name: "SellerVisitorID", Value: visitorID},
				{Name: "[Inserted]", ArrayValue: []interface{}{i + 1, item.ProductID, item.ProductCount, item.ProductFee, item.TozihatFaktor}},
				{Name: "[el_Inserted]", ArrayValue: []interface{}{visitorID, 1000, 0}},
			},
		})
	}

	return c.postJSON(ctx, saleOrderPath, payloads)
}

func (c *Client) PostInvoiceToSalePayment(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
	for _, item := range invoices {
		payloads = append(payloads, ParamsPayload{
			ID: salePaymentSelectPath,
			Params: []ParamEntry{
				{Name: "PaymentWayID", Value: item.SNoePardakht},
				{Name: "PaymentwayDesc", Value: item.TxtNoePardakht},
			},
		})
	}

	return c.postJSON(ctx, salePaymentSelectPath, payloads)
}

func (c *Client) PostInvoiceToSaleCenter(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
	for _, item := range invoices {
		payloads = append(payloads, ParamsPayload{
			ID: saleCenter4SaleSelectPath,
			Params: []ParamEntry{
				{Name: "StockId", Value: item.WareHouseID},
				{Name: "StockCode", Value: strconv.Itoa(item.WareHouseID)},
				{Name: "StockDesc", Value: item.NameAnbar},
			},
		})
	}

	return c.postJSON(ctx, saleCenter4SaleSelectPath, payloads)
}

func (c *Client) PostInvoiceToSalerSelect(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
	for _, item := range invoices {
		visitorID, err := domain.ParseVisitorCode(item.VisitorCode)
		if err != nil {
			return err
		}
		payloads = append(payloads, ParamsPayload{
			ID: salerSelectPath,
			Params: []ParamEntry{
				{Name: "SaleVisitorID", Value: visitorID},
				{Name: "SaleVisitorDesc", Value: item.VisitorName},
			},
		})
	}

	return c.postJSON(ctx, salerSelectPath, payloads)
}

func (c *Client) PostInvoiceToSaleProforma(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
	for i, item := range invoices {
		visitorCode, err := domain.ParseVisitorCode(item.VisitorCode)
		if err != nil {
			return err
		}
		payloads = append(payloads, ParamsPayload{
			ID: saleProformaPath,
			Params: []ParamEntry{
				{Name: "isrow", Value: "1"},
				{Name: "CustomerId", Value: item.CustomerID},
				{Name: "VoucherDate", Value: item.InvoiceDate},
				{Name: "SecondNumber", Value: item.InvoiceId},
				{Name: "VoucherDesc", Value: "ETL-Form Fararavand"},
				{Name: "StockId", Value: item.WareHouseID},
				{Name: "SaleTypeId", Value: item.SNoePardakht},
				{Name: "DeliveryCenterID", Value: item.SNoePardakht},
				{Name: "SaleCenterID", Value: item.WareHouseID},
				{Name: "PaymentWayID", Value: item.SNoePardakht},
				{Name: "SellerVisitorID", Value: visitorCode},
				{Name: "[Inserted]", ArrayValue: []interface{}{i + 1, item.ProductID, item.ProductCount, item.ProductFee, item.TozihatFaktor}},
				{Name: "[el_Inserted]", ArrayValue: []interface{}{visitorCode, 1000, 0}},
			},
		})
	}

	return c.postJSON(ctx, saleProformaPath, payloads)
}

func (c *Client) PostInvoiceToSaleTypeSelect(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
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

		payloads = append(payloads, ParamsPayload{
			ID: saleTypeSelectPath,
			Params: []ParamEntry{
				{Name: "BuySaleTypeID", Value: item.SNoePardakht},
				{Name: "BuySaleTypeCode", Value: strconv.Itoa(item.SNoePardakht)},
				{Name: "BuySaleTypeDesc", Value: item.TxtNoePardakht},
			},
		})
	}

	return c.postJSON(ctx, saleTypeSelectPath, payloads)
}

func (c *Client) PostBaseDataToDeliverCenterSaleSelect(ctx context.Context, baseData domain.BaseData) error {
	payloads := make([]ParamsPayload, 0, len(baseData.PaymentTypes))
	for _, item := range baseData.PaymentTypes {
		payloads = append(payloads, ParamsPayload{
			ID: deliverCenterSaleSelectPath,
			Params: []ParamEntry{
				{Name: "CentersID", Value: item.ID},
				{Name: "CentersCode", Value: strconv.Itoa(item.ID)},
			},
		})
	}

	return c.postJSON(ctx, deliverCenterSaleSelectPath, payloads)
}

func (c *Client) PostBaseDataToSaleCenterSelect(ctx context.Context, baseData domain.BaseData) error {
	payloads := make([]ParamsPayload, 0, len(baseData.Branches))
	for _, item := range baseData.Branches {
		payloads = append(payloads, ParamsPayload{
			ID: saleCenterSelectPath,
			Params: []ParamEntry{
				{Name: "CentersID", Value: item.ID},
				{Name: "CentersCode", Value: strconv.Itoa(item.ID)},
				{Name: "CenterDesc", Value: item.Name},
			},
		})
	}

	return c.postJSON(ctx, saleCenterSelectPath, payloads)
}

func (c *Client) PostInvoiceToSaleSellerVisitor(ctx context.Context, invoices []domain.Invoices) error {
	payloads := make([]ParamsPayload, 0, len(invoices))
	for _, item := range invoices {
		visitorID, err := domain.ParseVisitorCode(item.VisitorCode)
		if err != nil {
			return err
		}
		payloads = append(payloads, ParamsPayload{
			ID: saleSellerVisitorPath,
			Params: []ParamEntry{
				{Name: "CentersID", Value: visitorID},
				{Name: "CentersCode", Value: strconv.Itoa(visitorID)},
			},
		})
	}

	return c.postJSON(ctx, saleSellerVisitorPath, payloads)
}

func (c *Client) PostAddSubElementSelect(ctx context.Context) error {
	payload := ParamsPayload{
		ID: addSubElementSelectPath,
		Params: []ParamEntry{
			{Name: "IsStruct", Value: 1},
		},
	}

	return c.postJSON(ctx, addSubElementSelectPath, payload)
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

	if !c.breaker.Allow() {
		c.telemetry.RecordFailure(ctx, "aryan", endpoint, 0, "circuit_breaker_open")
		return fmt.Errorf("aryan %s: %w", endpoint, circuitbreaker.ErrOpen)
	}

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
		c.breaker.RecordFailure()
		span.SetAttributes(attribute.Int("http.status_code", result.StatusCode))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	c.breaker.RecordSuccess()
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

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		resBody, _ := io.ReadAll(res.Body)
		return res.StatusCode, fmt.Errorf("aryan request failed. status: %d, response: %s", res.StatusCode, strings.TrimSpace(string(resBody)))
	}

	return res.StatusCode, nil
}
