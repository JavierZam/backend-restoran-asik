package router

import (
	"restoran-asik/handler"

	"github.com/labstack/echo/v4"
)

func InitUserRoutes(e *echo.Echo, h *handler.UserService) {
	e.POST("/user/register", h.RegisterUser)
	e.POST("/user/login", h.Login)
	e.GET("/user/:id", h.GetUserByID)
}
