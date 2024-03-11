package syncdata

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

// BaseDataResponse is the response for the BaseData
type BaseDataResponse struct {
	Status      int             `json:"status"`
	NewBaseData models.BaseData `json:"new_base_data"`
}

type BaseData struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewBaseData(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *BaseData {
	return &BaseData{
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: config.Cfg.FararavandApp.Timeout,
		},
	}
}

func (b *BaseData) BaseData() error {
	var lastId int
	var pageNumber int
	var pageSize int = 1000

	for {

		lastBaseDataId, lastPageNumber, err := b.repos.Database.GetBaseDataProgress()
		if err == sql.ErrNoRows {
			lastId = 0
			pageNumber = 0
		} else {
			lastId = lastBaseDataId
			pageNumber = lastPageNumber + pageSize + 1
		}

		req, err := http.NewRequest(http.MethodGet, b.baseURL+
			fmt.Sprintf(GetBaseData, pageNumber, pageSize, lastId), nil)
		if err != nil {
			return err
		}

		req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

		res, err := b.httpClient.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("get invoice http request failed. status: %d, response: %v", res.StatusCode, res.Body)
		}

		response := new(BaseDataResponse)
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			return err
		}

		if len(response.NewBaseData.PaymentTypes) == 0 {
			break
		}

		if res.StatusCode != response.Status {
			return fmt.Errorf("get base data http request failed(body). status: %d, response: %v", response.Status, res.Body)
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("status code: %d", res.StatusCode)
			return fmt.Errorf(ErrNotOk)
		}

		err = b.fararavand.SyncBaseDataWithDeliverCenter(response.NewBaseData)
		if err != nil {
			fmt.Println("load SyncBaseDataWithDeliverCenter encountered an error: %w", err)
			return err
		}

		lastId = pageNumber + len(response.NewBaseData.PaymentTypes)

		err = b.repos.Database.InsertBaseDataProgress(lastId, pageNumber)
		if err != nil {
			return err
		}
		if len(response.NewBaseData.PaymentTypes) < pageSize {
			break
		}

	}

	return nil
}
