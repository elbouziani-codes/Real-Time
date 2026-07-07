package repository  

import (
		"database/sql"
		"context"
)


type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	db DBTX 
}


func New(db DBTX) *Queries   { // repository(for all cases)
		return &Queries{ db: db} 
}


