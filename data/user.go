package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
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

type GetUserRequest struct {
	Email    string
	Password string
}

func (s *DbSession) GetUser(request *restful.Request, response *restful.Response) {
	userRequest := GetUserRequest{}
	err := request.ReadEntity(&userRequest)
	if err != nil {
		formatError(new(GetUserRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Failed to initialize a transaction.\n%v", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, userRequest.Email, userRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	var u User
	err = tx.QueryRow("SELECT id, email, password_hash, created_at, updated_at FROM user_profile WHERE id = ?", userId).Scan(&u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Invalid query, user id %v is invalid.\n%v", userId, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
	u.PasswordHash = "redacted"
	response.WriteEntity(u)
	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		log.Println(errors.New(err).ErrorStack())
		return
	}
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

const bcryptCost = 12

func (s *DbSession) CreateUser(request *restful.Request, response *restful.Response) {
	userRequest := CreateUserRequest{}
	err := request.ReadEntity(&userRequest)
	if err != nil {
		formatError(new(CreateUserRequest), response)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcryptCost)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, failed to hash password."))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize a transaction: %v", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	insertUser, err := tx.Prepare("INSERT INTO user_profile VALUES(NULL, ?, utl_raw.cast_to_raw(?), default, default)")
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to prepare insert command for %v: %v", userRequest.Email, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
	_, err = insertUser.Exec(userRequest.Email, string(hashedPassword))
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to insert user %v: %v", userRequest.Email, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
}

type GetUserIdRequest struct {
	Email string
}

func (GetUserIdRequest) SwaggerDoc() map[string]string {
	return map[string]string{
		"":      "form to get user's id",
		"Email": "user to get the id of",
	}
}

type UserId struct {
	Id int
}

func (UserId) SwaggerDoc() map[string]string {
	return map[string]string{
		"":   "user's id",
		"Id": "user's id",
	}
}

func (s *DbSession) GetUserId(request *restful.Request, response *restful.Response) {
	userIdReq := GetUserIdRequest{}
	err := request.ReadEntity(&userIdReq)
	if err != nil {
		formatError(new(GetUserIdRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize a transaction: %v", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	uid := UserId{}
	err = tx.QueryRow("SELECT id FROM user_profile WHERE email = ?", userIdReq.Email).Scan(&uid.Id)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to find user with email %v: %v", userIdReq.Email, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
	response.WriteEntity(uid)
	tx.Commit()
}
