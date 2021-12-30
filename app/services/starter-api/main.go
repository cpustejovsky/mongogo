package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/ardanlabs/conf/v2"
	"github.com/cpustejovsky/mongogo/foundation/logger"
	"github.com/cpustejovsky/mongogo/routes"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/automaxprocs/maxprocs"

	"go.uber.org/zap"
)

/*
	TODO: Make sure environmental varibles can be loaded without godotenv and the conditional that replaces the conf value
*/

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("No .env file found")
	}
}

var build = "develop"

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
	// =========================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// Configuration
	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3001"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			Uri             string        `conf:"default:mongodb://localhost:27017/mongogo,mask"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "cpustejovsky MIT license",
		},
	}

	const prefix = "DEFAULT"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	mongoUriFromEnv := os.Getenv("MONGO_URI")
	if mongoUriFromEnv != "" {
		cfg.Web.Uri = mongoUriFromEnv
	}

	// DB Setup
	clientOptions := options.Client().
		ApplyURI(cfg.Web.Uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	log.Infow("Successfully connected to database!")

	srv := &http.Server{
		Addr:    cfg.Web.APIHost,
		Handler: routes.Routes(log, client),
	}

	// Server Start
	err = srv.ListenAndServe()
	log.Error(err)
	return nil
}
