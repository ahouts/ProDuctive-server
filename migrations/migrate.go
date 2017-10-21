package migrations

import (
	"fmt"
	"log"
	"path/filepath"

	"sort"

	"strings"

	"regexp"

	"github.com/ahouts/ProDuctive-server/data"
)

type MCon data.Conn

const m_history_name = "migration_history"

func (c *MCon) Up() {
	if !c.tableExist(m_history_name) {
		createMigrationsStr, err := migrationsCreate_migrationsSqlBytes()
		if err != nil {
			log.Fatalf("Failed to find create migrations command\n%v", err)
		}
		c.ExecContext(c.Ctx, fmt.Sprintf(string(createMigrationsStr), m_history_name))
	}

	migs := getUpMigs()
	sort.Strings(migs)
	for _, mig := range migs {
		if !c.mExist(mig) {
			migBytes, err := Asset(mig)
			if err != nil {
				log.Fatalln(err)
			}
			migStr := string(migBytes)
			for _, migPart := range strings.Split(migStr, ";") {
				_, err = c.ExecContext(c.Ctx, migPart)
				if err != nil {
					log.Fatalln(err)
				}
			}
			c.insertMig(migStr)
			_, err = c.ExecContext(c.Ctx, "COMMIT")
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func (c *MCon) Down() {
	migs := getDownMigs()
	sort.Sort(sort.Reverse(sort.StringSlice(migs)))
	for _, mig := range migs {
		if !c.mExist(mig) {
			migBytes, err := Asset(mig)
			if err != nil {
				log.Fatalln(err)
			}
			migStr := string(migBytes)
			for _, migPart := range strings.Split(migStr, ";") {
				_, err = c.ExecContext(c.Ctx, migPart)
				if err != nil {
					log.Fatalln(err)
				}
			}
			c.removeMig(migStr)
			_, err = c.ExecContext(c.Ctx, "COMMIT")
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func getUpMigs() []string {
	allMigs := AssetNames()
	reg := regexp.MustCompile(`migrations/\d{5}_.+\.up\.sql`)
	upMigs := make([]string, 0)
	for _, mig := range allMigs {
		if reg.MatchString(mig) {
			upMigs = append(upMigs, mig)
		}
	}
	return upMigs
}

func getDownMigs() []string {
	allMigs := AssetNames()
	reg := regexp.MustCompile(`migrations/\d{5}_.+\.down\.sql`)
	downMigs := make([]string, 0)
	for _, mig := range allMigs {
		if reg.MatchString(mig) {
			downMigs = append(downMigs, mig)
		}
	}
	return downMigs
}

func getFilename(fullName string) string {
	_, name := filepath.Split(fullName)
	parts := strings.Split(name, "_")
	return parts[0]
}

func (c *MCon) tableExist(tablename string) bool {
	s := fmt.Sprintf("select table_name from user_tables where table_name='%v'", strings.ToUpper(tablename))
	rows, err := c.QueryContext(c.Ctx, s)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Failed to serach tables.\n%v", err)
	}
	for rows.Next() {
		return true
	}
	return false
}

func (c *MCon) insertMig(migName string) {
	s := fmt.Sprintf("INSERT INTO migration_history VALUES('%v')", getFilename(migName))
	_, err := c.ExecContext(c.Ctx, s)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *MCon) removeMig(migName string) {
	s := fmt.Sprintf("DELETE FROM migration_history WHERE mig = '%v'", getFilename(migName))
	_, err := c.ExecContext(c.Ctx, s)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *MCon) getRunMigrations() []string {
	migs := make([]string, 0)
	query := "SELECT mig FROM migration_history"
	rows, err := c.QueryContext(c.Ctx, query)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Failed to find existing migrations.\n%v", err)
	}
	for rows.Next() {
		var mig string
		if err := rows.Scan(&mig); err != nil {
			log.Fatalln(err)
		}
		migs = append(migs, mig)
	}
	return migs
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
			log.Fatalf("%v is not 'Y'", res)
		}
	}
	return false
}
