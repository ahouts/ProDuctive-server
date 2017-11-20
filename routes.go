package main

import (
	"github.com/ahouts/ProDuctive-server/data"
	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
)

func setupRoutes(s *data.DbSession) {
	restful.Add(userWs(s))
	restful.Add(remindersWs(s))
	restful.Add(noteWs(s))
	restful.Add(projectWs(s))
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

	ws.Route(ws.PUT("/getid").To(s.GetUserId).
		Doc("get a user's id").
		Reads(new(data.GetUserIdRequest)).
		Writes(new(data.UserId)))

	ws.Route(ws.PUT("/").To(s.GetUser).
		Doc("get a user's info").
		Reads(new(data.GetUserRequest)).
		Writes(new(data.User)))

	ws.Route(ws.POST("/").To(s.CreateUser).
		Doc("create a user").
		Reads(new(data.CreateUserRequest)))

	ws.Filter(enableCORS)
	return ws
}

func remindersWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/reminder").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PUT("/").To(s.GetReminders).
		Doc("get a user's reminders").
		Reads(data.GetReminderRequest{}).
		Writes(new([]data.Reminder)))

	ws.Route(ws.PUT("/{reminder-id}").To(s.GetReminder).
		Doc("get a user's reminder").
		Param(ws.PathParameter("reminder-id", "id of the reminder").DataType("string")).
		Reads(data.GetReminderRequest{}).
		Writes(new(data.Reminder)))

	ws.Route(ws.POST("/").To(s.CreateReminder).
		Doc("create a reminder").
		Reads(new(data.CreateReminderRequest)))

	ws.Route(ws.POST("/{reminder-id}").To(s.UpdateReminder).
		Doc("update a reminder").
		Reads(new(data.UpdateReminderRequest)))

	ws.Route(ws.PUT("/delete/{reminder-id}").To(s.DeleteReminder).
		Doc("delete a reminder").
		Reads(new(data.DeleteReminderRequest)))

	ws.Filter(enableCORS)
	return ws
}

func noteWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/note").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PUT("/").To(s.GetNotes).
		Doc("get a user's notes").
		Reads(data.GetNoteRequest{}).
		Writes(new([]data.NoteMetadata)))

	ws.Route(ws.PUT("/{note-id}").To(s.GetNote).
		Doc("get a user's note").
		Param(ws.PathParameter("note-id", "id of the note").DataType("string")).
		Reads(data.GetNoteRequest{}).
		Writes(new(data.Note)))

	ws.Route(ws.POST("/{note-id}/add_user").To(s.AddUserToNote).
		Doc("add a user to a note").
		Param(ws.PathParameter("note-id", "id of the note").DataType("string")).
		Reads(data.AddUserToNoteRequest{}))

	ws.Route(ws.POST("/").To(s.CreateNote).
		Doc("create a note").
		Reads(new(data.CreateNoteRequest)))

	ws.Route(ws.POST("/{note-id}").To(s.UpdateNote).
		Doc("update a note").
		Param(ws.PathParameter("note-id", "id of the note").DataType("string")).
		Reads(new(data.UpdateNoteRequest)))

	ws.Route(ws.PUT("/delete/{note-id}").To(s.DeleteNote).
		Doc("delete a note").
		Param(ws.PathParameter("note-id", "id of the note").DataType("string")).
		Reads(new(data.DeleteNoteRequest)))

	ws.Filter(enableCORS)
	return ws
}

func projectWs(s *data.DbSession) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/project").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PUT("/").To(s.GetProjects).
		Doc("get a user's projects").
		Reads(data.GetProjectRequest{}).
		Writes(new([]data.ProjectMetadata)))
	//
	//ws.Route(ws.PUT("/{project-id}").To(s.GetProject).
	//	Doc("get a user's project").
	//	Param(ws.PathParameter("project-id", "id of the project").DataType("string")).
	//	Reads(data.GetProjectRequest{}).
	//	Writes(new(data.Project)))
	//
	//ws.Route(ws.PUT("/{project-id}/get_notes").To(s.GetNotesForProject).
	//	Doc("get notes for a project").
	//	Param(ws.PathParameter("project-id", "id of the project").DataType("string")).
	//	Reads(data.GetNotesForProjectRequest{}).
	//	Writes(new([]data.NoteMetadata)))
	//
	//ws.Route(ws.POST("/{project-id}/add_user").To(s.AddUserToProject).
	//	Doc("add a user to a project").
	//	Param(ws.PathParameter("project-id", "id of the project").DataType("string")).
	//	Reads(data.AddUserToProjectRequest{}))
	//
	ws.Route(ws.POST("/").To(s.CreateProject).
		Doc("create a project").
		Reads(new(data.CreateProjectRequest)))
	//
	//ws.Route(ws.POST("/{project-id}").To(s.UpdateNote).
	//	Doc("update a project").
	//	Param(ws.PathParameter("project-id", "id of the project").DataType("string")).
	//	Reads(new(data.UpdateProjectRequest)))
	//
	//ws.Route(ws.PUT("/delete/{project-id}").To(s.DeleteProject).
	//	Doc("delete a project").
	//	Param(ws.PathParameter("project-id", "id of the project").DataType("string")).
	//	Reads(new(data.DeleteProjectRequest)))

	ws.Filter(enableCORS)
	return ws
}

func enableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	resp.AddHeader("Access-Control-Allow-Origin", "*")
	resp.AddHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	chain.ProcessFilter(req, resp)
}
