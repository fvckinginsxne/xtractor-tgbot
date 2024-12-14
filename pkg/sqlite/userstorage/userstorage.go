package userstorage

import (
	"database/sql"

	"bot/pkg/tech/e"
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

func (s UserStorage) UsernameByUserID(userID int64) (string, error) {
	q := `SELECT username FROM users WHERE id = ?`

	var username string

	if err := s.db.QueryRow(q, userID).Scan(&username); err != nil {
		return "", e.Wrap("can't get username by user id", err)
	}

	return username, nil
}
