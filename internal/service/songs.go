package service

import (
	"context"

	"github.com/ekifel/playlist/internal/model"
	"github.com/ekifel/playlist/internal/repository"
)

type songService struct {
	songsRepo repository.Songs
}

func newSongService(songsRepo repository.Songs) *songService {
	return &songService{
		songsRepo: songsRepo,
	}
}

func (s *songService) GetSongs(ctx context.Context) ([]model.Song, error) {
	songs, err := s.songsRepo.GetSongs(ctx)
	if err != nil {
		return []model.Song{}, err
	}

	return songs, nil
}

func (s *songService) GetSongByID(ctx context.Context, id int) (model.Song, error) {
	song, err := s.songsRepo.GetSongByID(ctx, id)
	if err != nil {
		return model.Song{}, err
	}

	return song, nil
}

func (s *songService) SaveSong(ctx context.Context, song model.Song) (int, error) {
	id, err := s.songsRepo.SaveSong(ctx, song)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *songService) UpdateSong(ctx context.Context, song model.Song) error {
	err := s.songsRepo.UpdateSong(ctx, song)
	if err != nil {
		return err
	}

	return nil
}

func (s *songService) DeleteSong(ctx context.Context, id int) error {
	err := s.songsRepo.DeleteSong(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
