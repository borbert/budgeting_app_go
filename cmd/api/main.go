package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/borbert/budgeting_app_go/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func geEnvVars() {
	err := godotenv.Load("creds.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	var cfg config
	// geEnvVars()

	// flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	// flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	// flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DATABASE_URL"), "Postgres connection string")

	// flag.Parse()
	cfg.port = os.Getenv("PORT")
	cfg.env = os.Getenv("env")
	cfg.jwt.secret = os.Getenv("BUDGET_JWT")
	cfg.db.dsn = os.Getenv("DATABASE_URL")

	fmt.Println(cfg)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	err = srv.ListenAndServe()

	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
