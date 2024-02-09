package sync

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Contract interface {
	Invoices(ctx context.Context, inoviceId int) func(c echo.Context) error
}
