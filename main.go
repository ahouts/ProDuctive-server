package main

import (
	"log"

	"github.com/ahouts/ProDuctive-server/db"
	"github.com/emicklei/go-restful"
)

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

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServeTls(":443", "cert.pem", "key.pem", nil))
}

func (u db.User) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
}

/*
example of connecting to oracle db
import (
	"database/sql"

	_ "gopkg.in/rana/ora.v4"
)

func main() {
	dbConn := "system/" + os.Getenv("ORACLE_PWD") + "@oracledb:1521"
	db, err := sql.Open("ora", dbConn)
	defer db.Close()

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Set prefetch count
	ctx = ora.WithStmtCfg(ctx, ora.Cfg().StmtCfg.SetPrefetchCount(50000))
	rows, err := db.QueryContext(ctx, "SELECT * FROM user_objects")
	defer rows.Close()
}
*/
