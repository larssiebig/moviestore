package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.GET("/", s.HelloWorldHandler)
	r.GET("/movies", s.movieHandler)
	r.GET("/profile", s.profileHandler)

	r.POST("/movies", s.addMovieHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) movieHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.GetMovie())
}

func (s *Server) profileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Profile())
}

// function to post movies
func (s *Server) addMovieHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.AddMovie(c))
}
