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
	"erp-job/utility/logger"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type ProductResponse struct {
	Status      int               `json:"status"`
	NewProducts []models.Products `json:"new_products"`
}

type Product struct {
	log        *zap.SugaredLogger
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewProduct(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *Product {

	return &Product{
		log:        logger.Logger(),
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
			p.log.Errorw("get product request encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

		res, err := p.httpClient.Do(req)
		if err != nil {
			p.log.Errorw("get product response encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if res.StatusCode != http.StatusOK {
			p.log.Errorw("get product http request failed.",
				"error", err,
				"status:", res.StatusCode,
				"response", res.Body,
			)

			return fmt.Errorf("get product http request failed. status: %d, response: %v", res.StatusCode, res.Body)
		}

		response := new(ProductResponse)
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			p.log.Errorw("get product decode response encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if len(response.NewProducts) == 0 {
			break
		}

		err = p.fararavand.SyncProductsWithGoods(response.NewProducts)
		if err != nil {
			p.log.Errorw("load SyncProductsWithGoods encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		lastId = pageNumber + len(response.NewProducts)

		err = p.repos.Database.InsertProductProgress(lastId, pageNumber)
		if err != nil {
			p.log.Errorw("InsertProductProgress encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if len(response.NewProducts) < pageSize {
			break
		}

	}

	return nil
}
