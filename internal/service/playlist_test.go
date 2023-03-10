package service

import (
	"testing"
)

func TestNewPlaylist(t *testing.T) {
	songService := &songService{}
	p := newPlaylist(songService)

	if p.mu == nil {
		t.Errorf("Mutex was not initialized")
	}

	if p.timer == nil {
		t.Errorf("Timer was not initialized")
	}

	if p.nextSong == nil {
		t.Errorf("nextSong channel was not initialized")
	}

	if p.prevSong == nil {
		t.Errorf("prevSong channel was not initialized")
	}

	if p.start == nil {
		t.Errorf("start channel was not initialized")
	}

	if p.isPlay {
		t.Errorf("isPlay should be false initially")
	}

	if p.songService != songService {
		t.Errorf("songService was not set properly")
	}
}

func TestPlaylist_Play(t *testing.T) {
	p := newPlaylist(&songService{})
	p.isPlay = true

	_, err := p.Play()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	p.isPlay = false
	p.songs.head = nil
	_, err = p.Play()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPlaylist_Pause(t *testing.T) {
	p := newPlaylist(&songService{})

	_, err := p.Pause()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if !p.timer.IsPaused() {
		t.Errorf("Timer should be paused, but it is not")
	}
}

func TestPlaylist_Next(t *testing.T) {
	p := newPlaylist(&songService{})

	_, err := p.Next()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPlaylist_Prev(t *testing.T) {
	p := newPlaylist(&songService{})

	_, err := p.Prev()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
