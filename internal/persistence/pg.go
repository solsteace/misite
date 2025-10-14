package persistence

import "github.com/jmoiron/sqlx"

type Pg struct {
	db *sqlx.DB
}

func NewPg(db *sqlx.DB) Pg {
	return Pg{db: db}
}
