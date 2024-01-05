package router

import (
	"restoran-asik/handler"

	"github.com/labstack/echo/v4"
)

func InitReportRoutes(e *echo.Echo, h *handler.ReportService) {
	e.POST("/report", h.CreateReport)
	e.GET("/report/:date", h.GetReportByDate)
}
