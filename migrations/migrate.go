package migrations

import (
	"fmt"
	"log"
	"path/filepath"

	"sort"

	"io/ioutil"

	"github.com/ahouts/ProDuctive-server/data"
)

type MCon data.Conn

func (c *MCon) RunMigrations() {
	c.ExecContext(c.Ctx, `
CREATE TABLE migration_history (
	mig VARCHAR(50) PRIMARY KEY
)
`)
	files, err := filepath.Glob("./migrations/*.sql")
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(files)
	for _, file := range files {
		if !c.mExist(file) {
			migBytes, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatalln(err)
			}
			mig := string(migBytes)
			_, err = c.ExecContext(c.Ctx, mig)
			if err != nil {
				log.Fatalln(err)
			}
			c.insertMig(file)
		}
	}
}

func getFilename(fullName string) string {
	_, name := filepath.Split(fullName)
	return name
}

func (c *MCon) insertMig(migName string) {
	s := fmt.Sprintf("INSERT INTO migration_history VALUES('%v')", getFilename(migName))
	_, err := c.ExecContext(c.Ctx, s)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *MCon) mExist(migName string) bool {
	query := fmt.Sprintf(`select 'Y' from dual where exists (select 1 from migration_history where mig = '%v')`, getFilename(migName))
	rows, err := c.QueryContext(c.Ctx, query)
	if err != nil {
		log.Fatalf("Failed to check if migration %v exists\n%v", migName, err)
	}
	defer rows.Close()
	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			log.Fatalln(err)
		}
		if res == "Y" {
			return true
		} else {
			log.Fatalf("%v is not 'T'", res)
		}
	}
	return false
}
