package db

import "time"

type User struct {
	email         string
	password_hash string
	created_at    time.Time
	updated_at    time.Time
}
