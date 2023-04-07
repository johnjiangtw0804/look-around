package routes

import (
	"look-around/internal/database"
	"look-around/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(
	router *gin.Engine,
	logger *zap.Logger,
	db *database.GormDatabase,
) *gin.Engine {
	// Create Repo instances
	userRepo := repository.NewUserRepo(db)

	// Register handlers for no authentication API
	auth := NewAuth(userRepo, logger)
	noAuthRouters := router.Group("")
	noAuthRouters.POST("/api/auth/register", auth.register)
	noAuthRouters.POST("/api/auth/login", auth.login)

	// router group to add middle ware for authentication TODO
	return router
}
