package audiostorage

import (
	"database/sql"

	"bot/internal/service/extractor"
	"bot/pkg/postgres/urlstorage"
	"bot/pkg/tech/e"

	_ "github.com/lib/pq"
)

type AudioStorage struct {
	db         *sql.DB
	urlStorage *urlstorage.UrlStorage
}

func New(connURL string, urlStorage *urlstorage.UrlStorage) (*AudioStorage, error) {
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &AudioStorage{db: db, urlStorage: urlStorage}, nil
}

func (s *AudioStorage) Init() error {
	q := `create table if not exists audios(
			id serial primary key,
			audio_file bytea not null,
			title text not null,
			hash text not null,
			user_id integer not null,
			url_id integer not null,
			foreign key (user_id) references users(id),
			foreign key (url_id) references urls(id)
	);
	create index if not exists idx_audios_title on audios(title);
	create index if not exists idx_audios_url_id on audios(url_id);
	`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table audios", err)
	}

	return nil
}

func (s *AudioStorage) Save(audio *extractor.Audio, hash string, userID, urlID int64) (err error) {
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

	q := `insert into audios (audio_file, title, hash, user_id, url_id)
			values ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(q, audio.AudioFile, audio.Title, hash, userID, urlID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *AudioStorage) Remove(title string) (err error) {
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

	if err := s.urlStorage.Remove(tx, title); err != nil {
		return err
	}

	q := `delete from audios where title=$1 and 
			user_id in (select id from users)`

	if _, err := tx.Exec(q, title); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *AudioStorage) IsExists(videoURL string) (bool, error) {
	q := `select count(*) from audios a join users on a.user_id=users.id 
			join urls on a.url_id=urls.id where urls.video_url=$1`

	var count int

	if err := s.db.QueryRow(q, videoURL).Scan(&count); err != nil {
		return false, e.Wrap("can't check if audio is exists", err)
	}

	return count > 0, nil
}

func (s *AudioStorage) Title(hash string) (string, error) {
	q := `select title from audios where hash=$1`

	var title string

	if err := s.db.QueryRow(q, hash).Scan(&title); err != nil {
		return "", e.Wrap("can't get title by hash", err)
	}

	return title, nil
}

func (s *AudioStorage) Playlist(username string) ([]extractor.Audio, error) {
	q := `select a.audio_file, a.title from audios a
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

func scanRows(rows *sql.Rows) ([]extractor.Audio, error) {
	var audios []extractor.Audio

	for rows.Next() {
		var audio extractor.Audio

		err := rows.Scan(&audio.AudioFile, &audio.Title)
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
