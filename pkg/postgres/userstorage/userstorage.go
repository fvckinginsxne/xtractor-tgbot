package userstorage

import (
	"database/sql"
	"log"

	"bot/pkg/tech/e"

	_ "github.com/lib/pq"
)

type UserStorage struct {
	db *sql.DB
}

func New(connURL string) (*UserStorage, error) {
	log.Println("Connecting to database with URL:", connURL)

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
	q := `create table if not exists users (
			id serial primary key,
			username text not null
	)`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table users", err)
	}

	return nil
}

func (s *UserStorage) UsernameByUserID(userID int64) (string, error) {
	q := `select username from users where id=$1`

	var username string

	if err := s.db.QueryRow(q, userID).Scan(&username); err != nil {
		return "", e.Wrap("can't get username by user id", err)
	}

	return username, nil
}

func (s *UserStorage) UserIdByUsername(username string) (int64, error) {
	q := `select id from users where username=$1`

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
	q := `insert into users (username) values ($1) returning id`

	var userID int64
	err := s.db.QueryRow(q, username).Scan(&userID)
	if err != nil {
		return 0, e.Wrap("can't create user", err)
	}

	return userID, nil
}
