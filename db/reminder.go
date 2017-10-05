package db

import "time"

type Reminder struct {
	id         int
	message    string
	display_at time.Time
	created_at time.Time
	updated_at time.Time
}
