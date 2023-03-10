package repository

import (
	"context"

	"github.com/ekifel/playlist/internal/model"
	"github.com/jackc/pgx/v5"
)

type Songs interface {
	GetSongs(ctx context.Context) ([]model.Song, error)
	GetSongByID(ctx context.Context, id int) (model.Song, error)
	SaveSong(ctx context.Context, song model.Song) (int, error)
	UpdateSong(ctx context.Context, song model.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type Repositories struct {
	Songs Songs
}

func NewRepositories(db *pgx.Conn) *Repositories {
	return &Repositories{
		Songs: newSongsRepo(db),
	}
}
