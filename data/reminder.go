package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"strconv"

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

type GetReminderRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func (s *DbSession) GetReminders(request *restful.Request, response *restful.Response) {
	reminderRequest := GetReminderRequest{}
	err := request.ReadEntity(&reminderRequest)

	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", reminderRequest))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	id, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	rows, err := tx.Query("SELECT id, user_id, body, created_at, updated_at FROM reminder WHERE user_id = :1", id)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to query db for reminders: %v", err))
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
			response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("error while loading row from database, check logs"))
			tx.Rollback()
			return
		}
		reminders = append(reminders, r)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error while loading row from database: %v\n", errors.New(err).ErrorStack())
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("error while loading row from database, check logs"))
		tx.Rollback()
		return
	}

	rows.Close()
	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		log.Println(errors.New(err).ErrorStack())
		return
	}

	response.WriteEntity(reminders)
}

func (s *DbSession) GetReminder(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("reminder-id")
	reminderId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, reminder id %v is invalid.\n%v", idStr, err))
		return
	}

	reminderRequest := GetReminderRequest{}
	err = request.ReadEntity(&reminderRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", reminderRequest))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	res := Reminder{}
	err = tx.QueryRow("SELECT id, user_id, body, created_at, updated_at FROM reminder WHERE user_id = :1 AND id = :2", userId, reminderId).
		Scan(&res.Id, &res.UserId, &res.Body, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to query db for reminder %v: %v", reminderId, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		log.Println(errors.New(err).ErrorStack())
		return
	}
	response.WriteEntity(res)
}

type CreateReminderRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Body     string `json:"Body"`
}

func (s *DbSession) CreateReminder(request *restful.Request, response *restful.Response) {
	reminderRequest := CreateReminderRequest{}
	err := request.ReadEntity(&reminderRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", reminderRequest))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO reminder VALUES(null, :1, :2, default, default)", userId, reminderRequest.Body)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to create reminder: %v", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		log.Println(errors.New(err).ErrorStack())
		return
	}
}

type UpdateReminderRequest struct {
	Email    string
	Password string
	Body     string
}

func (s *DbSession) UpdateReminder(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("reminder-id")
	reminderId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, reminder id %v is invalid.\n%v", idStr, err))
		return
	}

	reminderRequest := UpdateReminderRequest{}
	err = request.ReadEntity(&reminderRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", reminderRequest))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	var reminderUserId int
	err = tx.QueryRow("SELECT user_id FROM reminder WHERE id=:1", reminderId).Scan(&reminderUserId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to query db for reminder id %v: %v", reminderId, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	if reminderUserId != userId {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("reminder %v is not user %v's reminder", reminderId, userId))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE reminder SET body = :1, updated_at = :2 WHERE id = :3", reminderRequest.Body, time.Now(), reminderId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to update reminder: %v", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		log.Println(errors.New(err).ErrorStack())
		return
	}
}

type DeleteReminderRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func (s *DbSession) DeleteReminder(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("reminder-id")
	reminderId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, reminder id %v is invalid.\n%v", idStr, err))
		return
	}

	reminderRequest := GetReminderRequest{}
	err = request.ReadEntity(&reminderRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", reminderRequest))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, reminderRequest.Email, reminderRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	res := Reminder{}
	err = tx.QueryRow("SELECT id, user_id, body, created_at, updated_at FROM reminder WHERE user_id = :1 AND id = :2", userId, reminderId).
		Scan(&res.Id, &res.UserId, &res.Body, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to query db for reminder %v: %v", reminderId, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	if res.UserId != userId {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user %v doesn't have permission to delete reminder %v: %v", userId, reminderId, err))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("DELETE FROM reminder WHERE id = :1", reminderId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to delete reminder %v: %v", reminderId, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		log.Println(errors.New(err).ErrorStack())
		return
	}
}
