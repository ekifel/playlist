package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ekifel/playlist/internal/model"
	"github.com/sirupsen/logrus"
)

type Playlist struct {
	songs       SongList
	mu          *sync.Mutex
	timer       *Timer
	nextSong    chan bool
	prevSong    chan bool
	start       chan bool
	isPlay      bool
	songService *songService
}

type SongList struct {
	head    *model.Song
	current *model.Song
	tail    *model.Song
}

func newPlaylist(songService *songService) *Playlist {
	return &Playlist{
		mu:          &sync.Mutex{},
		timer:       newTimer(),
		nextSong:    make(chan bool),
		prevSong:    make(chan bool),
		start:       make(chan bool),
		isPlay:      false,
		songService: songService,
	}
}

func (p *Playlist) Run() error {
	savedSongs, err := p.songService.GetSongs(context.Background())
	if err != nil {
		return fmt.Errorf("error initialy running playlist: %s", err.Error())
	}

	for _, song := range savedSongs {
		newSong := &model.Song{ID: song.ID, Name: song.Name, Duration: song.Duration}
		if p.songs.head == nil {
			p.songs.head = newSong
			p.songs.tail = newSong
		} else {
			p.songs.tail.Next = newSong
			newSong.Prev = p.songs.tail
			p.songs.tail = newSong
			newSong.Next = p.songs.head
			p.songs.head.Prev = newSong
		}
	}

	go func() {
		for {
			if p.isPlay {
				if p.timer.IsPaused() {
					p.timer.duration = time.Duration(p.songs.current.Duration-p.timer.timePassed) * time.Second
				} else {
					p.timer.duration = time.Duration(p.songs.current.Duration) * time.Second
				}

				p.timer.Start()
				logrus.Infof("Song '%s' started to play\n", p.songs.current.Name)
			}

			select {
			case <-p.start:
				p.isPlay = true
				if p.songs.current == nil {
					p.songs.current = p.songs.head
				}

			case <-p.timer.Done():
				logrus.Infof("Song '%s' is over!\n", p.songs.current.Name)

				if p.songs.current.Next != nil {
					p.songs.current = p.songs.current.Next
				} else {
					p.songs.current = p.songs.head
				}

			case <-p.timer.Paused():
				p.isPlay = false
				logrus.Infof("Song '%s' was stopped, time left: %v\n",
					p.songs.current.Name, p.songs.current.Duration-int(p.timer.timePassed))

			case <-p.nextSong:
				if p.songs.current != nil {
					if p.songs.current.Next != nil {
						p.songs.current = p.songs.current.Next
					}

				} else {
					if p.songs.head != nil {
						if p.songs.head.Next != nil {
							p.songs.current = p.songs.head.Next
						}
					}
				}

				if !p.timer.paused {
					p.timer.Stop()
				}
				logrus.Infof("Switch song to '%s'", p.songs.current.Name)

			case <-p.prevSong:
				if p.songs.current != nil {
					if p.songs.current.Prev != nil {
						p.songs.current = p.songs.current.Prev
					}
				} else {
					if p.songs.head != nil {
						if p.songs.head.Prev != nil {
							p.songs.current = p.songs.head.Prev
						}
					}
				}

				if !p.timer.paused {
					p.timer.Stop()
				}
				logrus.Infof("Switch song to '%s'", p.songs.current.Name)
			}
		}
	}()

	return nil
}

func (p *Playlist) Play() (string, error) {
	if p.isPlay {
		return "", fmt.Errorf("The music is already playing")
	}

	if p.songs.head == nil {
		return "", fmt.Errorf("There are no songs in the playlist")
	}

	p.start <- true

	return "Music started to play", nil
}

func (p *Playlist) Pause() (string, error) {
	if !p.isPlay {
		return "", fmt.Errorf("The music was already stopped")
	}

	p.timer.Pause()

	return "Music was stopped", nil
}

func (p *Playlist) Next() (string, error) {
	if p.songs.head == nil {
		return "", fmt.Errorf("There are no songs in the playlist")
	}

	p.nextSong <- true

	return "Current track switched to the next song", nil
}

func (p *Playlist) Prev() (string, error) {
	if p.songs.head == nil {
		return "", fmt.Errorf("There are no songs in the playlist")
	}

	p.prevSong <- true

	return "Current track switched to the prev song", nil
}

func (p *Playlist) AddSong(song model.Song) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.songs.head == nil {
		p.songs.head = &song
		p.songs.tail = &song
	} else {
		p.songs.tail.Next = &song
		song.Prev = p.songs.tail
		p.songs.tail = &song
		song.Next = p.songs.head
		p.songs.head.Prev = &song
	}

	return fmt.Sprintf("Song '%s' was added to playlist", song.Name), nil
}

func (p *Playlist) UpdateSong(ctx context.Context, song model.Song) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.songs.current != nil {
		if p.songs.current.ID == song.ID {
			return "", fmt.Errorf("Impossible to delete a song that continuing to play")
		}
	}

	err := p.songService.UpdateSong(ctx, song)
	if err != nil {
		return "", err
	}

	s := p.songs.head
	for s.ID != song.ID {
		s = s.Next
	}

	s.Name = song.Name
	s.Duration = song.Duration

	return "Song was updated", nil
}

func (p *Playlist) DeleteSong(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.songs.current != nil {
		if p.songs.current.ID == id {
			return fmt.Errorf("Impossible to delete a song that continuing to play")
		}
	}

	err := p.songService.DeleteSong(ctx, id)
	if err != nil {
		return err
	}

	s := p.songs.head
	for s.ID != id {
		s = s.Next
	}

	if s == p.songs.head {
		if s.Next != nil {
			p.songs.head = s.Next
			p.songs.tail.Next = p.songs.head
		} else {
			p.songs.head = nil
		}

		return nil
	}

	if s.Prev != nil && s.Next != nil {
		s.Prev.Next = s.Next
	}

	if s.Next != nil {
		s.Next.Prev = s.Prev
	} else {
		s.Prev.Next = nil
	}

	s.Prev = nil
	s.Next = nil

	return nil
}
