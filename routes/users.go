package routes

import (
	"look-around/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type userHandler struct {
	router   *gin.Engine
	logger   *zap.Logger
	userRepo repository.UserRepo
}

func NewUserHandler(router *gin.Engine, logger *zap.Logger, userRepo repository.UserRepo) *userHandler {
	return &userHandler{
		router:   router,
		logger:   logger,
		userRepo: userRepo,
	}
}
