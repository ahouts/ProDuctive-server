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

type CreateUserRequest struct {
	Email    string
	Password string
}

const bcryptCost = 12

func (c *Conn) GetUser(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("user-id")
	var u User
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, user id %v is invalid.\n%v", idStr, err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	ctx := InitContext()
	defer ctx.Done()
	err = c.QueryRowContext(ctx, "SELECT id, email, password_hash, created_at, updated_at FROM user_profile WHERE id = :1", id).Scan(&u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, user id %v is invalid.\n%v", idStr, err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	response.WriteEntity(u)
}

func (c *Conn) CreateUser(request *restful.Request, response *restful.Response) {
	userRequest := CreateUserRequest{}
	err := request.ReadEntity(&userRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Request invalid, request must match format %v.", userRequest))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcryptCost)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Request invalid, failed to hash password %v.", userRequest.Password))
		return
	}

	tx, err := c.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Failed to initialize a transaction.\n%v", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	_, err = tx.Exec("INSERT INTO user_profile VALUES(null, :1, utl_raw.cast_to_raw(:2), default, default)", userRequest.Email, string(hashedPassword))
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Failed to create user %v, %v\nerr: %v", userRequest.Email, userRequest.Password, err))
		tx.Rollback()
		return
	}
	tx.Commit()
}
