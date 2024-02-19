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

// BaseDataResponse is the response for the BaseData
type BaseDataResponse struct {
	Status      int             `json:"status"`
	NewBaseData models.BaseData `json:"new_base_data"`
}

type BaseDataRequest struct {
	LastId      int `json:"LastId"`
	PageSize    int `json:"PageSize"`
	PageNumeber int `json:"PageNumeber"`
}

// NewBaseDataRequest is the BaseDataResponse factory method
func NewBaseDataRequest(lastid int, pageSize int, pageNumber int) ProductRequest {
	return ProductRequest{
		LastId:      lastid,
		PageSize:    pageSize,
		PageNumeber: pageNumber,
	}
}

type BaseData struct {
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewBaseData(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface, requestTimeout time.Duration) *BaseData {
	return &BaseData{
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (b *BaseData) BaseData() error {

	request := new(BaseDataRequest)

	req, err := http.NewRequest(http.MethodGet, b.baseURL+
		fmt.Sprintf("/GetBaseData?PageNumeber=%d&PageSize=%d&LastId=%d/", request.PageNumeber, request.PageSize, request.LastId), nil)
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

	if res.StatusCode != response.Status {
		return fmt.Errorf("get base data http request failed(body). status: %d, response: %v", response.Status, res.Body)
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("status code: %d", res.StatusCode)
		return fmt.Errorf(utility.ErrNotOk)
	}

	if request.LastId <= 0 {
		return fmt.Errorf("validation.required %d", http.StatusBadRequest)
	}

	err = b.fararavand.SyncBaseDataWithDeliverCenter(response.NewBaseData)
	if err != nil {
		fmt.Println("load SyncBaseDataWithDeliverCenter encountered an error: %w", err)
		return err
	}

	return nil
}
