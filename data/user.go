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
	query := fmt.Sprintf("SELECT id, email, password_hash, created_at, updated_at FROM user_profile WHERE id=%v", id)
	fmt.Println(query)
	ctx := InitContext()
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalln(errors.New(err).ErrorStack())
	}
	defer rows.Close()
	defer ctx.Done()
	users := make([]User, 0)
	for rows.Next() {
		fmt.Println("got hrere")
		var u User
		err = rows.Scan(&u.id, &u.email, &u.password_hash, &u.created_at, &u.updated_at)
		if err != nil {
			log.Fatalf("Error while loading data from database: %v\n", errors.New(err).ErrorStack())
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
