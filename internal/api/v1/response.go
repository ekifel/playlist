package v1

import (
	"github.com/ekifel/playlist/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type response struct {
	Message string `json:"message"`
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}

type songsResponse struct {
	Songs []model.Song `json:"songs"`
}
