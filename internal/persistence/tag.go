package persistence

type PgTag struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}
