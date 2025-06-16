package main

import (
	"context"
	"log"

	"github.com/ductong169z/shorten-url/config"
	"github.com/ductong169z/shorten-url/internal/server"
	"github.com/ductong169z/shorten-url/pkg/cache/redis"
	"github.com/ductong169z/shorten-url/pkg/database/mysql"
	"github.com/ductong169z/shorten-url/pkg/logger"

	_ "github.com/ductong169z/shorten-url/docs" // Swagger docs import
)

func main() {
	log.Println("Starting api server")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	ctx := context.Background()
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(ctx, "AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Repository
	mysqlDB, err := mysql.New(&cfg.MySQL)
	if err != nil {
		appLogger.Fatalf(ctx, "MySQL init: %s", err)
	}

	rdb, err := redis.NewClient(&cfg.Redis)
	if err != nil {
		appLogger.Fatalf(ctx, "RedisCluster init: %s", err)
	}

	s := server.NewServer(
		cfg,
		mysqlDB,
		server.Logger(appLogger),
		server.Redis(rdb),
	)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
