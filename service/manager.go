package service


import (
	"database/sql"
)

type repo interface {
	Exec(query string, args ...any)  (sql.Result, error)
}

type Servece struct {
	repository repo
}

func LoadRepo(repo repo) *Servece {
	return &Servece{
		repository: repo,
	}
}


