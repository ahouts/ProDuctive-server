package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/ahouts/ProDuctive-server/migrations"
	"github.com/go-errors/errors"
)

type dbConfig struct {
	Username string
	Password string
	Hostname string
	Port     uint
	Name     string
}

type configuration struct {
	DB dbConfig
}

func serve(cfgFile string, port int) {
	s := initDb(cfgFile)
	migrations.Up(s)

	setupRoutes(s)

	fmt.Printf("Environment successfully configured, beginning to serve on port %v.\n", strconv.Itoa(port))

	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(port), "cert.pem", "key.pem", nil))
}

func dropDb(cfgFile string) {
	s := initDb(cfgFile)
	migrations.Down(s)
}

func initDb(cfgFile string) *data.DbSession {
	cfg := getCfg(cfgFile)

	db := cfgDb(cfg)

	s := &data.DbSession{DB: db}
	return s
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
	dbConn := cfg.DB.Username + `/` + cfg.DB.Password + `@ ` + cfg.DB.Hostname + `:` + strconv.FormatInt(int64(cfg.DB.Port), 10) + "/" + cfg.DB.Name
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database...\n%v", errors.New(err).ErrorStack())
	}
	return db
}
