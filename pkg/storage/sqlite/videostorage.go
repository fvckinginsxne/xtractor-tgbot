package sqlite

import (
	"database/sql"

	"bot/internal/core"
	"bot/lib/e"

	_ "github.com/mattn/go-sqlite3"
)

type AudioStorage struct {
	db *sql.DB
}

func New(path string) (*AudioStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &AudioStorage{db: db}, nil
}

func (s AudioStorage) Save(p *core.Audio) error {
	q := "INSERT INTO audios (url, user_name) VALUES (?, ?)"

	if _, err := s.db.Exec(q, p.URL, p.Username); err != nil {
		return e.Wrap("cant't save page", err)
	}

	return nil
}

func (s AudioStorage) Remove(p *core.Audio) error {
	q := "DELETE FROM audios WHERE url = ? AND user_name = ?"

	if _, err := s.db.Exec(q, p.URL, p.Username); err != nil {
		return e.Wrap("can't remove page", err)
	}

	return nil
}

func (s AudioStorage) IsExists(p *core.Audio) (bool, error) {
	q := "SELECT COUNT(*) FROM audios WHERE url = ? AND user_name = ?"

	var count int

	err := s.db.QueryRow(q, p.URL, p.Username).Scan(&count)
	if err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s AudioStorage) Init() error {
	q := "CREATE TABLE IF NOT EXISTS audios (url TEXT, user_name TEXT)"

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table", err)
	}

	return nil
}
