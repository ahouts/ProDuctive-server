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

type Project struct {
	Id        int
	Title     string
	OwnerId   int
	UserIds   []int
	NoteIds   []int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProjectMetadata struct {
	Id    int
	Title string
}

type GetProjectRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func (s *DbSession) GetProjects(request *restful.Request, response *restful.Response) {
	projectRequest := GetProjectRequest{}
	err := request.ReadEntity(&projectRequest)

	if err != nil {
		formatError(new(GetProjectRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, projectRequest.Email, projectRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	projectMetadata := make([]ProjectMetadata, 0)

	rows, err := tx.Query("SELECT id, title FROM project WHERE id in (select * from table(get_projects_for_user(:1)))", userId)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to find projects for user: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	for rows.Next() {
		metadata := ProjectMetadata{}
		err = rows.Scan(&metadata.Id, &metadata.Title)
		if err != nil {
			response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("error while loading row from database, check logs"))
			log.Printf("error while loading row from database: %v\n", errors.New(err).ErrorStack())
			tx.Rollback()
			return
		}
		projectMetadata = append(projectMetadata, metadata)
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

	response.WriteEntity(projectMetadata)
}

type CreateProjectRequest struct {
	Email    string
	Password string
	Title    string
}

func (s *DbSession) CreateProject(request *restful.Request, response *restful.Response) {
	projectRequest := CreateProjectRequest{}
	err := request.ReadEntity(&projectRequest)
	if err != nil {
		formatError(new(CreateProjectRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, projectRequest.Email, projectRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO project VALUES(null, :1, :2, default, default)", projectRequest.Title, userId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to create project: %v", err))
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

func (s *DbSession) GetProject(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("project-id")
	projectId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, project id %v is invalid.\n%v", idStr, err))
		return
	}

	projectRequest := GetProjectRequest{}
	err = request.ReadEntity(&projectRequest)
	if err != nil {
		formatError(new(GetProjectRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, projectRequest.Email, projectRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	var hasPermission int
	err = tx.QueryRow("SELECT * FROM TABLE(permission_for_project(:1, :2))", userId, projectId).Scan(&hasPermission)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		log.Println(errors.New(err).ErrorStack())
		return
	}

	if hasPermission == 0 {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user does not have permission to view project %v", projectId))
		tx.Rollback()
		return
	}

	res := Project{}
	res.NoteIds = make([]int, 0)
	res.UserIds = make([]int, 0)
	err = tx.QueryRow("SELECT id, title, owner_id, created_at, updated_at FROM project WHERE id = :1", projectId).
		Scan(&res.Id, &res.Title, &res.OwnerId, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to query db for project %v: %v", projectId, err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	rows, err := tx.Query("select * from table(get_users_for_project(:1))", projectId)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to find users for project: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	for rows.Next() {
		var uid int
		err = rows.Scan(&uid)
		if err != nil {
			response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("error while loading row from database, check logs"))
			log.Printf("error while loading row from database: %v\n", errors.New(err).ErrorStack())
			tx.Rollback()
			return
		}
		res.UserIds = append(res.UserIds, uid)
	}
	if err = rows.Err(); err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("error while loading row from database, check logs"))
		log.Printf("error while loading row from database: %v\n", errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}
	rows.Close()

	rows, err = tx.Query("select * from table(get_notes_for_project(:1))", projectId)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to find notes for project: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}
	for rows.Next() {
		var nid int
		err = rows.Scan(&nid)
		if err != nil {
			response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("error while loading row from database, check logs"))
			log.Printf("error while loading row from database: %v\n", errors.New(err).ErrorStack())
			tx.Rollback()
			return
		}
		res.NoteIds = append(res.NoteIds, nid)
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
		log.Println(errors.New(err).ErrorStack())
		return
	}
	response.WriteEntity(res)
}

type AddUserToProjectRequest struct {
	Email     string
	Password  string
	NewUserId int
}

func (s *DbSession) AddUserToProject(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("project-id")
	projectId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, project id %v is invalid.\n%v", idStr, err))
		return
	}

	projectRequest := AddUserToProjectRequest{}
	err = request.ReadEntity(&projectRequest)
	if err != nil {
		formatError(new(AddUserToProjectRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context.\n%v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, projectRequest.Email, projectRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	var hasPermission int
	err = tx.QueryRow("SELECT * FROM TABLE(permission_for_project(:1, :2))", userId, projectId).Scan(&hasPermission)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		log.Println(errors.New(err).ErrorStack())
		return
	}

	if hasPermission == 0 {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user does not have permission to add user to project %v", projectId))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("INSERT INTO project_user VALUES(:1, :2, 1)", projectRequest.NewUserId, projectId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to add user to project: %v", err))
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

type DeleteProjectRequest struct {
	Email    string
	Password string
}

func (s *DbSession) DeleteProject(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("project-id")
	projectId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, project id %v is invalid.\n%v", idStr, err))
		return
	}

	projectRequest := DeleteProjectRequest{}
	err = request.ReadEntity(&projectRequest)
	if err != nil {
		formatError(new(DeleteProjectRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, projectRequest.Email, projectRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		tx.Rollback()
		return
	}

	var hasPermission int
	err = tx.QueryRow("SELECT * FROM TABLE(permission_for_project(:1, :2))", userId, projectId).Scan(&hasPermission)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		log.Println(errors.New(err).ErrorStack())
		return
	}

	if hasPermission == 0 {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user does not have permission to view project %v", projectId))
		tx.Rollback()
		return
	}

	_, err = tx.Exec("DELETE FROM project WHERE id = :1", projectId)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to delete project %v: %v", projectId, err))
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

func (s *DbSession) GetNotesForProject(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("project-id")
	projectId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("Invalid query, project id %v is invalid.\n%v", idStr, err))
		return
	}

	projectRequest := GetProjectRequest{}
	err = request.ReadEntity(&projectRequest)

	if err != nil {
		formatError(new(GetProjectRequest), response)
		return
	}

	tx, err := s.InitTransaction()
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to initialize database context: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		return
	}

	userId, err := AuthUser(tx, projectRequest.Email, projectRequest.Password)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("failed to authenticate request: %v.", err))
		log.Println(errors.New(err).ErrorStack())
		tx.Rollback()
		return
	}

	var hasPermission int
	err = tx.QueryRow("SELECT * FROM TABLE(permission_for_project(:1, :2))", userId, projectId).Scan(&hasPermission)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		log.Println(errors.New(err).ErrorStack())
		return
	}

	if hasPermission == 0 {
		response.WriteErrorString(http.StatusBadRequest, fmt.Sprintf("user does not have permission to view project %v", projectId))
		tx.Rollback()
		return
	}

	noteMetadata := make([]NoteMetadata, 0)

	rows, err := tx.Query("SELECT id, title, owner_id, project_id FROM note WHERE id in (select * from table(get_notes_for_project(:1)))", projectId)
	if err != nil {
		tx.Rollback()
		response.WriteErrorString(http.StatusInternalServerError, fmt.Sprintf("failed to find notes for project: %v.", err))
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
