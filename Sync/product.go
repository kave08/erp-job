package sync

import (
	"encoding/json"
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"erp-job/utility"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ProductResponse struct {
	Status      int               `json:"status"`
	NewProducts []models.Products `json:"new_products"`
}

type ProductRequest struct {
	LastId      int `json:"LastId"`
	PageSize    int `json:"PageSize"`
	PageNumeber int `json:"PageNumeber"`
}

func NewProductRequest(lastid int, pageSize int, pageNumber int) ProductRequest {
	return ProductRequest{
		LastId:      lastid,
		PageSize:    pageSize,
		PageNumeber: pageNumber,
	}
}

type Product struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewProduct(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface, requestTimeout time.Duration) *Product {

	return &Product{
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (p Product) Products() error {

	request := new(ProductRequest)

	req, err := http.NewRequest(http.MethodGet, p.baseURL+
		fmt.Sprintf("/GetProducts?PageNumeber=%d&PageSize=%d&LastId=%d/", request.PageNumeber, request.PageSize, request.LastId), nil)
	if err != nil {
		return err
	}

	req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

	res, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("get products http request failed. status: %d, response: %v", res.StatusCode, res.Body)
	}

	response := new(ProductResponse)
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return err
	}

	if res.StatusCode != response.Status {
		return fmt.Errorf("get products http request failed(body). status: %d, response: %v", response.Status, res.Body)
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("status code: %d", res.StatusCode)
		return fmt.Errorf(utility.ErrNotOk)
	}

	err = p.fararavand.SyncProductsWithGoods(response.NewProducts)
	if err != nil {
		return fmt.Errorf("load SyncProductsWithGoods encountered an error: %w", err)
	}

	return nil
}
