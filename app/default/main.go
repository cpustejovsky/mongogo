package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cpustejovsky/mongogo/foundation/logger"
	"github.com/cpustejovsky/mongogo/internal/data/sys/database"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("No .env file found")
	}
}

type db struct {
	uri string `mongodb://localhost:27017/example`
}

type config struct {
	db
}

func main() {
	// Construct the application logger.
	log, err := logger.New("SALES-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg := config{
		db{
			uri: os.Getenv("MONGO_URI"),
		},
	}

	client, err := database.Open(database.Config{URI: cfg.db.uri})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	client.Database("example")
	log.Infow("Database connected")
	return nil
}
