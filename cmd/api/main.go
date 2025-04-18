package main

import (
	"context"
	"log"

	"github.com/ductong169z/blog-api/config"
	"github.com/ductong169z/blog-api/internal/server"
	"github.com/ductong169z/blog-api/pkg/database/mysql"
	"github.com/ductong169z/blog-api/pkg/logger"
)

func main() {
	log.Println("Starting api server")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	// token, _ := utils.GenerateJWTToken(&models.User{
	// 	UserID: uuid.New(),
	// }, cfg)

	// fmt.Println("token: ", token)

	ctx := context.Background()
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(ctx, "AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Repository
	mysqlDB, err := mysql.New(&cfg.MySQL)
	if err != nil {
		appLogger.Fatalf(ctx, "MySQL init: %s", err)
	}

	// rdb, err := redis.NewClient(&cfg.Redis)
	// if err != nil {
	// 	appLogger.Fatalf(ctx, "RedisCluster init: %s", err)
	// }

	s := server.NewServer(
		cfg,
		mysqlDB,
		server.Logger(appLogger),
		// server.Redis(rdb),
	)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
