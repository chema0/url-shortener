package repository

import (
	"database/sql"
	"time"

	"github.com/chema0/url-shortener/config"
)

type URL struct {
	ID        int64     `db:"id" json:"id"`
	URL       string    `db:"url" json:"url"`
	Hash      string    `db:"hash" json:"hash"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
	OwnerID   string    `db:"owner_id" json:"owner_id"`
}

type URLs interface {
	CreateURL(url URL) (URL, error)
	ListURLs(ownerID string) ([]URL, error)
	DeleteURL(id string) error
	GetURLByHash(hash string) (URL, error)
}

type URLsRepository struct {
	db  *sql.DB
	cfg *config.DatabaseConfig
}

func NewURLs(db *sql.DB, cfg *config.DatabaseConfig) URLs {
	return &URLsRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *URLsRepository) CreateURL(url URL) (URL, error) {
	// r.db.Exec("INSERT INTO url (url, hash, expired_at, owner_id) VALUES ($1, $2, $3, $4) RETURNING ")
	row, err := r.db.Query("INSERT INTO url (url, hash, expired_at, owner_id) VALUES ($1, $2, $3, $4) RETURNING id, url, hash, created_at, expired_at, owner_id",
		url.URL, url.Hash, url.ExpiredAt, url.OwnerID)

	// FIXME: check url return on error
	if err != nil {
		return url, err
	}

	row.Scan(&url)
	return url, err
}

// TODO: implement
func (r *URLsRepository) ListURLs(ownerID string) ([]URL, error) {
	return []URL{}, nil
}

// TODO: implement
func (r *URLsRepository) DeleteURL(id string) error {
	return nil
}

// TODO: implement
func (r *URLsRepository) GetURLByHash(hash string) (URL, error) {
	return URL{}, nil
}
