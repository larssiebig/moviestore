package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

type Movie struct {
	Title string `json:"title"`
	Year  string `json:"year"`
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	body := Movie{}
	data, err := c.GetRawData()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Movie not defined: %v", err))
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Bad Input: %v", err))

	}

	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255),
		year VARCHAR(255)
	)`)

	if err != nil {
		log.Fatalf(fmt.Sprintf("Couldnt create table Movie: %v", err))
	}

	_, err = s.db.Exec(`INSERT INTO movies (title, year) VALUES ($1, $2)`, body.Title, body.Year)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Couldnt insert movie: %v", err))
	}

	return map[string]string{
		"message": "Movie added",
	}
}
