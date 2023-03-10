package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initPlaylistRoutes(api *gin.RouterGroup) {
	api.POST("/play", h.play)
	api.POST("/pause", h.pause)
	api.POST("/next", h.next)
	api.POST("/prev", h.prev)
}

func (h *Handler) play(c *gin.Context) {
	message, err := h.services.Playlist.Play()
	if err != nil {
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: message,
	})
}

func (h *Handler) pause(c *gin.Context) {
	message, err := h.services.Playlist.Pause()
	if err != nil {
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: message,
	})
}

func (h *Handler) next(c *gin.Context) {
	message, err := h.services.Playlist.Next()
	if err != nil {
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: message,
	})
}

func (h *Handler) prev(c *gin.Context) {
	message, err := h.services.Playlist.Prev()
	if err != nil {
		errorResponse(c, 500, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: message,
	})
}
