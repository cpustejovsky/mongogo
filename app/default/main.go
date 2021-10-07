package main

import (
	"log"
	"os"

	"github.com/cpustejovsky/mongogo/internal/data/sys/database"
	"github.com/joho/godotenv"
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

	cfg := config{
		db{
			uri: os.Getenv("MONGO_URI"),
		},
	}

	//likely move to DB
	database.Open(database.Config{URI: cfg.db.uri})

}
