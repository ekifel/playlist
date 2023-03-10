package service

import (
	"context"

	"github.com/ekifel/playlist/internal/model"
	"github.com/ekifel/playlist/internal/repository"
)

type ISongs interface {
	GetSongs(ctx context.Context) ([]model.Song, error)
	GetSongByID(ctx context.Context, id int) (model.Song, error)
	SaveSong(ctx context.Context, song model.Song) (int, error)
	UpdateSong(ctx context.Context, song model.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type IPlaylist interface {
	Run() error
	Play() (string, error)
	Pause() (string, error)
	Next() (string, error)
	Prev() (string, error)
	AddSong(song model.Song) (string, error)
	DeleteSong(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, song model.Song) (string, error)
}

type Services struct {
	Songs    ISongs
	Playlist IPlaylist
}

func NewServices(deps Deps) *Services {
	songService := newSongService(deps.Repos.Songs)

	return &Services{
		Songs:    songService,
		Playlist: newPlaylist(songService),
	}
}

type Deps struct {
	Repos *repository.Repositories
}
