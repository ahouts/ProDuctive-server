package migrations

import (
	"fmt"
	"log"
	"path/filepath"

	"sort"

	"strings"

	"regexp"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/go-errors/errors"
)

type MCon data.Conn

const m_history_name = "migration_history"

func (c *MCon) Up() {
	if !c.tableExist(m_history_name) {
		createMigrationsStr, err := migrationsCreate_migrationsSqlBytes()
		if err != nil {
			log.Fatalf("Failed to find create migrations command\n%v", errors.New(err).ErrorStack())
		}
		ctx := data.InitContext()
		c.ExecContext(ctx, fmt.Sprintf(string(createMigrationsStr), m_history_name))
		ctx.Done()
	}

	migs := getUpMigs()
	sort.Strings(migs)
	for _, mig := range migs {
		if !c.mExist(mig) {
			migBytes, err := Asset(mig)
			if err != nil {
				log.Fatalln(errors.New(err).ErrorStack())
			}
			migStr := string(migBytes)
			ctx := data.InitContext()
			_, err = c.ExecContext(ctx, migStr)
			if err != nil {
				log.Fatalln(mig, errors.New(err).ErrorStack())
			}
			ctx.Done()
			c.insertMig(mig)
			ctx = data.InitContext()
			_, err = c.ExecContext(ctx, "COMMIT")
			ctx.Done()
			if err != nil {
				log.Fatalln(errors.New(err).ErrorStack())
			}
		}
	}
}

func (c *MCon) Down() {
	migs := getDownMigs()
	sort.Sort(sort.Reverse(sort.StringSlice(migs)))
	for _, mig := range migs {
		if c.mExist(mig) {
			migBytes, err := Asset(mig)
			if err != nil {
				log.Fatalln(errors.New(err).ErrorStack())
			}
			migStr := string(migBytes)
			ctx := data.InitContext()
			_, err = c.ExecContext(ctx, migStr)
			if err != nil {
				log.Fatalln(errors.New(err).ErrorStack())
			}
			ctx.Done()
			c.removeMig(mig)
			ctx = data.InitContext()
			_, err = c.ExecContext(ctx, "COMMIT")
			ctx.Done()
			if err != nil {
				log.Fatalln(errors.New(err).ErrorStack())
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
	rows, err := c.Query(s)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Failed to serach tables.\n%v", errors.New(err).ErrorStack())
	}
	for rows.Next() {
		return true
	}
	return false
}

func (c *MCon) insertMig(migName string) {
	s := fmt.Sprintf("INSERT INTO %v VALUES('%v')", m_history_name, getFilename(migName))
	ctx := data.InitContext()
	_, err := c.ExecContext(ctx, s)
	ctx.Done()
	if err != nil {
		log.Fatalln(errors.New(err).ErrorStack())
	}
}

func (c *MCon) removeMig(migName string) {
	s := fmt.Sprintf("DELETE FROM %v WHERE mig = '%v'", m_history_name, getFilename(migName))
	ctx := data.InitContext()
	_, err := c.ExecContext(ctx, s)
	ctx.Done()
	if err != nil {
		log.Fatalln(errors.New(err).ErrorStack())
	}
}

func (c *MCon) getRunMigrations() []string {
	migs := make([]string, 0)
	query := fmt.Sprintf("SELECT mig FROM %v", m_history_name)
	ctx := data.InitContext()
	rows, err := c.QueryContext(ctx, query)
	defer ctx.Done()
	defer rows.Close()
	if err != nil {
		log.Fatalf("Failed to find existing migrations.\n%v", errors.New(err).ErrorStack())
	}
	for rows.Next() {
		var mig string
		if err := rows.Scan(&mig); err != nil {
			log.Fatalln(errors.New(err).ErrorStack())
		}
		migs = append(migs, mig)
	}
	return migs
}

func (c *MCon) mExist(migName string) bool {
	query := fmt.Sprintf(`select 'Y' from dual where exists (select 1 from %v where mig = '%v')`, m_history_name, getFilename(migName))
	ctx := data.InitContext()
	rows, err := c.QueryContext(ctx, query)
	defer ctx.Done()
	defer rows.Close()
	if err != nil {
		log.Fatalf("Failed to check if migration %v exists\n%v", migName, errors.New(err).ErrorStack())
	}
	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			log.Fatalln(errors.New(err).ErrorStack())
		}
		if res == "Y" {
			return true
		} else {
			log.Fatalf("%v is not 'Y'", res)
		}
	}
	return false
}
