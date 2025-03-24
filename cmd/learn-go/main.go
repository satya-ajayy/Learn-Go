package main

import (
	// Go Internal Packages
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	// Local Packages
	config "learn-go/config"
	xhttp "learn-go/http"
	handlers "learn-go/http/handlers"
	mongodb "learn-go/repositories/mongodb"
	redis "learn-go/repositories/redis"
	health "learn-go/services/health"
	orders "learn-go/services/orders"
	students "learn-go/services/students"

	// External Packages
	"github.com/alecthomas/kingpin/v2"
	_ "github.com/jsternberg/zap-logfmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"go.uber.org/zap"
)

// InitializeServer sets up an HTTP server with defined handlers.
// Repositories are initialized, creates the services, and subsequently constructs
// handlers for the services
func InitializeServer(ctx context.Context, k config.Config, logger *zap.Logger) (*xhttp.Server, error) {
	// Mongo Connection
	mongoClient, err := mongodb.Connect(ctx, k.Mongo.URI)
	if err != nil {
		return nil, err
	}

	// Redis Connection
	redisClient, err := redis.Connect(ctx, logger, k.Redis.URI, k.Redis.Password)
	if err != nil {
		return nil, err
	}

	// Init repos, services && handlers
	studentsRepo := mongodb.NewStudentsRepository(mongoClient)
	ordersRepo := redis.NewOrdersRepository(redisClient)

	healthSvc := health.NewService(logger, mongoClient)
	studentsSvc := students.NewService(studentsRepo)
	ordersSvc := orders.NewService(ordersRepo)

	studentsHandler := handlers.NewStudentsHandler(studentsSvc)
	ordersHandler := handlers.NewOrdersHandler(ordersSvc)

	server := xhttp.NewServer(k.Prefix, logger, studentsHandler, ordersHandler, healthSvc)
	return server, nil
}

// LoadConfig loads the default configuration and overrides it with the config file
// specified by the path defined in the config flag
func LoadConfig() *koanf.Koanf {
	configPathMsg := "Path to the application config file"
	configPath := kingpin.Flag("config", configPathMsg).Short('c').Default("config.yml").String()

	kingpin.Parse()
	k := koanf.New(".")
	_ = k.Load(rawbytes.Provider(config.DefaultConfig), yaml.Parser())
	if *configPath != "" {
		_ = k.Load(file.Provider(*configPath), yaml.Parser())
	}

	return k
}

func main() {
	k := LoadConfig()

	// Unmarshalling config into struct
	appKonf := config.Config{}
	err := k.Unmarshal("", &appKonf)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Validate the config loaded
	if err = appKonf.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	if !appKonf.IsProdMode {
		k.Print()
	}

	cfg := zap.NewProductionConfig()
	cfg.Encoding = "logfmt"
	_ = cfg.Level.UnmarshalText([]byte(appKonf.Logger.Level))
	cfg.InitialFields = make(map[string]any)
	cfg.InitialFields["host"], _ = os.Hostname()
	cfg.InitialFields["service"] = appKonf.Application
	cfg.OutputPaths = []string{"stdout"}
	logger, _ := cfg.Build()
	defer func() {
		_ = logger.Sync()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, err := InitializeServer(ctx, appKonf, logger)
	if err != nil {
		logger.Fatal("Cannot initialize server", zap.Error(err))
	}
	if err := srv.Listen(ctx, appKonf.Listen); err != nil {
		logger.Fatal("Cannot listen", zap.Error(err))
	}
}
