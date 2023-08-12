package aryan

import (
	"erp-job/config"
	"erp-job/logics"
	"erp-job/models/aryan"
	"erp-job/models/fararavand"
	"erp-job/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type AryanInterface interface {
	GetProducts() ([]aryan.SaleCustomer, error)
}

type Aryan struct {
	restyClient *resty.Client
	baseUrl     string
	repos       *repository.Repository
}

func NewAryan() AryanInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.ApiKey)

	return &Aryan{
		restyClient: c,
		baseUrl:     config.Cfg.BaseURL,
	}
}

func NewLogics(repos *repository.Repository) *Aryan {
	return &Aryan{
		repos: repos,
	}
}

// GetProducts gets all products data from the first ERP
func (a *Aryan) GetProducts() ([]aryan.SaleCustomer, error) {
	var newProducts []fararavand.Products

	resp, err := a.restyClient.R().
		SetResult(newProducts).
		Get(
			fmt.Sprintf("%s/%s", a.baseUrl, logics.FGetProducts),
		)
	if err != nil {
		return nil, err
	}

	// get last product id from response --100
	lastId := newProducts[len(newProducts)-1].ID
	// get last product id from data --80
	pId, err := f.repos.Database.GetProduct()
	if err != nil {
		return nil, err
	}
	// fetch new product id
	if lastId > pId {
		newProducts = newProducts[pId:]
		//insert new product id into db
		//
		err = f.repos.Database.InsertProduct(lastId)
		if err != nil {
			return nil, err
		}
		return newProducts, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code: %d", resp.StatusCode())
		return nil, fmt.Errorf(logics.ErrNotOk)
	}

	return newProducts, nil
}
