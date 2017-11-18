package main

import (
	"github.com/ahouts/ProDuctive-server/data"
	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
)

func setupRoutes(s *data.DbSession) {
	restful.Add(userWs(s))
	restful.Add(createUserWs(s))
	restful.Add(getRemindersWs(s))
	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "./swagger-dist"}
	swagger.InstallSwaggerService(config)
}

func userWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{user-id}").To(s.GetUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(data.User{}))

	return ws
}

func createUserWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/create_user").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/").To(s.CreateUser).
		Doc("create a user").
		Reads(data.CreateUserRequest{}))

	return ws
}

func getRemindersWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/reminders").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/").To(s.GetReminders).
		Doc("get a user's reminders").
		Reads(data.GetRemindersRequest{}).
		Writes([]data.Reminder{}))
	return ws
}
