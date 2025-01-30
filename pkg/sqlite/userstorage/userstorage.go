package userstorage

import (
	"database/sql"
	"log"

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

func (s *UserStorage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL
	)`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table users", err)
	}

	return nil
}

func (s *UserStorage) UsernameByUserID(userID int64) (string, error) {
	q := `SELECT username FROM users WHERE id = ?`

	var username string

	if err := s.db.QueryRow(q, userID).Scan(&username); err != nil {
		return "", e.Wrap("can't get username by user id", err)
	}

	return username, nil
}

func (s *UserStorage) UserIdByUsername(username string) (int64, error) {
	q := `SELECT id FROM users WHERE username = ?`

	var userID int64

	err := s.db.QueryRow(q, username).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Println("create new user")

		userID, err = s.createUser(username)
		if err != nil {
			return 0, e.Wrap("can't get userID by username", err)
		}

		log.Println("user successfully created")
	} else if err != nil {
		return 0, e.Wrap("can't get userID by username", err)
	}

	return userID, nil
}

func (s *UserStorage) createUser(username string) (int64, error) {
	q := `INSERT INTO users (username) VALUES (?)`

	result, err := s.db.Exec(q, username)
	if err != nil {
		return 0, e.Wrap("can't create user", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, e.Wrap("can't create user", err)
	}

	return userID, nil
}
