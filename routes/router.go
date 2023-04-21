package routes

import (
	"look-around/envconfig"
	api "look-around/external/api"
	"look-around/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(
	router *gin.Engine,
	logger *zap.Logger,
	db *repository.GormDatabase,
	env *envconfig.Env,
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
	userHandler := NewUserHandler(logger, userRepo, api.NewMapUtilities(env.GOOGLE_MAP_API_KEY), api.NewEventsSearcher(env.TICKET_MASTER_API_KEY))

	authRouters.GET("/api/user/events", userHandler.listEvents)
	authRouters.POST("/api/user/events/like", userHandler.likeEvent)
	authRouters.POST("/api/user/events/dislike", userHandler.dislikeEvent)
	authRouters.GET("/api/user/events/recommend", userHandler.recommendEvents)
	return router
}
