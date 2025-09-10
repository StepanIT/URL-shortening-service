package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

type URLData struct {
	ShortURL    string
	OriginalURL string
}

type PostrgesStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostrgesStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	s := &PostrgesStorage{db: db}
	if err := s.initSchema(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *PostrgesStorage) initSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url TEXT UNIQUE NOT NULL,
			original_url TEXT UNIQUE NOT NULL,
			user_id TEXT NOT NULL
		);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostrgesStorage) Save(shortURL, originalURL, userID string) (string, error) {
	_, err := s.db.Exec(
		`INSERT INTO urls (short_url, original_url, user_id)
         VALUES ($1, $2, $3)`,
		shortURL, originalURL, userID,
	)
	if err != nil {
		// проверяем, что это ошибка уникальности
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			// достаем уже существующий short_url
			var existingShort string
			row := s.db.QueryRow(
				"SELECT short_url FROM urls WHERE original_url = $1",
				originalURL,
			)
			if scanErr := row.Scan(&existingShort); scanErr != nil {
				return "", scanErr
			}
			return existingShort, fmt.Errorf("url already exists")
		}
		return "", err
	}
	return shortURL, nil
}

func (s *PostrgesStorage) Get(shortURL string) (string, error) {
	var originalURL string
	err := s.db.QueryRow(
		"SELECT original_url FROM urls WHERE short_url = $1",
		shortURL,
	).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return originalURL, err
}

func (s *PostrgesStorage) GetAllByUser(userID string) (map[string]string, error) {
	rows, err := s.db.Query(
		"SELECT short_url, original_url FROM urls WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make(map[string]string)
	for rows.Next() {
		var shortURL, originalURL string
		if err := rows.Scan(&shortURL, &originalURL); err != nil {
			return nil, err
		}
		urls[shortURL] = originalURL
	}
	return urls, rows.Err()
}

func (s *PostrgesStorage) Ping() error {
	return s.db.Ping()
}
