package syncdata

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

// BaseDataResponse is the response for the BaseData
type BaseDataResponse struct {
	Status      int             `json:"status"`
	NewBaseData models.BaseData `json:"new_base_data"`
}

type BaseData struct {
	log        *zap.SugaredLogger
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewBaseData(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface) *BaseData {
	return &BaseData{
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
			b.log.Errorw("get base data request encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)

		res, err := b.httpClient.Do(req)
		if err != nil {
			b.log.Errorw("get base data response encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)
			return err
		}

		if res.StatusCode != http.StatusOK {
			b.log.Errorw("get base data http request failed.",
				"error", err,
				"status:", res.StatusCode,
				"response", res.Body,
			)

			return fmt.Errorf("get invoice http request failed. status: %d, response: %v", res.StatusCode, res.Body)
		}

		response := new(BaseDataResponse)
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			b.log.Errorw("get base data decode response encountered an error: ",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		if len(response.NewBaseData.PaymentTypes) == 0 {
			break
		}

		err = b.fararavand.SyncBaseDataWithDeliverCenter(response.NewBaseData)
		if err != nil {
			b.log.Errorw("load SyncBaseDataWithDeliverCenter encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}

		lastId = pageNumber + len(response.NewBaseData.PaymentTypes)

		err = b.repos.Database.InsertBaseDataProgress(lastId, pageNumber)
		if err != nil {
			b.log.Errorw("InsertBaseDataProgress encountered an error:",
				"error", err,
				"last_id", lastId,
				"page_number", pageNumber,
			)

			return err
		}
		if len(response.NewBaseData.PaymentTypes) < pageSize {
			break
		}

	}

	return nil
}
