package main

import (
	"log"
	"look-around/envconfig"
	"look-around/internal/database"

	"github.com/gin-gonic/gin"
)

var env envconfig.Env

func main() {
	var err error
	err = envconfig.Process(&env)
	if err != nil {
		log.Fatal("Error: failed to load config from env vars ", err)
	}
	gin.SetMode(gin.ReleaseMode)
	// logger.Printf("Connecting to database... %s\n", env.DB_DSN)

	db, err := database.NewGormDatabase(env.DATABASE_URL, env.Debug)
	if err != nil {
		log.Fatal("Error: failed to connect to database ", err)
	}
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Error: failed to migrate database ", err)
	}
	log.Println("Finished migrating database")
}
