package entitystore

import "database/sql"

type txOrDB interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...interface{}) (sql.Result, error)
}
