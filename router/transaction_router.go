package router

import (
	"restoran-asik/handler"

	"github.com/labstack/echo/v4"
)

func InitTransactionRoutes(e *echo.Echo, h *handler.TransactionService) {
	e.POST("/transaction", h.CreateTransaction)
	e.GET("/unpaid-transactions", h.GetUnpaidTransactions)
	e.POST("/confirm-payment/:id", h.ConfirmPayment)
	e.GET("/transaction/:id", h.GetTransactionByID)
	e.GET("/transaction", h.GetAllTransactions)
}
