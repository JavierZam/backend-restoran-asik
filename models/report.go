// models/report.go
package models

import "encoding/json"

type Report struct {
	ID           int             `json:"id"`
	TotalSales   float64         `json:"total_sales"`
	ReportDate   string          `json:"report_date"`
	ProductsSold json.RawMessage `json:"products_sold"`
}

// ReportRequest represents the request structure for creating a report.
type ReportRequest struct {
	TotalSales   float64        `json:"total_sales"`
	ReportDate   string         `json:"report_date"`
	ProductsSold map[string]int `json:"products_sold"`
}

// ReportResponse represents the response structure for getting a report.
type ReportResponse struct {
	TotalSales   float64        `json:"total_sales"`
	ReportDate   string         `json:"report_date"`
	ProductsSold map[string]int `json:"products_sold"`
}
