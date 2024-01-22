package handlers

import "github.com/labstack/echo/v4"

func MapProductRoutes(productRouteGroup *echo.Group, p ProductHandlers) {
	productRouteGroup.POST("", p.Create())
	productRouteGroup.PUT("/:id", p.Update())
	productRouteGroup.DELETE("/:id", p.Delete())
	productRouteGroup.GET("/:id", p.GetById())
	productRouteGroup.GET("", p.GetAll())
}
