package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
)

type User struct {
	id            int
	email         string
	password_hash string
	created_at    time.Time
	updated_at    time.Time
}

func (db *Conn) GetUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	rows, err := db.QueryContext(db.Ctx, "SELECT id, email, password_hash, created_at, updated_at FROM user_profile WHERE id=?", id)
	if err != nil {
		log.Fatalln(err.(*errors.Error).ErrorStack())
	}
	defer rows.Close()
	users := make([]User, 0)
	for rows.Next() {
		var u User
		err = rows.Scan(&u.id, &u.email, &u.password_hash, &u.created_at, &u.updated_at)
		if err != nil {
			log.Fatalf("Error while loading data from database: %v\n", err)
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		response.WriteErrorString(http.StatusNoContent, fmt.Sprintf("User with id %v could not be found.", id))
	} else if len(users) > 1 {
		response.WriteErrorString(http.StatusConflict, fmt.Sprintf("Found multiple users with the same id %v, exiting...", id))
	} else {
		response.WriteEntity(users[0])
	}
}
