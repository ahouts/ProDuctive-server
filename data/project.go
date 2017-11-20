package data

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
)

type Project struct {
	Id        int
	title     string
	OwnerId   int
	UserIds   []int
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
