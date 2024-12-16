package audiostorage

import (
	"database/sql"
	"log"

	"bot/internal/core"
	"bot/pkg/sqlite/userstorage"
	"bot/pkg/tech/e"

	_ "github.com/mattn/go-sqlite3"
)

type AudioStorage struct {
	db          *sql.DB
	userStorage *userstorage.UserStorage
}

func New(path string, userStrorage *userstorage.UserStorage) (*AudioStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &AudioStorage{db: db, userStorage: userStrorage}, nil
}

func (s AudioStorage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS audios (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			data BLOB NOT NULL,
			title TEXT NOT NULL,
			uuid TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	)`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table audios", err)
	}

	return nil
}

func (s AudioStorage) SaveAudio(audio *core.Audio, username string, uuid string) error {
	userID, err := s.userStorage.UserIdByUsername(username)
	if err != nil {
		return e.Wrap("can't save audio", err)
	}

	q := `INSERT INTO audios (url, data, title, uuid, user_id) 
		  VALUES (?, ?, ?, ?, ?)`

	_, err = s.db.Exec(q, audio.URL, audio.Data, audio.Title, uuid, userID)
	if err != nil {
		return e.Wrap("can't save audio", err)
	}

	return nil
}

func (s AudioStorage) RemoveAudio(title, username string) error {
	userId, err := s.userStorage.UserIdByUsername(username)
	if err != nil {
		return e.Wrap("can't remove audio", err)
	}

	log.Println("Audio title: ", title)
	log.Println("Username: ", username)

	q := `DELETE FROM audios WHERE user_id = ? AND title = ?`

	if _, err := s.db.Exec(q, userId, title); err != nil {
		return e.Wrap("can't remove audio", err)
	}

	return nil
}

func (s AudioStorage) IsExists(audio *core.Audio, username string) (bool, error) {
	userID, err := s.userStorage.UserIdByUsername(username)
	if err != nil {
		return false, e.Wrap("can't check if audio is exists", err)
	}

	q := `SELECT COUNT(*) FROM audios WHERE url = ? AND user_id = ?`

	var count int

	if err := s.db.QueryRow(q, audio.URL, userID).Scan(&count); err != nil {
		return false, e.Wrap("can't check if audio is exists", err)
	}

	return count > 0, nil
}

func (s AudioStorage) TitleAndUsernameByUUID(uuid string) (title, username string, err error) {
	q := `SELECT title, user_id FROM audios WHERE uuid = ?`

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

func (s AudioStorage) Playlist(username string) ([]core.Audio, error) {
	q := `SELECT a.url, a.data, a.title FROM audios a
		  JOIN users u ON a.user_id = u.id  WHERE u.username = ?`

	rows, err := s.db.Query(q, username)
	if err != nil {
		return nil, e.Wrap("can't join tables by user id", err)
	}
	defer func() { _ = rows.Close() }()

	var audios []core.Audio

	for rows.Next() {
		var audio core.Audio

		err := rows.Scan(&audio.URL, &audio.Data, &audio.Title)
		if err != nil {
			return nil, e.Wrap("can't scan audio", err)
		}

		audios = append(audios, audio)
	}

	if err = rows.Err(); err != nil {
		return nil, e.Wrap("error during rows iteration", err)
	}

	return audios, nil
}
