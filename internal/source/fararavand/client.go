package fararavand

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	getProductsPath  = "/GetProducts?PageNumber=%d&PageSize=%d&LastId=%d/"
	getInvoicesPath  = "/GetInvoices?PageNumber=%d&PageSize=%d&LastId=%d/"
	getCustomersPath = "/GetCustomers?PageNumber=%d&PageSize=%d&LastId=%d/"
	getBaseDataPath  = "/GetBaseData?PageNumber=%d&PageSize=%d&LastId=%d/"
)

type Client struct {
	baseURL     string
	apiKey      string
	httpClient  *http.Client
	retryPolicy retry.Policy
	telemetry   *observability.Telemetry
	log         *zap.SugaredLogger
}

type invoicesResponse struct {
	Status      int               `json:"status"`
	NewInvoices []domain.Invoices `json:"new_invoice"`
}

type customersResponse struct {
	Status       int                `json:"status"`
	NewCustomers []domain.Customers `json:"new_customer"`
}

type productsResponse struct {
	Status      int               `json:"status"`
	NewProducts []domain.Products `json:"new_products"`
}

type baseDataResponse struct {
	Status      int             `json:"status"`
	NewBaseData domain.BaseData `json:"new_base_data"`
}

func NewClient(cfg config.FararavandApp, telemetry *observability.Telemetry, log *zap.SugaredLogger) *Client {
	return newClient(cfg, telemetry, log, retry.DefaultPolicy())
}

func newClient(cfg config.FararavandApp, telemetry *observability.Telemetry, log *zap.SugaredLogger, policy retry.Policy) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		retryPolicy: policy,
		telemetry:   telemetry,
		log:         log,
	}
}

func (c *Client) FetchInvoices(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Invoices, error) {
	var response invoicesResponse

	if err := c.get(ctx, "GetInvoices", fmt.Sprintf(getInvoicesPath, pageNumber, pageSize, lastID), &response); err != nil {
		return nil, err
	}

	return response.NewInvoices, nil
}

func (c *Client) FetchCustomers(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Customers, error) {
	var response customersResponse

	if err := c.get(ctx, "GetCustomers", fmt.Sprintf(getCustomersPath, pageNumber, pageSize, lastID), &response); err != nil {
		return nil, err
	}

	return response.NewCustomers, nil
}

func (c *Client) FetchProducts(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Products, error) {
	var response productsResponse

	if err := c.get(ctx, "GetProducts", fmt.Sprintf(getProductsPath, pageNumber, pageSize, lastID), &response); err != nil {
		return nil, err
	}

	return response.NewProducts, nil
}

func (c *Client) FetchBaseData(ctx context.Context, pageNumber, pageSize, lastID int) (domain.BaseData, error) {
	var response baseDataResponse

	if err := c.get(ctx, "GetBaseData", fmt.Sprintf(getBaseDataPath, pageNumber, pageSize, lastID), &response); err != nil {
		return domain.BaseData{}, err
	}

	return response.NewBaseData, nil
}

func (c *Client) get(ctx context.Context, endpointGroup, path string, target interface{}) error {
	ctx, span := c.telemetry.Tracer("erp-job/source/fararavand").Start(ctx, "fararavand.get",
		trace.WithAttributes(
			attribute.String("component", "fararavand"),
			attribute.String("endpoint_group", endpointGroup),
			attribute.String("http.method", http.MethodGet),
		),
	)
	defer span.End()

	observer := func(attempt retry.Attempt) {
		observability.LogHTTPAttempt(c.log, ctx, "fararavand", endpointGroup, attempt)
		if attempt.WillRetry {
			c.telemetry.RecordRetry(ctx, "fararavand", endpointGroup)
		}
		if attempt.Error != nil && !attempt.WillRetry {
			c.telemetry.RecordFailure(ctx, "fararavand", endpointGroup, attempt.StatusCode, observability.ClassifyHTTPError(attempt.StatusCode, attempt.Error))
		}
	}

	result, err := retry.Do(ctx, c.retryPolicy, observer, func(ctx context.Context) (int, error) {
		return c.getOnce(ctx, path, target)
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

func (c *Client) getOnce(ctx context.Context, path string, target interface{}) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("ApiKey", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return res.StatusCode, fmt.Errorf("fararavand request failed. status: %d, response: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		body, _ := io.ReadAll(res.Body)
		return res.StatusCode, fmt.Errorf("fararavand returned non-JSON response (content-type: %s): %s", contentType, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return res.StatusCode, err
	}

	return res.StatusCode, nil
}
