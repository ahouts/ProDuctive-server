package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
)

type User struct {
	Id           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (User) SwaggerDoc() map[string]string {
	return map[string]string{
		"":             "User data",
		"Id":           "user id",
		"Email":        "user email, unique",
		"PasswordHash": "hashed password using bcrypt",
		"CreatedAt":    "timestamp the user was created at",
		"UpdatedAt":    "timestamp the user was last updated",
	}
}

func (s *DbSession) GetUser(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("user-id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, user id %v is invalid.\n%v", idStr, err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Failed to initialize a transaction.\n%v", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	var u User
	err = tx.QueryRow("SELECT id, email, password_hash, created_at, updated_at FROM user_profile WHERE id = :1", id).Scan(&u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, user id %v is invalid.\n%v", idStr, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
	response.WriteEntity(u)
	tx.Commit()
}

type CreateUserRequest struct {
	Email    string
	Password string
}

func (CreateUserRequest) SwaggerDoc() map[string]string {
	return map[string]string{
		"":         "Form to create a new user",
		"Email":    "user email, must be unique",
		"Password": "User password. Supports infinite length but only the first 72 characters will be used",
	}
}

type UserCreated struct {
	Email string
	Id int
}

const bcryptCost = 12

func (s *DbSession) CreateUser(request *restful.Request, response *restful.Response) {
	userRequest := CreateUserRequest{}
	err := request.ReadEntity(&userRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v.", userRequest))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcryptCost)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, failed to hash password %v.", userRequest.Password))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to initialize a transaction: %v", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	insertUser, err := tx.Prepare("INSERT INTO user_profile VALUES(null, :1, utl_raw.cast_to_raw(:2), default, default)")
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to prepare insert command for %v: %v", userRequest.Email, err))
		tx.Rollback()
		return
	}
	_, err = insertUser.Exec(userRequest.Email, string(hashedPassword))
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to insert user %v: %v", userRequest.Email, err))
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		return
	}
}

// returns id if successful, error otherwise
func AuthUser(tx *sql.Tx, email, password string) (int, error) {
	var u User
	err := tx.QueryRow("SELECT id, password_hash FROM user_profile WHERE email = :1", email).Scan(&u.Id, &u.PasswordHash)
	if err != nil {
		log.Println(errors.New(err).ErrorStack())
		return 0, fmt.Errorf("invalid query, user email %v is invalid: %v", email, err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return 0, err
	}
	return u.Id, nil
}
