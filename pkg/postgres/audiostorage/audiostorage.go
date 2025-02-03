package audiostorage

import (
	"database/sql"

	"bot/internal/core"
	"bot/pkg/postgres/userstorage"
	"bot/pkg/tech/e"

	_ "github.com/lib/pq"
)

type AudioStorage struct {
	db          *sql.DB
	userStorage *userstorage.UserStorage
}

func New(connURL string, userStrorage *userstorage.UserStorage) (*AudioStorage, error) {
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &AudioStorage{db: db, userStorage: userStrorage}, nil
}

func (s *AudioStorage) Init() error {
	q := `create table if not exists audios(
			id serial primary key,
			url text not null,
			data bytea not null,
			title text not null,
			uuid text not null,
			user_id integer not null,
			foreign key (user_id) references users (id) on delete cascade
	)`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table audios", err)
	}

	return nil
}

func (s *AudioStorage) SaveAudio(audio *core.Audio, username string, uuid string) (err error) {
	defer func() { err = e.Wrap("can't save audio", err) }()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userID, err := s.userStorage.UserIdByUsername(username)
	if err != nil {
		return err
	}

	q := `insert into audios (url, data, title, uuid, user_id)
			values ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(q, audio.URL, audio.Data, audio.Title, uuid, userID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *AudioStorage) RemoveAudio(title, username string) (err error) {
	defer func() { err = e.Wrap("can't remove audio", err) }()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userId, err := s.userStorage.UserIdByUsername(username)
	if err != nil {
		return err
	}

	q := `delete from audios where user_id=$1 and title=$2`

	if _, err := s.db.Exec(q, userId, title); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *AudioStorage) IsExists(audio *core.Audio, username string) (bool, error) {
	userID, err := s.userStorage.UserIdByUsername(username)
	if err != nil {
		return false, e.Wrap("can't check if audio is exists", err)
	}

	q := `select count(*) from audios where url=$1 and user_id=$2`

	var count int

	if err := s.db.QueryRow(q, audio.URL, userID).Scan(&count); err != nil {
		return false, e.Wrap("can't check if audio is exists", err)
	}

	return count > 0, nil
}

func (s *AudioStorage) TitleAndUsernameByUUID(uuid string) (title, username string, err error) {
	q := `select title, user_id from audios where uuid = $1`

	var userID int64

	if err := s.db.QueryRow(q, uuid).Scan(&title, &userID); err != nil {
		return "", "", e.Wrap("can't get title by uuid", err)
	}

	username, err = s.userStorage.UsernameByUserID(userID)
	if err != nil {
		return "", "", e.Wrap("can't get username by uuid", err)
	}

	return title, username, nil
}

func (s *AudioStorage) Playlist(username string) ([]core.Audio, error) {
	q := `select a.url, a.data, a.title from audios a
		  join users u on a.user_id = u.id where u.username = $1`

	rows, err := s.db.Query(q, username)
	if err != nil {
		return nil, e.Wrap("can't get user's playlist", err)
	}
	defer func() { _ = rows.Close() }()

	audios, err := scanRows(rows)
	if err != nil {
		return nil, e.Wrap("can't get user's playlist", err)
	}

	return audios, nil
}

func scanRows(rows *sql.Rows) ([]core.Audio, error) {
	var audios []core.Audio

	for rows.Next() {
		var audio core.Audio

		err := rows.Scan(&audio.URL, &audio.Data, &audio.Title)
		if err != nil {
			return nil, e.Wrap("can't scan audio", err)
		}

		audios = append(audios, audio)
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap("error during rows iteration", err)
	}

	return audios, nil
}
