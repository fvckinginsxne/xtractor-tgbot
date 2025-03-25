package urlstorage

import (
	"database/sql"

	_ "github.com/lib/pq"

	"bot/internal/tech/e"
)

type UrlStorage struct {
	db *sql.DB
}

func New(connURL string) (*UrlStorage, error) {
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &UrlStorage{db: db}, nil
}

func (s *UrlStorage) Init() error {
	q := `create table if not exists urls(
			id serial primary key,
			video_url text not null
	);
	create index if not exists idx_urls_video_url on urls(video_url);
	`

	if _, err := s.db.Exec(q); err != nil {
		return e.Wrap("can't create table urls", err)
	}

	return nil
}

func (s *UrlStorage) Save(url string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, e.Wrap("can't save video url", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	q := `insert into urls (video_url) values ($1) returning id`

	var urlID int64

	if err := tx.QueryRow(q, url).Scan(&urlID); err != nil {
		return 0, e.Wrap("can't save video url", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, e.Wrap("can't save video url", err)
	}

	return urlID, nil
}

func (s *UrlStorage) Remove(tx *sql.Tx, title string) (err error) {
	q := `delete from urls where id in (
			select a.url_id from audios a join users u on a.user_id=u.id
				where a.title=$1
	)`

	if _, err := tx.Exec(q, title); err != nil {
		return e.Wrap("can't remove video url", err)
	}

	return nil
}
