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
)

const (
	getProductsPath  = "/GetProducts?PageNumeber=%d&PageSize=%d&LastId=%d/"
	getInvoicesPath  = "/GetInvoices?PageNumeber=%d&PageSize=%d&LastId=%d/"
	getCustomersPath = "/GetCustomers?PageNumeber=%d&PageSize=%d&LastId=%d/"
	getBaseDataPath  = "/GetBaseData?PageNumeber=%d&PageSize=%d&LastId=%d/"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
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

func NewClient(cfg config.FararavandApp) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *Client) FetchInvoices(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Invoices, error) {
	var response invoicesResponse

	if err := c.get(ctx, fmt.Sprintf(getInvoicesPath, pageNumber, pageSize, lastID), &response); err != nil {
		return nil, err
	}

	return response.NewInvoices, nil
}

func (c *Client) FetchCustomers(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Customers, error) {
	var response customersResponse

	if err := c.get(ctx, fmt.Sprintf(getCustomersPath, pageNumber, pageSize, lastID), &response); err != nil {
		return nil, err
	}

	return response.NewCustomers, nil
}

func (c *Client) FetchProducts(ctx context.Context, pageNumber, pageSize, lastID int) ([]domain.Products, error) {
	var response productsResponse

	if err := c.get(ctx, fmt.Sprintf(getProductsPath, pageNumber, pageSize, lastID), &response); err != nil {
		return nil, err
	}

	return response.NewProducts, nil
}

func (c *Client) FetchBaseData(ctx context.Context, pageNumber, pageSize, lastID int) (domain.BaseData, error) {
	var response baseDataResponse

	if err := c.get(ctx, fmt.Sprintf(getBaseDataPath, pageNumber, pageSize, lastID), &response); err != nil {
		return domain.BaseData{}, err
	}

	return response.NewBaseData, nil
}

func (c *Client) get(ctx context.Context, path string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("ApiKey", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("fararavand request failed. status: %d, response: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return err
	}

	return nil
}
