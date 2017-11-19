package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
)

type Note struct {
	Id        int
	Title     string
	Body      string
	OwnerId   int
	ProjectId sql.NullInt64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetNoteRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type NoteMetadata struct {
	Id        int
	Title     string
	OwnerId   int
	ProjectId sql.NullInt64
}

func (s *DbSession) GetNotes(request *restful.Request, response *restful.Response) {
	noteRequest := GetNoteRequest{}
	err := request.ReadEntity(&noteRequest)

	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", noteRequest))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, noteRequest.Email, noteRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	noteMetadata := make([]NoteMetadata, 0)

	rows, err := tx.Query("SELECT id, title, owner_id, project_id FROM note WHERE id in (select * from table(get_notes_for_user(:1)))", userId)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to find notes for user: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	for rows.Next() {
		metadata := NoteMetadata{}
		err = rows.Scan(&metadata.Id, &metadata.Title, &metadata.OwnerId, &metadata.ProjectId)
		if err != nil {
			response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("error while loading row from database, check logs"))
			log.Printf("error while loading row from database: %v\n", errors.New(err).ErrorStack())
			tx.Rollback()
			return
		}
		noteMetadata = append(noteMetadata, metadata)
	}
	if err = rows.Err(); err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("error while loading row from database, check logs"))
		log.Printf("error while loading row from database: %v\n", errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
	rows.Close()

	err = tx.Commit()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to commit change: %v", err))
		tx.Rollback()
		return
	}

	response.WriteEntity(noteMetadata)
}

func (s *DbSession) GetNote(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("note-id")
	noteId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, note id %v is invalid.\n%v", idStr, err))
		return
	}

	noteRequest := GetNoteRequest{}
	err = request.ReadEntity(&noteRequest)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("request invalid, must match format: %v", noteRequest))
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, noteRequest.Email, noteRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	var hasPermission int
	err = tx.QueryRow("SELECT * FROM TABLE(user_has_permission_for_note(:1, :2))", userId, noteId).Scan(&hasPermission)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		log.Println(errors.New(err).ErrorStack())
		return
	}

	if hasPermission == 0 {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user does not have permission to view note %v", noteId))
		tx.Rollback()
		return
	}

	res := Note{}
	err = tx.QueryRow("SELECT id, title, body, owner_id, project_id, created_at, updated_at FROM note WHERE id = :1", noteId).
		Scan(&res.Id, &res.Title, &res.Body, &res.OwnerId, &res.ProjectId, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to query db for note %v: %v", noteId, err))
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

type CreateNoteRequest struct {
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	Title     string `json:"Title"`
	Body      string `json:"Body"`
	ProjectId int    `json:"ProjectId"`
}

func (s *DbSession) CreateNote(request *restful.Request, response *restful.Response) {
	reminderRequest := CreateNoteRequest{}
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

	if reminderRequest.ProjectId == 0 {
		_, err = tx.Exec("INSERT INTO note VALUES(null, :1, :2, :3, :4, default, default)",
			reminderRequest.Title, reminderRequest.Body, userId, nil)
	} else {
		_, err = tx.Exec("INSERT INTO note VALUES(null, :1, :2, :3, :4, default, default)",
			reminderRequest.Title, reminderRequest.Body, userId, reminderRequest.ProjectId)
	}

	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
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
