package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct{}

func SetupStats(r *gin.RouterGroup) *StatsHandler {
	s := &StatsHandler{}
	s.initRoutes(r)
	return s
}

func (s *StatsHandler) initRoutes(r *gin.RouterGroup) {
	stats := r.Group("stats")

	stats.GET("/ping", s.ping)
}

func (s *StatsHandler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now().Unix(),
	})
}
