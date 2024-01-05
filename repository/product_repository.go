package repository

import (
	"database/sql"
	"fmt"
	"log"
	"restoran-asik/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (repo *ProductRepository) GetAllProducts() ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM products"
	rows, err := repo.DB.Query(query)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to get products: %v", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Failed to scan product row: %v", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (repo *ProductRepository) GetProductByID(id int) (*models.Product, error) {
	query := "SELECT id, name, price, stock FROM products WHERE id = $1"
	row := repo.DB.QueryRow(query, id)

	var product models.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Failed to get product by ID: %v", err)
	}

	return &product, nil
}

func (repo *ProductRepository) UpdateProduct(id int, product models.Product) error {
	query := "UPDATE products SET name=$1, price=$2, stock=$3 WHERE id=$4"
	_, err := repo.DB.Exec(query, product.Name, product.Price, product.Stock, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to update product: %v", err)
	}

	return nil
}

func (repo *ProductRepository) AddProduct(product models.Product) error {
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err := repo.DB.QueryRow(query, product.Name, product.Price, product.Stock).Scan(&product.ID)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to add product: %v", err)
	}

	return nil
}

func (repo *ProductRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := repo.DB.Exec(query, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to delete product: %v", err)
	}

	return nil
}
