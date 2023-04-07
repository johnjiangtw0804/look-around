package main

import (
	"log"
	"look-around/envconfig"
	"look-around/internal/database"
	"look-around/routes"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var env envconfig.Env
var logger *zap.Logger

func main() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Error: failed to initialize logger ", err)
	}
	err = envconfig.Process(&env)
	if err != nil {
		log.Fatal("Error: failed to load config from env vars ", err)
	}
	gin.SetMode(gin.ReleaseMode)

	db, err := database.NewGormDatabase(env.DATABASE_URL, env.Debug)
	if err != nil {
		log.Fatal("Error: failed to connect to database ", err)
	}
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Error: failed to migrate database ", err)
	}
	logger.Info("Finished migrating database")

	server := routes.Register(gin.Default(), logger, db)
	go func() {
		server.Run(":8080")
	}()
	logger.Info("Server started")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutdown servers...")
}
