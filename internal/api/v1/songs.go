package v1

import (
	"net/http"
	"strconv"

	"github.com/ekifel/playlist/internal/model"
	"github.com/ekifel/playlist/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

type songInput struct {
	Name     string `json:"name" binding:"required"`
	Duration int    `json:"duration" binding:"required"`
}

func (h *Handler) initSongsRoutes(api *gin.RouterGroup) {
	api.GET("/songs", h.getAllSongs)
	api.POST("/songs", h.addSong)
	api.GET("/songs/:id", h.getSongByID)
	api.PATCH("/songs/:id", h.updateSong)
	api.DELETE("/songs/:id", h.deleteSong)
}

func (h *Handler) getAllSongs(c *gin.Context) {
	ctx := context.Background()

	songs, err := h.services.Songs.GetSongs(ctx)
	if err != nil {
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, songsResponse{
		Songs: songs,
	})
}

func (h *Handler) addSong(c *gin.Context) {
	ctx := context.Background()

	var inp songInput
	if err := c.BindJSON(&inp); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	song := model.Song{
		Name:     inp.Name,
		Duration: inp.Duration,
	}

	id, err := h.services.Songs.SaveSong(ctx, song)
	if err != nil {
		errorResponse(c, 500, err.Error())

		return
	}

	song.ID = id

	message, err := h.services.Playlist.AddSong(song)
	if err != nil {
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: message,
	})
}

func (h *Handler) getSongByID(c *gin.Context) {
	ctx := context.Background()

	id, _ := c.Params.Get("id")
	songID, _ := strconv.Atoi(id)

	song, err := h.services.Songs.GetSongByID(ctx, songID)
	if err != nil {
		if err == pgx.ErrNoRows {
			errorResponse(c, 404, err.Error())
		}
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, songsResponse{
		Songs: []model.Song{song},
	})
}

func (h *Handler) updateSong(c *gin.Context) {
	ctx := context.Background()

	id, _ := c.Params.Get("id")
	songID, _ := strconv.Atoi(id)

	var inp songInput
	if err := c.BindJSON(&inp); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	song := model.Song{
		ID:       songID,
		Name:     inp.Name,
		Duration: inp.Duration,
	}

	message, err := h.services.Playlist.UpdateSong(ctx, song)
	if err != nil {
		if objNotFound, ok := err.(*repository.ObjNotFound); ok {
			errorResponse(c, 404, objNotFound.Msg)

			return
		}
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: message,
	})
}

func (h *Handler) deleteSong(c *gin.Context) {
	ctx := context.Background()

	id, _ := c.Params.Get("id")
	songID, _ := strconv.Atoi(id)

	err := h.services.Playlist.DeleteSong(ctx, songID)
	if err != nil {
		if objNotFound, ok := err.(*repository.ObjNotFound); ok {
			errorResponse(c, 404, objNotFound.Msg)

			return
		}
		errorResponse(c, 500, err.Error())

		return
	}

	c.Status(http.StatusOK)
}
