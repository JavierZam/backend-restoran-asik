package repository

import (
	"database/sql"
	"fmt"
	"log"
	"restoran-asik/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(user models.User) error {
	query := "INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id"
	row := repo.DB.QueryRow(query, user.Username, user.Password, user.Role)

	if err := row.Scan(&user.ID); err != nil {
		log.Println(err)
		return fmt.Errorf("Failed to create user: %v", err)
	}

	return nil
}

func (repo *UserRepository) GetUserByUsernamePassword(username string, password string) (*models.User, error) {
	query := "SELECT id, username, role FROM users WHERE username = $1 AND password = $2"
	row := repo.DB.QueryRow(query, username, password)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetUserByID(userID int) (*models.User, error) {
	query := "SELECT id, username, role FROM users WHERE id = $1"
	row := repo.DB.QueryRow(query, userID)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
