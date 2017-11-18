package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
)

type Reminder struct {
	Id        int
	UserId    int
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetRemindersRequest struct {
	Email    string
	Password string
}

func (s *DbSession) GetReminders(request *restful.Request, response *restful.Response) {
	reminderRequest := GetRemindersRequest{}
	err := request.ReadEntity(&reminderRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Request invalid, must match format:\n%v", reminderRequest))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	id, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Failed to authenticate request.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	rows, err := tx.Query("SELECT id, user_id, body, created_at, updated_at FROM reminder WHERE user_id = :1", id)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Failed to query db for reminders.\n%v", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	reminders := make([]Reminder, 0)
	for rows.Next() {
		var r Reminder
		err = rows.Scan(&r.Id, &r.UserId, &r.Body, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			log.Printf("Error while loading row from database: %v\n", errors.New(err).ErrorStack())
			response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Error while loading row from database, check logs."))
			tx.Rollback()
			return
		}
		reminders = append(reminders, r)
	}

	rows.Close()
	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Failed to commit change.\n%v", err))
		tx.Rollback()
		return
	}

	response.WriteEntity(reminders)
}
