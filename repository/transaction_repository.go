package repository

import (
	"database/sql"
	"fmt"
	"log"
	"restoran-asik/models"
)

type TransactionRepository struct {
	DB          *sql.DB
	ProductRepo *ProductRepository
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		DB:          db,
		ProductRepo: NewProductRepository(db),
	}
}

func (repo *TransactionRepository) CreateTransaction(newTransaction models.Transaction) (int, error) {
	product, err := repo.ProductRepo.GetProductByID(newTransaction.ProductID)

	newTransaction.TotalAmount = float64(newTransaction.Quantity) * product.Price

	if newTransaction.PaymentType == "tunai" {
		newTransaction.IsPaid = false
	} else {
		newTransaction.IsPaid = true
	}

	var transactionID int
	err = repo.DB.QueryRow("INSERT INTO transactions (product_id, user_id, quantity, payment_type, total_amount, is_paid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		newTransaction.ProductID,
		newTransaction.UserID,
		newTransaction.Quantity,
		newTransaction.PaymentType,
		newTransaction.TotalAmount,
		newTransaction.IsPaid,
	).Scan(&transactionID)

	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("Failed to create transaction: %v", err)
	}

	return transactionID, nil
}

func (repo *TransactionRepository) GetUnpaidTransactions() ([]models.Transaction, error) {
	query := "SELECT id, product_id, user_id, quantity, payment_type, total_amount, is_paid FROM transactions WHERE is_paid = false"
	rows, err := repo.DB.Query(query)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to get unpaid transactions: %v", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.ProductID, &transaction.UserID, &transaction.Quantity, &transaction.PaymentType, &transaction.TotalAmount, &transaction.IsPaid)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Failed to scan transaction row: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repo *TransactionRepository) UpdateTransactionIsPaid(transactionID int) error {
	query := "UPDATE transactions SET is_paid = true WHERE id = $1"
	_, err := repo.DB.Exec(query, transactionID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to update transaction is_paid: %v", err)
	}

	return nil
}

func (repo *TransactionRepository) GetTransactionByID(id int) (*models.Transaction, error) {
	query := "SELECT id, product_id, user_id, quantity, payment_type, total_amount, is_paid FROM transactions WHERE id = $1"
	row := repo.DB.QueryRow(query, id)

	var transaction models.Transaction
	err := row.Scan(&transaction.ID, &transaction.ProductID, &transaction.UserID, &transaction.Quantity, &transaction.PaymentType, &transaction.TotalAmount, &transaction.IsPaid)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to get transaction by ID: %v", err)
	}

	return &transaction, nil
}

func (repo *TransactionRepository) GetAllTransactions() ([]models.Transaction, error) {
	query := "SELECT id, product_id, user_id, quantity, payment_type, total_amount, is_paid FROM transactions"
	rows, err := repo.DB.Query(query)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to get transactions: %v", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.ProductID, &transaction.UserID, &transaction.Quantity, &transaction.PaymentType, &transaction.TotalAmount, &transaction.IsPaid)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Failed to scan transaction row: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repo *TransactionRepository) GetTransactionsByDateAndPaidStatus(date string, isPaid bool) ([]models.Transaction, error) {
	query := "SELECT id, product_id, user_id, quantity, payment_type, total_amount, is_paid, transaction_date FROM transactions WHERE DATE(transaction_date) = $1 AND is_paid = $2"
	rows, err := repo.DB.Query(query, date, isPaid)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to get transactions: %v", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.ProductID, &transaction.UserID, &transaction.Quantity, &transaction.PaymentType, &transaction.TotalAmount, &transaction.IsPaid, &transaction.TransactionDate)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Failed to scan transaction row: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
