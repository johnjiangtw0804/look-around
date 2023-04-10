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
	authHandler := NewAuthHandler(logger, userRepo)
	noAuthRouters := router.Group("")
	noAuthRouters.POST("/api/auth/register", authHandler.register)
	noAuthRouters.POST("/api/auth/login", authHandler.login)

	authRouters := router.Group("", authenticate(userRepo, logger))
	// router group to add middle ware for authentication
	userHandler := NewUserHandler(logger, userRepo)

	authRouters.GET("/api/user/events", userHandler.listEvents)
	authRouters.POST("/api/user/events/:id/like", userHandler.likeEvent)
	authRouters.POST("/api/user/events/:id/dislike", userHandler.dislikeEvent)
	return router
}
