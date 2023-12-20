package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	GetMovie() map[string]string
	Profile() map[string]string
	AddMovie(ctx *gin.Context) map[string]string
}

type service struct {
	db *sql.DB
}

type Movie struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title"`
	Year  string `json:"year"`
}

func newMovie(id int, title, year string) Movie {
	return Movie{id, title, year}
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func ConnectDB() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	return s
}

func (s *service) GetMovie() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "Those are your movies",
	}
}

func (s *service) Profile() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "Thats your profile",
	}
}

func (s *service) AddMovie(c *gin.Context) map[string]string {
	var newMovie Movie
	if err := c.BindJSON(&newMovie); err != nil {
		log.Fatalf(fmt.Sprintf("invalid JSON: %v", err))
		makeGinResponse(c, http.StatusNotFound, "invalid JSON")
	}

	sql := `INSERT INTO movies (title, year) VALUES ($1, $2)`
	result, err := s.db.Exec(sql, newMovie.Title, newMovie.Year)
	if err != nil {
		log.Fatalf(fmt.Sprintf("unable to execute the query: %v", err))
		makeGinResponse(c, http.StatusInternalServerError, "unable to execute the query")
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Fatalf(fmt.Sprintf("unable to get affected rows: %v", err))
		makeGinResponse(c, http.StatusInternalServerError, "unable to get affected rows")
	}

	if n == 0 {
		err := "could not insert movie"
		log.Fatalf(fmt.Sprintf("could not insert movie: %v", err))
		makeGinResponse(c, http.StatusInternalServerError, err)
	}

	m := "Movie added"
	makeGinResponse(c, http.StatusOK, m)

	return map[string]string{
		"message": m,
	}
}

func makeGinResponse(c *gin.Context, statusCode int, value string) {
	c.JSON(statusCode, gin.H{
		"message":    value,
		"statusCode": statusCode,
	})
}
