package migrations

import (
	"fmt"
	"log"
	"path/filepath"

	"sort"

	"strings"

	"regexp"

	"database/sql"

	"github.com/ahouts/ProDuctive-server/data"
	"github.com/go-errors/errors"
)

const mHistoryName = "migration_history"

func Up(s *data.DbSession) {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatalf("failed to init context %v", err)
	}
	if !tableExist(tx, mHistoryName) {
		createMigrationsStr, err := migrationsCreate_migrationsSqlBytes()
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to find create migrations command\n%v", errors.New(err).ErrorStack())
		}
		_, err = tx.Exec(fmt.Sprintf(string(createMigrationsStr), mHistoryName))
		if err != nil {
			tx.Rollback()
			log.Fatalln(errors.New(err).ErrorStack())
		}
	}

	migs := getUpMigs()
	sort.Strings(migs)
	for _, mig := range migs {
		if !mExist(tx, mig) {
			migBytes, err := Asset(mig)
			if err != nil {
				tx.Rollback()
				log.Fatalln(errors.New(err).ErrorStack())
			}
			migStr := string(migBytes)
			_, err = tx.Exec(migStr)
			if err != nil {
				tx.Rollback()
				log.Fatalln(mig, errors.New(err).ErrorStack())
			}
			insertMig(tx, mig)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("failed to commit changes, %v", err)
	}
}

func Down(s *data.DbSession) {
	migs := getDownMigs()
	sort.Sort(sort.Reverse(sort.StringSlice(migs)))

	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatalf("failed to init context %v", err)
	}
	for _, mig := range migs {
		if mExist(tx, mig) {
			migBytes, err := Asset(mig)
			if err != nil {
				tx.Rollback()
				log.Fatalln(errors.New(err).ErrorStack())
			}
			migStr := string(migBytes)
			_, err = tx.Exec(migStr)
			if err != nil {
				tx.Rollback()
				log.Fatalln(errors.New(err).ErrorStack())
			}
			removeMig(tx, mig)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to commit changes, %v", err)
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

func tableExist(tx *sql.Tx, tablename string) bool {
	s := fmt.Sprintf("select table_name from user_tables where table_name='%v'", strings.ToUpper(tablename))
	_, err := tx.Query(s)
	if err != nil {
		return false
	}
	return true
}

func insertMig(tx *sql.Tx, migName string) {
	s := fmt.Sprintf("INSERT INTO %v VALUES('%v')", mHistoryName, getFilename(migName))
	_, err := tx.Exec(s)
	if err != nil {
		tx.Rollback()
		log.Fatalln(errors.New(err).ErrorStack())
	}
}

func removeMig(tx *sql.Tx, migName string) {
	s := fmt.Sprintf("DELETE FROM %v WHERE mig = '%v'", mHistoryName, getFilename(migName))
	_, err := tx.Exec(s)
	if err != nil {
		tx.Rollback()
		log.Fatalln(errors.New(err).ErrorStack())
	}
}

func getRunMigrations(tx *sql.Tx) []string {
	migs := make([]string, 0)
	query := fmt.Sprintf("SELECT mig FROM %v", mHistoryName)
	rows, err := tx.Query(query)
	defer rows.Close()
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to find existing migrations.\n%v", errors.New(err).ErrorStack())
	}
	for rows.Next() {
		var mig string
		if err := rows.Scan(&mig); err != nil {
			tx.Rollback()
			log.Fatalln(errors.New(err).ErrorStack())
		}
		migs = append(migs, mig)
	}
	return migs
}

func mExist(tx *sql.Tx, migName string) bool {
	query := fmt.Sprintf(`select 'Y' from dual where exists (select 1 from %v where mig = '%v')`, mHistoryName, getFilename(migName))
	rows, err := tx.Query(query)
	defer rows.Close()
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to check if migration %v exists\n%v", migName, errors.New(err).ErrorStack())
	}
	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			tx.Rollback()
			log.Fatalln(errors.New(err).ErrorStack())
		}
		if res == "Y" {
			return true
		} else {
			tx.Rollback()
			log.Fatalf("%v is not 'Y'", res)
		}
	}
	return false
}
