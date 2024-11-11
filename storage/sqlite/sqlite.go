package sqlite

import (
	"database/sql"

	"bot/lib/e"
	"bot/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &Storage{db: db}, nil
}

func (s Storage) Save(p *storage.Page) error {
	q := "INSERT INTO pages (url, user_name) VALUES (?, ?)"

	if _, err := s.db.Exec(q, p.URL, p.Username); err != nil {
		return e.Wrap("cant't save page", err)
	}

	return nil
}

func (s Storage) PickRandom(username string) (*storage.Page, error) {
	q := "SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1"

	var url string
	err := s.db.QueryRow(q, username).Scan(&url)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPage
	}

	if err != nil {
		return nil, e.Wrap("can't get random page", err)
	}

	return &storage.Page{
		URL:      url,
		Username: username,
	}, nil
}

func (s Storage) Remove(p *storage.Page) error {
	q := "DELETE FROM pages WHERE url = ? AND user_name = ?"

	if _, err := s.db.Exec(q, p.URL, p.Username); err != nil {
		return e.Wrap("can't remove page", err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	q := "SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?"

	var count int

	err := s.db.QueryRow(q, p.URL, p.Username).Scan(&count)
	if err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s Storage) Init() error {
	q := "CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)"

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table", err)
	}

	return nil
}
