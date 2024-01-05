package router

import (
	"restoran-asik/handler"

	"github.com/labstack/echo/v4"
)

func InitProductRoutes(e *echo.Echo, h *handler.ProductService) {
	e.GET("/product", h.GetProducts)
	e.GET("/product/:id", h.GetProductByID)
	e.PUT("/product/:id", h.UpdateProduct)
	e.POST("/product", h.AddProduct)
	e.DELETE("/product/:id", h.DeleteProduct)
}
