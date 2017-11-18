package main

import (
	"github.com/ahouts/ProDuctive-server/data"
	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
)

func setupRoutes(s *data.DbSession) {
	restful.Add(userWs(s))
	restful.Add(getUserIdWs(s))
	restful.Add(remindersWs(s))
	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "./swagger-dist"}
	swagger.InstallSwaggerService(config)
}

func userWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/user").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/").To(s.CreateUser).
		Doc("create a user").
		Reads(new(data.CreateUserRequest)))

	ws.Route(ws.GET("/{user-id}").To(s.GetUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(new(data.User)))

	ws.Filter(enableCORS)
	return ws
}

func getUserIdWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/get_user_id").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/").To(s.GetUserId).
		Doc("get a user's id").
		Reads(new(data.GetUserIdRequest)).
		Writes(new(data.UserId)))

	ws.Filter(enableCORS)
	return ws
}

func remindersWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/reminder").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PUT("/get").To(s.GetReminders).
		Doc("get a user's reminders").
		Reads(data.GetReminderRequest{}).
		Writes(new([]data.Reminder)))

	ws.Route(ws.PUT("/get/{reminder-id}").To(s.GetReminder).
		Doc("get a user's reminder").
		Param(ws.PathParameter("reminder-id", "id of the reminder").DataType("string")).
		Reads(data.GetReminderRequest{}).
		Writes(new(data.Reminder)))

	ws.Route(ws.POST("/").To(s.CreateReminder).
		Doc("create a reminder").
		Reads(new(data.CreateReminderRequest)))

	ws.Route(ws.PUT("/").To(s.UpdateReminder).
		Doc("update a reminder").
		Reads(new(data.UpdateReminderRequest)))

	ws.Filter(enableCORS)
	return ws
}

func enableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	resp.AddHeader("Access-Control-Allow-Origin", "*")
	resp.AddHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	chain.ProcessFilter(req, resp)
}
