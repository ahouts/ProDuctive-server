package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "gopkg.in/rana/ora.v4"

	"net/http"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/ahouts/ProDuctive-server/tunnel"
	"github.com/emicklei/go-restful"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	"github.com/miquella/ask"
	"golang.org/x/crypto/ssh"
	"gopkg.in/rana/ora.v4"
)

func main() {
	sshHostname, err := ask.Ask("SSH Hostname: ")
	if err != nil {
		log.Fatalf(string(err))
	}
	sshUsername, err := ask.Ask("SSH Username: ")
	if err != nil {
		log.Fatalf(string(err))
	}
	sshPasswd, err := ask.HiddenAsk("SSH Password: ")
	if err != nil {
		log.Fatalf(string(err))
	}
	dbHostname, err := ask.Ask("Database Hostname: ")
	if err != nil {
		log.Fatalf(string(err))
	}
	dbUsername, err := ask.Ask("Database Username: ")
	if err != nil {
		log.Fatalf(string(err))
	}
	dbPasswd, err := ask.HiddenAsk("Database Password: ")
	if err != nil {
		log.Fatalf(string(err))
	}

	createTunnel(sshHostname, sshUsername, sshPasswd, dbHostname)

	dbConn := dbUsername + "/" + dbPasswd + "@localhost:48620"
	db, err := sql.Open("ora", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database...\n %v", err)
	}
	defer db.Close()

	// Set timeout
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// Set prefetch count
	ctx = ora.WithStmtCfg(ctx, ora.Cfg().StmtCfg.SetPrefetchRowCount(50000))
	rows, err := db.QueryContext(ctx, "SELECT * FROM user_objects")
	defer rows.Close()

	dbDriver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "migration_table",
		DatabaseName:    "",
	})
	if err != nil {

	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"oracle",
		dbDriver,
	)

	m.Up()

	ws := new(restful.WebService)
	c := &data.Conn{DB: *db, Ctx: ctx}
	configureRoutes(c, ws)

	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}

func createTunnel(sshHostname, username, password, dbHostname string) {
	localEndpoint := &tunnel.Endpoint{
		Host: "localhost",
		Port: 48620,
	}

	serverEndpoint := &tunnel.Endpoint{
		Host: sshHostname,
		Port: 22,
	}

	remoteEndpoint := &tunnel.Endpoint{
		Host: dbHostname,
		Port: 1521,
	}

	tun := &tunnel.SSHTunnel{
		Config: &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
		},
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	go tun.Start()
}

func configureRoutes(c *data.Conn, ws *restful.WebService) {
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{user-id}").To(c.GetUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "id of the user").DataType("int")).
		Writes(data.User{}))
}
