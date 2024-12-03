package audiostorage

import (
	"database/sql"
	"log"

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

func (s AudioStorage) Save(a *core.Audio, username string) error {
	userID, err := s.GetOrCreateUserID(username)
	if err != nil {
		return e.Wrap("can't get user id", err)
	}

	q := `INSERT INTO audios (url, data, title, user_id) VALUES (?, ?, ?, ?)`

	if _, err = s.db.Exec(q, a.URL, a.Data, a.Title, userID); err != nil {
		return e.Wrap("can't save audio", err)
	}

	return nil
}

func (s AudioStorage) Remove(a *core.Audio) error {
	q := `DELETE FROM audios WHERE url = ?`

	if _, err := s.db.Exec(q, a.URL); err != nil {
		return e.Wrap("can't remove audio", err)
	}

	return nil
}

func (s AudioStorage) IsExists(a *core.Audio, username string) (bool, error) {
	userID, err := s.GetOrCreateUserID(username)
	if err != nil {
		return false, e.Wrap("can't check if audio is exists", err)
	}

	q := `SELECT COUNT(*) FROM audios WHERE url = ? AND user_id = ?`

	var count int

	if err = s.db.QueryRow(q, a.URL, userID).Scan(&count); err != nil {
		return false, e.Wrap("can't check if audio is exists", err)
	}

	return count > 0, nil
}

func (s AudioStorage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS audios (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			data BLOB NOT NULL,
			title TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	)`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table audios", err)
	}

	return nil
}

func (s AudioStorage) GetOrCreateUserID(username string) (int64, error) {
	q := `SELECT id FROM users WHERE username = ?`

	var userID int64

	err := s.db.QueryRow(q, username).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("create new user")

		q = `INSERT INTO users (username) VALUES (?)`

		result, err := s.db.Exec(q, username)
		if err != nil {
			return 0, err
		}

		userID, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}
