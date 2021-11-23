package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/cpustejovsky/mongogo/routes"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Addr  string
	Uri   string
	Pprof string
}

var logger = &logrus.Logger{
	Out:       os.Stdout,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
}

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		logger.Info("No .env file found")
	}
}

func main() {
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
	logger.Info("Successfully connected to database!")

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: routes.Routes(logger, client),
	}
	logger.WithField(
		"address", cfg.Addr,
	).Info("Starting server")

	go func() {
		logger.Info(http.ListenAndServe(cfg.Pprof, nil))
	}()

	// Server Start
	err = srv.ListenAndServe()
	logger.Error(err)
}
