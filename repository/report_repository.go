// repository/report_repository.go
package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"restoran-asik/models"
	"time"
)

type ReportRepository struct {
	DB *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (repo *ReportRepository) CreateReport(request models.ReportRequest) error {
	productsSoldJSON, err := json.Marshal(request.ProductsSold)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to marshal products_sold: %v", err)
	}

	query := "INSERT INTO reports (total_sales, report_date, products_sold) VALUES ($1, $2, $3) RETURNING id"
	_, err = repo.DB.Exec(query, request.TotalSales, request.ReportDate, productsSoldJSON)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to add report: %v", err)
	}

	return nil
}

func (repo *ReportRepository) GetReportByDate(reportDate string) (*models.ReportResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", reportDate)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to parse date: %v", err)
	}

	formattedDate := parsedDate.Format("2006-01-02")

	query := "SELECT total_sales, report_date, products_sold FROM reports WHERE report_date = $1 ORDER BY id DESC LIMIT 1"
	row := repo.DB.QueryRow(query, formattedDate)

	var response models.ReportResponse
	var productsSoldJSON []byte
	err = row.Scan(&response.TotalSales, &response.ReportDate, &productsSoldJSON)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to get report by date: %v", err)
	}

	err = json.Unmarshal(productsSoldJSON, &response.ProductsSold)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to unmarshal products_sold: %v", err)
	}

	return &response, nil
}
