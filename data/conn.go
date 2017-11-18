package data

import (
	"context"
	"database/sql"
	"time"
)

type DbSession struct {
	DB *sql.DB
}

func initContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func (s *DbSession) InitTransaction() (*sql.Tx, error) {
	tx, err := s.DB.BeginTx(initContext(), nil)
	return tx, err
}
