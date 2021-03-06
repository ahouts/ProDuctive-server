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
		formatError(new(GetNoteRequest), response)
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
		formatError(new(GetNoteRequest), response)
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
	Email     string        `json:"Email"`
	Password  string        `json:"Password"`
	Title     string        `json:"Title"`
	Body      string        `json:"Body"`
	ProjectId sql.NullInt64 `json:"ProjectId"`
}

func (s *DbSession) CreateNote(request *restful.Request, response *restful.Response) {
	noteRequest := CreateNoteRequest{}
	err := request.ReadEntity(&noteRequest)
	if err != nil {
		formatError(new(CreateNoteRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, noteRequest.Email, noteRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO note VALUES(null, :1, :2, :3, :4, default, default)",
		noteRequest.Title, noteRequest.Body, userId, noteRequest.ProjectId)
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

type DeleteNoteRequest struct {
	Email    string
	Password string
}

func (s *DbSession) DeleteNote(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("note-id")
	noteId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, note id %v is invalid.\n%v", idStr, err))
		return
	}

	noteRequest := DeleteNoteRequest{}
	err = request.ReadEntity(&noteRequest)
	if err != nil {
		formatError(new(DeleteNoteRequest), response)
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

	_, err = tx.Exec("DELETE FROM note WHERE id = :1", noteId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to delete note %v: %v", noteId, err))
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

type UpdateNoteRequest struct {
	Email     string
	Password  string
	Title     string
	Body      string
	OwnerId   int
	ProjectId sql.NullInt64
}

func (s *DbSession) UpdateNote(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("note-id")
	noteId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, note id %v is invalid.\n%v", idStr, err))
		return
	}

	noteRequest := UpdateNoteRequest{}
	err = request.ReadEntity(&noteRequest)
	if err != nil {
		formatError(new(UpdateNoteRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
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
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user does not have permission to update note %v", noteId))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE note SET title = :1, body = :2, owner_id = :3, project_id = :4 WHERE id = :5", noteRequest.Title, noteRequest.Body, noteRequest.OwnerId, noteRequest.ProjectId, noteId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to delete note: %v", err))
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

type AddUserToNoteRequest struct {
	Email     string
	Password  string
	NewUserId int
}

func (s *DbSession) AddUserToNote(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("note-id")
	noteId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, note id %v is invalid.\n%v", idStr, err))
		return
	}

	noteRequest := AddUserToNoteRequest{}
	err = request.ReadEntity(&noteRequest)
	if err != nil {
		formatError(new(AddUserToNoteRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
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
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user does not have permission to add user to note %v", noteId))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO note_user VALUES(:1, :2, 1)", noteRequest.NewUserId, noteId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to add user to note: %v", err))
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
