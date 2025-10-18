package persistence

type PgSerie struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}
