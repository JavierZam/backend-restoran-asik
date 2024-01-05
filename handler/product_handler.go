package handler

import (
	"log"
	"net/http"
	"restoran-asik/models"
	"restoran-asik/repository"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductService struct {
	ProductRepo repository.ProductRepository
}

func NewProductHandler(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepo: productRepo,
	}
}

func (h *ProductService) GetProducts(c echo.Context) error {
	products, err := h.ProductRepo.GetAllProducts()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get products"})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *ProductService) GetProductByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	product, err := h.ProductRepo.GetProductByID(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductService) UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var updatedProduct models.Product
	if err := c.Bind(&updatedProduct); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	err = h.ProductRepo.UpdateProduct(id, updatedProduct)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Product updated successfully"})
}

func (h *ProductService) AddProduct(c echo.Context) error {
	var newProduct models.Product
	if err := c.Bind(&newProduct); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	err := h.ProductRepo.AddProduct(newProduct)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add product"})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Product added successfully"})
}

func (h *ProductService) DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	err = h.ProductRepo.DeleteProduct(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
