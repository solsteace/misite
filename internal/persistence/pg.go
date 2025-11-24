package persistence

import "github.com/jmoiron/sqlx"

type Pg struct {
	db *sqlx.DB
}

func NewPg(db *sqlx.DB) Pg {
	return Pg{db: db}
}

type ExplorationQueryParam struct {
	Page    int
	Limit   int
	Include struct {
		Keyword []string
		Tag     []int
		Serie   []int
	}
	Exclude struct {
		Keyword []string
		Tag     []int
		Serie   []int
	}
}
