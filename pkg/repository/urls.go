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
	GetURLByHash(hash string) (URL, error)
	DeleteURL(id int64) error
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

const (
	hashLength = 8
)

func (r *URLsRepository) CreateURL(url URL) (URL, error) {
	url.Hash = generateHash(hashLength)
	url.ExpiredAt = time.Now().AddDate(0, 1, -1)

	query := `
		INSERT INTO urls (url, hash, expired_at, owner_id) VALUES ($1, $2, $3, $4) 
		RETURNING id, url, hash, created_at, expired_at, owner_id;`
	err := r.db.QueryRow(query,
		url.URL, url.Hash, url.ExpiredAt, url.OwnerID).Scan(&url.ID, &url.URL, &url.Hash, &url.CreatedAt, &url.ExpiredAt, &url.OwnerID)

	if err != nil {
		return url, err
	}

	return url, err
}

func (r *URLsRepository) ListURLs(ownerID string) ([]URL, error) {
	rows, err := r.db.Query("SELECT * FROM urls WHERE owner_id = $1;", ownerID)

	urls := []URL{}
	defer rows.Close()
	for rows.Next() {
		var url URL
		err := rows.Scan(&url.ID, &url.URL, &url.Hash, &url.CreatedAt, &url.ExpiredAt, &url.OwnerID)
		if err != nil {
			return []URL{}, err
		}
		urls = append(urls, url)
	}

	return urls, err
}

func (r *URLsRepository) GetURLByHash(hash string) (URL, error) {
	url := URL{}
	err := r.db.QueryRow("SELECT * FROM urls WHERE hash = $1;",
		hash).Scan(&url.ID, &url.URL, &url.Hash, &url.CreatedAt, &url.ExpiredAt, &url.OwnerID)
	return url, err
}

func (r *URLsRepository) DeleteURL(id int64) error {
	_, err := r.db.Exec("DELETE FROM urls WHERE id = $1;", id)
	return err
}
