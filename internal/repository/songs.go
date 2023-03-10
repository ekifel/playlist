package repository

import (
	"context"
	"fmt"

	"github.com/ekifel/playlist/internal/model"
	"github.com/jackc/pgx/v5"
)

type songsRepo struct {
	db *pgx.Conn
}

func newSongsRepo(db *pgx.Conn) *songsRepo {
	return &songsRepo{
		db: db,
	}
}

func (s *songsRepo) GetSongs(ctx context.Context) ([]model.Song, error) {
	songs := []model.Song{}
	rows, err := s.db.Query(ctx, "select * from songs")
	if err != nil {
		return []model.Song{}, fmt.Errorf("error occurred while getting list of songs: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var id, duration int

		err = rows.Scan(&id, &name, &duration)
		if err != nil {
			return []model.Song{}, fmt.Errorf("error occurred while scanning row: %s", err.Error())
		}

		songs = append(songs, model.Song{
			ID:       id,
			Name:     name,
			Duration: duration,
		})
	}

	return songs, nil
}

func (s *songsRepo) GetSongByID(ctx context.Context, id int) (model.Song, error) {
	var song model.Song
	err := s.db.QueryRow(ctx, "SELECT id, name, duration FROM songs WHERE id = $1", id).Scan(&song.ID, &song.Name, &song.Duration)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Song{}, fmt.Errorf("song with ID %d not found", id)
		}

		return model.Song{}, fmt.Errorf("error retrieving song: %s", err.Error())
	}

	return song, nil
}

func (s *songsRepo) SaveSong(ctx context.Context, song model.Song) (int, error) {
	var id int
	err := s.db.QueryRow(ctx, "INSERT INTO songs (name, duration) VALUES ($1, $2) RETURNING id", song.Name, song.Duration).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating song: %s", err.Error())
	}

	return id, nil
}

func (s *songsRepo) UpdateSong(ctx context.Context, song model.Song) error {
	result, err := s.db.Exec(ctx, "UPDATE songs SET name = $1, duration = $2 WHERE id = $3", song.Name, song.Duration, song.ID)
	if err != nil {
		return fmt.Errorf("error updating song: %s", err.Error())
	}

	if result.RowsAffected() == 0 {
		return &ObjNotFound{fmt.Sprintf("song with id: %v doesn't exist", song.ID)}
	}

	return nil
}

func (s *songsRepo) DeleteSong(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, "DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting song: %s", err.Error())
	}

	if result.RowsAffected() == 0 {
		return &ObjNotFound{fmt.Sprintf("song with id: %v doesn't exist", id)}
	}

	return nil
}
