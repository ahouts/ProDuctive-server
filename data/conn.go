package data

import (
	"context"
	"database/sql"
)

type Conn struct {
	sql.DB
	Ctx context.Context
}
