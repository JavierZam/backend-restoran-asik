package handler

import (
	"log"
	"net/http"
	"restoran-asik/models"
	"restoran-asik/repository"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionService struct {
	TransactionRepo repository.TransactionRepository
	ProductRepo     repository.ProductRepository
	UserRepo        repository.UserRepository
}

func NewTransactionHandler(transactionRepo repository.TransactionRepository, productRepo repository.ProductRepository, userRepo repository.UserRepository) *TransactionService {
	return &TransactionService{
		TransactionRepo: transactionRepo,
		ProductRepo:     productRepo,
		UserRepo:        userRepo,
	}
}

func (ts *TransactionService) CreateTransaction(c echo.Context) error {
	var newTransaction models.Transaction
	if err := c.Bind(&newTransaction); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	user, err := ts.UserRepo.GetUserByID(newTransaction.UserID)
	if err != nil || user.Role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	product, err := ts.ProductRepo.GetProductByID(newTransaction.ProductID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to get product information"})
	}

	newTransaction.TotalAmount = float64(newTransaction.Quantity) * product.Price

	transactionID, err := ts.TransactionRepo.CreateTransaction(newTransaction)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create transaction"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Transaction created successfully", "id": strconv.Itoa(transactionID)})
}

func (ts *TransactionService) GetUnpaidTransactions(c echo.Context) error {
	transactions, err := ts.TransactionRepo.GetUnpaidTransactions()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get unpaid transactions"})
	}

	return c.JSON(http.StatusOK, transactions)
}

func (ts *TransactionService) ConfirmPayment(c echo.Context) error {
	transactionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transaction ID"})
	}

	// Call repository function to update is_paid to true
	err = ts.TransactionRepo.UpdateTransactionIsPaid(transactionID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to confirm payment"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Payment confirmed successfully"})
}

func (ts *TransactionService) GetTransactionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transaction ID"})
	}

	transaction, err := ts.TransactionRepo.GetTransactionByID(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get transaction by ID"})
	}

	return c.JSON(http.StatusOK, transaction)
}

func (ts *TransactionService) GetAllTransactions(c echo.Context) error {
	transactions, err := ts.TransactionRepo.GetAllTransactions()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get transactions"})
	}

	return c.JSON(http.StatusOK, transactions)
}
