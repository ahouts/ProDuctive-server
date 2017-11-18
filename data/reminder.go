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

	fmt.Printf("User id: %v\n", id)

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
		fmt.Printf("Scanned reminder %v\n", r)
		reminders = append(reminders, r)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error while loading row from database: %v\n", errors.New(err).ErrorStack())
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Error while loading row from database, check logs."))
		tx.Rollback()
		return
	}

	rows.Close()
	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("Failed to commit change.\n%v", err))
		tx.Rollback()
		return
	}

	fmt.Printf("All reminders %v\n", reminders)

	response.WriteEntity(reminders)
}

type CreateReminderRequest struct {
	Email    string
	Password string
	Body     string
}

func (s *DbSession) CreateReminder(request *restful.Request, response *restful.Response) {
	reminderRequest := CreateReminderRequest{}
	err := request.ReadEntity(&reminderRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", reminderRequest))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	id, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO reminder VALUES(null, :1, :2, default, default)", id, reminderRequest.Body)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed insert reminder: %v", err))
		log.Println(errors.New(err).ErrorStack())
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

type UpdateReminderRequest struct {
	Email      string
	Password   string
	ReminderId int
	Body       string
}

func (s *DbSession) UpdateReminder(request *restful.Request, response *restful.Response) {
	reminderRequest := UpdateReminderRequest{}
	err := request.ReadEntity(&reminderRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", reminderRequest))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	id, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	var reminderUserId int
	err = tx.QueryRow("SELECT user_id FROM reminder WHERE id=:1", reminderRequest.ReminderId).Scan(&reminderUserId)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to query db for reminder id: %v", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	if reminderUserId != id {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("reminder %v is not user %v's id", reminderRequest.ReminderId, id))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE reminder SET body = :1, updated_at = :2 WHERE id = :3", reminderRequest.Body, time.Now(), reminderRequest.ReminderId)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to update reminder: %v", err))
		log.Println(errors.New(err).ErrorStack())
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
