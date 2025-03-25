package userstorage

import (
	"database/sql"

	_ "github.com/lib/pq"

	"bot/internal/tech/e"
)

type UserStorage struct {
	db *sql.DB
}

func New(connURL string) (*UserStorage, error) {
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &UserStorage{db: db}, nil
}

func (s *UserStorage) Init() error {
	q := `create table if not exists users(
			id serial primary key,
			username text not null
	);
	create index if not exists idx_users_username on users(username);
	`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table users", err)
	}

	return nil
}

func (s *UserStorage) Save(username string) (int64, error) {
	q := "select id from users where username=$1"

	var userID int64

	err := s.db.QueryRow(q, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.createUser(username)
		} else {
			return 0, e.Wrap("can't get user", err)
		}
	}

	return userID, nil
}

func (s *UserStorage) createUser(username string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, e.Wrap("can't create new user", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	q := `insert into users (username) values ($1) returning id`

	var userID int64

	if err := tx.QueryRow(q, username).Scan(&userID); err != nil {
		return 0, e.Wrap("can't create new user", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, e.Wrap("can't create new user", err)
	}

	return userID, nil
}

func (s *UserStorage) Username(hash string) (string, error) {
	q := `select u.username from users u join audios a on 
			u.id=a.user_id where a.hash=$1`

	var username string

	if err := s.db.QueryRow(q, hash).Scan(&username); err != nil {
		return "", e.Wrap("can't get username by hash", err)
	}

	return username, nil
}
