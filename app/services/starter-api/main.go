package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/cpustejovsky/mongogo/foundation/logger"
	"github.com/cpustejovsky/mongogo/routes"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Config struct {
	Addr  string
	Uri   string
	Pprof string
}

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {

	log, err := logger.New("DEFAULT-API")
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// Flag and Config Setup
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":5000", "HTTP network address")
	flag.StringVar(&cfg.Uri, "uri", "mongodb://localhost:27017/mongogo", "MongoDB URI")
	flag.StringVar(&cfg.Pprof, "pprof", ":4000", "Pprof host and port")
	flag.Parse()

	// Environemntal Variables
	mongoUriFromEnv := os.Getenv("MONGO_URI")
	if mongoUriFromEnv != "" {
		cfg.Uri = mongoUriFromEnv
	}

	// DB Setup
	clientOptions := options.Client().
		ApplyURI(cfg.Uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	log.Infow("Successfully connected to database!")

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: routes.Routes(log, client),
	}

	// Server Start
	err = srv.ListenAndServe()
	log.Error(err)
	return nil
}
