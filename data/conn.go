package data

import (
	"context"
	"database/sql"
	"time"

	"gopkg.in/rana/ora.v4"
)

const dbPrefetchRowCount = 50000

type Conn struct {
	sql.DB
}

func InitContext() context.Context {
	// Set timeout
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// Set prefetch count
	ctx = ora.WithStmtCfg(ctx, ora.Cfg().StmtCfg.SetPrefetchRowCount(dbPrefetchRowCount))
	return ctx
}

func (c *Conn) InitTransaction() (*sql.Tx, error) {
	ctx := InitContext()
	return c.BeginTx(ctx, nil)
}
