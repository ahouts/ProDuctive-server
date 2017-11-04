package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

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

type CreateUserRequest struct {
	Email    string
	Password string
}

const bcryptCost = 12

func (db *Conn) GetUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	query := fmt.Sprintf("SELECT id, email, password_hash, created_at, updated_at FROM user_profile WHERE id=%v", id)
	ctx := InitContext()
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalln(errors.New(err).ErrorStack())
	}
	defer rows.Close()
	defer ctx.Done()
	users := make([]User, 0)
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
		if err != nil && !strings.Contains(err.Error(), "unsupported Scan, storing driver.Value type <nil> into type *time.Time") {
			log.Printf("Error while loading data from database: %v\n", errors.New(err).ErrorStack())
			response.WriteErrorString(http.StatusConflict, fmt.Sprintf("Invalid data in database, check logs."))
			return
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("User with id %v could not be found.", id))
	} else if len(users) > 1 {
		response.WriteErrorString(http.StatusConflict, fmt.Sprintf("Found multiple users with the same id %v, exiting...", id))
	} else {
		response.WriteEntity(users[0])
	}
}

func (db *Conn) CreateUser(request *restful.Request, response *restful.Response) {
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

	ctx := InitContext()
	_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO user_profile VALUES(null, '%v', utl_raw.cast_to_raw('%v'), default, default)", userRequest.Email, string(hashedPassword)))
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Failed to create user %v, %v\nerr: %v", userRequest.Email, userRequest.Password, err))
		return
	}
}
