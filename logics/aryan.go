package logics

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type AryanInterface interface {
	PostSaleOrder(fp []models.FararavandProducts) (*resty.Response, error)
}

type Aryan struct {
	restyClient *resty.Client
	baseUrl     string
	repos       *repository.Repository
}

func NewAryan(repos *repository.Repository) AryanInterface {
	c := resty.New().
		SetHeader("ApiKey", config.Cfg.ApiKey)

	return &Aryan{
		restyClient: c,
		baseUrl:     config.Cfg.BaseURL,
		repos:       repos,
	}
}

// PostSalesOrder Post all sale order data to the secound ERP
func (a *Aryan) PostSaleOrder(fp []models.FararavandProducts) (*resty.Response, error) {
	var newGoods []models.AryanGoods

	for _, item := range fp {
		newGoods = append(newGoods, models.AryanGoods{
			Level1: item.BrandID,
			// TODO: fix this
		})
	}

	res, err := a.restyClient.R().SetBody(newGoods).Post("asdasdasdasd")
	if err != nil {
		// TOOD: handle erro
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostSaleCustomer Post all sale customer data to the secound ERP
func (a *Aryan) PostSaleCustomer(fp []models.FararavandProducts) (*resty.Response, error) {
	var newGoods []models.AryanGoods

	for _, item := range fp {
		newGoods = append(newGoods, models.AryanGoods{
			Level1: item.BrandID,
			// TODO: fix this
		})
	}

	res, err := a.restyClient.R().SetBody(newGoods).Post("asdasdasdasd")
	if err != nil {
		// TOOD: handle erro
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostSaleCustomer Post all sale customer data to the secound ERP
func (a *Aryan) PostSaleTypeSelect(fp []models.FararavandProducts) (*resty.Response, error) {
	var newGoods []models.AryanGoods

	for _, item := range fp {
		newGoods = append(newGoods, models.AryanGoods{
			Level1: item.BrandID,
			// TODO: fix this
		})
	}

	res, err := a.restyClient.R().SetBody(newGoods).Post("asdasdasdasd")
	if err != nil {
		// TOOD: handle erro
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}
