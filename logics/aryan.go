package logics

import (
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type AryanInterface interface {
	PostSaleFactor(fp []models.Fararavand) (*resty.Response, error)
	PostSaleCustomer(fp []models.Fararavand) (*resty.Response, error)
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
func (a *Aryan) PostSaleFactor(fp []models.Fararavand) (*resty.Response, error) {
	var newSaleFactor models.Aryan

	for _, item := range fp {
		newSaleFactor = append(newSaleFactor, models.Aryan{
			SaleFactor{
				CustomerId: item.CustomerId,
			},
		})
	}
	
	res, err := a.restyClient.R().SetBody(newSaleFactor).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}

// PostSaleCustomer Post all sale customer data to the secound ERP
func (a *Aryan) PostSaleCustomer(fp []models.FararavandCustomers) (*resty.Response, error) {
	var newSaleCustomer []models.AryanSaleCustomer

	for _, item := range fp {
		newSaleCustomer = append(newSaleCustomer, models.AryanSaleCustomer{
			CustomerID:   item.CustomerId,
			CustomerCode: strconv.Itoa(item.CustomerCodePosti),
		})
	}

	res, err := a.restyClient.R().SetBody(newSaleCustomer).Post("asdasdasdasd")
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		fmt.Println(res.Body())
	}

	return res, nil
}
