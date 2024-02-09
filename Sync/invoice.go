package sync

import (
	"context"
	"erp-job/config"
	"erp-job/models"
	"erp-job/repository"
	"erp-job/services/aryan"
	"erp-job/services/fararavand"
	"log"
	"net/http"
	"strconv"
	"time"

	"erp-job/helpers"

	"github.com/labstack/echo/v4"
)

const (
	handleGetInvoiceEndpoint = "GetInvoices"
)

type InvoiceRequest struct {
	LatestId int `json:"latest_id"`
	Invoices []models.Invoices
}

// InvoiceResponse is the response for the invoice
type InvoiceResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewInvoiceResponse is the InvoiceResponse factory method
func NewInvoiceResponse(status int, message string) *InvoiceResponse {
	return &InvoiceResponse{
		Status:  status,
		Message: message,
	}
}

type Invoice struct {
	e          *echo.Echo
	baseURL    string
	httpClient *http.Client
	repos      *repository.Repository
	aryan      aryan.AryanInterface
	fararavand fararavand.FararavandInterface
}

func NewInvoice(repos *repository.Repository, fr fararavand.FararavandInterface, ar aryan.AryanInterface, requestTimeout time.Duration) *Invoice {
	e := echo.New()
	e.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				req.URL.Scheme = "https"
				req.URL.Host = config.Cfg.FararavandApp.BaseURL
				req.Header.Set("ApiKey", config.Cfg.FararavandApp.APIKey)
				return next(c)
			}
		},
	)
	return &Invoice{
		e:          e,
		baseURL:    config.Cfg.FararavandApp.BaseURL,
		repos:      repos,
		aryan:      ar,
		fararavand: fr,
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (cl *client) Invoices(ctx context.Context, inoviceId int) func(c echo.Context) error {
	return func(c echo.Context) error {

		url := helpers.EncodeURLParams(cl.baseURL+handleGetInvoiceEndpoint, map[string]string{
			"inovice_id": strconv.Itoa(inoviceId),
		})

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}

		response := new(InvoiceResponse)

		err = cl.call(req, response)

		request := new(InvoiceRequest)

		if err := c.Bind(request); err != nil {
			log.Printf(
				"@ERP.sync.invoice.Invoices",
				"message", "bind invoice",
				"error", err,
			)
			return c.JSON(http.StatusBadRequest, NewInvoiceResponse(http.StatusBadRequest, "validation.required"))
		}

		if request.LatestId <= 0 {
			log.Printf(
				"@ERP.sync.invoice.Invoices",
				"message", "invalid payload",
			)

			return c.JSON(http.StatusBadRequest, NewInvoiceResponse(http.StatusBadRequest, "validation.required"))
		}

		// i.e.GET(handleGetInvoiceEndpoint+"/:id", func(c echo.Context) error {
		// 	id := c.Param("latest_id")
		// 	return c.String(http.StatusOK, "User ID: "+id)
		// })

		// if i.e.AcquireContext().Response().Status != http.StatusOK {
		// 	log.Printf("status code: %d", i.e.AcquireContext().Response().Status)
		// 	return fmt.Errorf(utility.ErrNotOk)
		// }

		// err = i.fararavand.SyncInvoicesWithSaleFactor(request.Invoices)
		// if err != nil {
		// 	fmt.Println("Load SyncInvoicesWithSaleFactor encountered an error", err.Error())
		// 	return err
		// }

		// err = i.fararavand.SyncInvoicesWithSaleOrder(request.Invoices)
		// if err != nil {
		// 	fmt.Println("Load SyncInvoicesWithSaleOrder encountered an error", err.Error())
		// 	return err
		// }

		// err = i.fararavand.SyncInvoicesWithSalePayment(request.Invoices)
		// if err != nil {
		// 	fmt.Println("Load SyncInvoicesWithSalePayment encountered an error", err.Error())
		// 	return err
		// }

		// err = i.fararavand.SyncInvoicesWithSalerSelect(request.Invoices)
		// if err != nil {
		// 	fmt.Println("Load SyncInvoicesWithSalerSelect encountered an error", err.Error())
		// 	return err
		// }

		// err = i.fararavand.SyncInvoicesWithSaleProforma(request.Invoices)
		// if err != nil {
		// 	fmt.Println("Load SyncInvoicesWithSaleProforma encountered an error", err.Error())
		// 	return err
		// }

		// err = i.fararavand.SyncInvoicesWithSaleCenter(request.Invoices)
		// if err != nil {
		// 	fmt.Println("Load SyncInvoicesWithSaleCenter encountered an error", err.Error())
		// 	return err
		// }

		// err = i.fararavand.SyncInvoiceWithSaleTypeSelect(request.Invoices)
		// if err != nil {
		// 	fmt.Println("Load SyncInvoicesWithSaleCenter encountered an error", err.Error())
		// 	return err
		// }

		return nil
	}
}
