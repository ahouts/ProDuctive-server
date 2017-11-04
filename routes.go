package main

import (
	"github.com/ahouts/ProDuctive-server/data"
	restful "github.com/emicklei/go-restful"
)

func setupRoutes(c *data.Conn) {
	restful.Add(userWs(c))
	restful.Add(createUserWs(c))
}

func userWs(c *data.Conn) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{user-id}").To(c.GetUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(data.User{}))

	return ws
}

func createUserWs(c *data.Conn) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/create_user").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/").To(c.CreateUser).
		Doc("create a user").
		Reads(data.CreateUserRequest{}))

	return ws
}
