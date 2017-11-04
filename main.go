package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	_ "github.com/mattn/go-oci8"

	"io/ioutil"

	"encoding/json"

	"strconv"

	"path/filepath"

	"net/http"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/ahouts/ProDuctive-server/migrations"
	"github.com/ahouts/ProDuctive-server/tunnel"
	"github.com/emicklei/go-restful"
	"github.com/go-errors/errors"
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

const localPort = 40841
const oraclePort = 1521
const sshPort = 22
const webPort = 444

func main() {
	cfgFile, err := ask.Ask("config file? (./config.json): ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	if cfgFile == "" {
		cfgFile = "./config.json"
	}
	cfg := getCfg(cfgFile)

	// this function returns a channel that returns a value once its done
	// by pulling a value out of the channel, we block until the tunnel is ready
	<-cfg.createTunnel()

	db := cfgDb(cfg)
	defer db.Close()

	c := &data.Conn{DB: *db}
	mConn := migrations.MCon(*c)
	mConn.Up()

	restful.Add(userWs(c))

	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(webPort), "cert.pem", "key.pem", nil))
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
		log.Fatalf("Failed to connect to database...\n%v", err, errors.New(err).ErrorStack())
	}
	return db
}
