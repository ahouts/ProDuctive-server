package data

import "time"

type Note struct {
	id         int
	name       string
	content    string
	created_at time.Time
	updated_at time.Time
}
