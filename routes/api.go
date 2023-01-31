package api

import (
	"github.com/labstack/echo/v4"
)


func SetRouting(e echo.Echo) error{

	g := e.Group("users")
	g.GET("",nil)

	return nil
}