package v1

import (
	"github.com/ekifel/playlist/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initSongsRoutes(v1)
		h.initPlaylistRoutes(v1)
	}
}
