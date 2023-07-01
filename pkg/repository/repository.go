package repository

import (
	"database/sql"

	"github.com/chema0/url-shortener/config"
)

type Repositories struct {
	URLs URLs
}

func New(db *sql.DB, cfg *config.DatabaseConfig) (*Repositories, error) {
	return &Repositories{
		URLs: NewURLs(db, cfg),
	}, nil
}
