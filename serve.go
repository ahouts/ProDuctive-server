package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-oci8"

	"fmt"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/ahouts/ProDuctive-server/migrations"
	"github.com/ahouts/ProDuctive-server/tunnel"
	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/ssh"
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

const localPort = 40841
const oraclePort = 1521
const sshPort = 22

func serve(cfgFile string, port int) {
	c, mConn := initDb(cfgFile)
	mConn.Up()

	restful.Add(userWs(c))

	fmt.Printf("Environment successfully configured, beggining to serve on port %v.\n", strconv.Itoa(port))

	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(port), "cert.pem", "key.pem", nil))
}

func dropDb(cfgFile string) {
	_, mConn := initDb(cfgFile)
	mConn.Down()
}

func initDb(cfgFile string) (*data.Conn, migrations.MCon) {
	cfg := getCfg(cfgFile)

	// this function returns a channel that returns a value once its done
	// by pulling a value out of the channel, we block until the tunnel is ready
	<-cfg.createTunnel()

	db := cfgDb(cfg)
	defer db.Close()

	c := &data.Conn{DB: *db}
	mConn := migrations.MCon(*c)
	return c, mConn
}

func (cfg *configuration) createTunnel() chan (bool) {
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
	ready := make(chan (bool))
	go tun.Start(ready)
	return ready
}

func getCfg(cfgFile string) configuration {
	fullPath, err := filepath.Abs(cfgFile)
	if err != nil {
		log.Fatalf("Failed to resolve path %v\n%v", cfgFile, errors.New(err).ErrorStack())
	}

	file, e := ioutil.ReadFile(fullPath)
	if e != nil {
		log.Fatalf("File error: %v\n", errors.New(e).ErrorStack())
	}
	cfg := configuration{}
	json.Unmarshal(file, &cfg)
	return cfg
}

func cfgDb(cfg configuration) *sql.DB {
	dbConn := cfg.Db.Username + `/` + cfg.Db.Password + `@localhost:` + strconv.Itoa(localPort) + "/" + cfg.DbName
	db, err := sql.Open("oci8", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database...\n%v", errors.New(err).ErrorStack())
	}
	return db
}
