package userstorage

import (
	"database/sql"

	"bot/lib/e"
)

type UserStorage struct {
	db *sql.DB
}

func New(path string) (*UserStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &UserStorage{db: db}, nil
}

func (s UserStorage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL
	)`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table users", err)
	}

	return nil
}
