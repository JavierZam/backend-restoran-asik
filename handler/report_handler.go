package handler

import (
	"log"
	"net/http"
	"restoran-asik/models"
	"restoran-asik/repository"
	"time"

	"github.com/labstack/echo/v4"
)

type ReportService struct {
	ReportRepo      repository.ReportRepository
	TransactionRepo repository.TransactionRepository
	ProductRepo     repository.ProductRepository
}

func NewReportHandler(reportRepo repository.ReportRepository, productRepo repository.ProductRepository, transactionRepo repository.TransactionRepository) *ReportService {
	return &ReportService{
		ReportRepo:      reportRepo,
		TransactionRepo: transactionRepo,
		ProductRepo:     productRepo,
	}
}

func (rs *ReportService) CreateReport(c echo.Context) error {
	var request models.ReportRequest
	if err := c.Bind(&request); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	today := time.Now().Format("2006-01-02")
	transactions, err := rs.TransactionRepo.GetTransactionsByDateAndPaidStatus(today, true)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get transactions"})
	}

	totalSales := 0.0
	productsSold := make(map[string]int)

	for _, transaction := range transactions {
		totalSales += transaction.TotalAmount

		product, err := rs.ProductRepo.GetProductByID(transaction.ProductID)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get product information"})
		}

		productsSold[product.Name] += transaction.Quantity
	}

	request = models.ReportRequest{
		TotalSales:   totalSales,
		ReportDate:   today,
		ProductsSold: productsSold,
	}

	err = rs.ReportRepo.CreateReport(request)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create report"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Report created successfully"})
}

func (rs *ReportService) GetReportByDate(c echo.Context) error {
	reportDate := c.Param("date")

	report, err := rs.ReportRepo.GetReportByDate(reportDate)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get report by date"})
	}

	log.Printf("Report: %+v", report)

	return c.JSON(http.StatusOK, report)
}
