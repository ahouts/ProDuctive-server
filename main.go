package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "gopkg.in/rana/ora.v4"

	"net/http"

	"fmt"
	"io/ioutil"
	"os"

	"encoding/json"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/ahouts/ProDuctive-server/tunnel"
	"github.com/emicklei/go-restful"
	"github.com/miquella/ask"
	"golang.org/x/crypto/ssh"
	"gopkg.in/rana/ora.v4"
)

type loginInfo struct {
	hostname string
	username string
	password string
}

type configuration struct {
	ssh    loginInfo
	db     loginInfo
	dbName string
}

const localPort = 48620
const oraclePort = 1521
const sshPort = 22
const dbPrefetchRowCount = 50000

func main() {
	cfgFile, err := ask.Ask("config file? (./config.json): ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	if cfgFile == "" {
		cfgFile = "./config.json"
	}
	file, e := ioutil.ReadFile(cfgFile)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	cfg := configuration{}
	json.Unmarshal(file, &cfg)

	createTunnel(cfg)

	dbConn := cfg.db.username + "/" + cfg.db.password + "@\"localhost:" + string(localPort) + "/" + cfg.dbName + "\""
	db, err := sql.Open("ora", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database...\n %v", err)
	}
	defer db.Close()

	// Set timeout
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// Set prefetch count
	ctx = ora.WithStmtCfg(ctx, ora.Cfg().StmtCfg.SetPrefetchRowCount(dbPrefetchRowCount))
	rows, err := db.QueryContext(ctx, "SELECT * FROM user_objects")
	defer rows.Close()

	// do migrations

	ws := new(restful.WebService)
	c := &data.Conn{DB: *db, Ctx: ctx}
	configureRoutes(c, ws)

	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}

func createTunnel(cfg configuration) {
	localEndpoint := &tunnel.Endpoint{
		Host: "localhost",
		Port: localPort,
	}

	serverEndpoint := &tunnel.Endpoint{
		Host: cfg.ssh.hostname,
		Port: sshPort,
	}

	remoteEndpoint := &tunnel.Endpoint{
		Host: cfg.db.hostname,
		Port: oraclePort,
	}

	tun := &tunnel.SSHTunnel{
		Config: &ssh.ClientConfig{
			User: cfg.ssh.username,
			Auth: []ssh.AuthMethod{
				ssh.Password(cfg.ssh.password),
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
