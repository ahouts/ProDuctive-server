package main

import "github.com/emicklei/go-restful"
import "github.com/ahouts/ProDuctive-server/db"


func main() {
	ws := new(restful.WebService)
	ws.
	Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(db.User{}))
	...

	func (u db.User) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	...
	}
}
