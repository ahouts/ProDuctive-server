package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func formatError(format interface{}, response *restful.Response) {
	b, err := json.MarshalIndent(format, "", "  ")
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to find format for route.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", string(b)))
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
