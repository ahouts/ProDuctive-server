package main

import (
	"github.com/ahouts/ProDuctive-server/data"
	restful "github.com/emicklei/go-restful"
)

func configureRoutes(c *data.Conn, ws *restful.WebService) {
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/users/{user-id}").To(c.GetUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "id of the user").DataType("int")).
		Writes(data.User{}))
}
