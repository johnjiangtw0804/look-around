package routes

import (
	"look-around/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	listEvents(ctx *gin.Context)
	likeEvent(ctx *gin.Context)
	dislikeEvent(ctx *gin.Context)
}

func NewUserHandler(logger *zap.Logger, repo repository.UserRepo) UserHandler {
	return &userHandler{
		logger:   logger,
		userRepo: repo,
	}
}

type userHandler struct {
	logger   *zap.Logger
	userRepo repository.UserRepo
}

func (u *userHandler) listEvents(ctx *gin.Context) {
	

}

func (u *userHandler) likeEvent(ctx *gin.Context) {

}

func (u *userHandler) dislikeEvent(ctx *gin.Context) {

}
