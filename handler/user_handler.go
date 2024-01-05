package handler

import (
	"log"
	"net/http"
	"restoran-asik/models"
	"restoran-asik/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserService struct {
	UserRepo  repository.UserRepository
	JWTSecret string
}

func NewUserHandler(userRepo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

func (us *UserService) RegisterUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	user.Role = "customer"

	err := us.UserRepo.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (us *UserService) Login(c echo.Context) error {
	var loginRequest models.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	log.Println("Login Attempt for User:", loginRequest.Username)

	user, err := us.UserRepo.GetUserByUsernamePassword(loginRequest.Username, loginRequest.Password)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials - user not found or password does not match"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(us.JWTSecret))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate JWT token"})
	}

	response := map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (us *UserService) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user"})
	}

	user, err := us.UserRepo.GetUserByID(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}
