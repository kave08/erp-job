package syncdata

//rename pkg name

import (
	"database/sql"
	"encoding/json"
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"fmt"
	"log"
	"net/http"
)

type ProductResponse struct {
	Status      int               `json:"status"`
	NewProducts []models.Products `json:"new_products"`
}

type Product struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewProduct(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Product {

	return &Product{
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (p Product) Products() error {
	var lastId int
	var pageNumber int
	var pageSize int = 1000

	for {
		lastProductId, lastPageNumber, err := p.repos.Database.GetProductProgress()
		if err == sql.ErrNoRows {
			lastId = 0
			pageNumber = 0
		} else {
			lastId = lastProductId
			pageNumber = lastPageNumber + pageSize + 1
		}

		req, err := http.NewRequest(http.MethodGet, p.baseURL+
			fmt.Sprintf(GetProducts, pageNumber, pageSize, lastId), nil)
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

		if len(response.NewProducts) == 0 {
			break
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("status code: %d", res.StatusCode)
			return fmt.Errorf(ErrNotOk)
		}

		err = p.fararavand.SyncProductsWithGoods(response.NewProducts)
		if err != nil {
			return fmt.Errorf("load SyncProductsWithGoods encountered an error: %w", err)
		}

		lastId = pageNumber + len(response.NewProducts)

		err = p.repos.Database.InsertProductProgress(lastId, pageNumber)
		if err != nil {
			return err
		}

		if len(response.NewProducts) < pageSize {
			break
		}

	}

	return nil
}
