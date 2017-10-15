package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	_ "github.com/mattn/go-oci8"

	"net/http"

	"io/ioutil"

	"encoding/json"

	"strconv"

	"path/filepath"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/ahouts/ProDuctive-server/migrations"
	"github.com/ahouts/ProDuctive-server/tunnel"
	"github.com/emicklei/go-restful"
	"github.com/miquella/ask"
	"golang.org/x/crypto/ssh"
	"gopkg.in/rana/ora.v4"
)

type loginInfo struct {
	Hostname string
	Username string
	Password string
}

type configuration struct {
	Ssh    loginInfo
	Db     loginInfo
	DbName string
}

const localPort = 1521
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
	fullPath, err := filepath.Abs(cfgFile)
	if err != nil {
		log.Fatalf("Failed to resolve path %v\n%v", cfgFile, err)
	}

	file, e := ioutil.ReadFile(fullPath)
	if e != nil {
		log.Fatalf("File error: %v\n", e)
	}
	cfg := configuration{}
	json.Unmarshal(file, &cfg)

	createTunnel(cfg)

	dbConn := cfg.Db.Username + `/` + cfg.Db.Password + `@localhost:` + strconv.Itoa(localPort) + "/" + cfg.DbName
	db, err := sql.Open("oci8", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database...\n%v", err)
	}
	defer db.Close()

	// Set timeout
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// Set prefetch count
	ctx = ora.WithStmtCfg(ctx, ora.Cfg().StmtCfg.SetPrefetchRowCount(dbPrefetchRowCount))

	c := &data.Conn{DB: *db, Ctx: ctx}
	mConn := migrations.MCon(*c)
	mConn.RunMigrations()

	ws := new(restful.WebService)
	configureRoutes(c, ws)

	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}

func createTunnel(cfg configuration) {
	localEndpoint := &tunnel.Endpoint{
		Host: "localhost",
		Port: localPort,
	}

	serverEndpoint := &tunnel.Endpoint{
		Host: cfg.Ssh.Hostname,
		Port: sshPort,
	}

	remoteEndpoint := &tunnel.Endpoint{
		Host: cfg.Db.Hostname,
		Port: oraclePort,
	}

	tun := &tunnel.SSHTunnel{
		Config: &ssh.ClientConfig{
			User: cfg.Ssh.Username,
			Auth: []ssh.AuthMethod{
				ssh.Password(cfg.Ssh.Password),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
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
