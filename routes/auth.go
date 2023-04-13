package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	repo "look-around/repository"
	schema "look-around/repository/schema"
	"look-around/token"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler interface {
	login(ctx *gin.Context)
	register(ctx *gin.Context)
}

type authHandler struct {
	repo   repo.UserRepo
	logger *zap.Logger
}

func NewAuthHandler(logger *zap.Logger, repo repo.UserRepo) AuthHandler {
	return &authHandler{
		repo:   repo,
		logger: logger,
	}
}

const _ctxKey_UserID = "userID"
const jwtExpPeriod = 7 * 24 * time.Hour

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (a *authHandler) login(ctx *gin.Context) {
	req := &loginReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("Warn: invalid request body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user, err := a.repo.SelectUserByUsername(req.Username)
	if err != nil {
		a.logger.Warn("Warn: user not found")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := checkPassword(user.Password, req.Password); err != nil {
		a.logger.Warn("Warn: invalid password")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	signedToken, err := token.GenJWT(
		user.ID.String(),
		user.UserName,
		(float64)(time.Now().Add(jwtExpPeriod).Unix()))
	if err != nil {
		a.logger.Error("failed to sign jwt", zap.String("user", req.Username), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}

func checkPassword(storedPassword, loginPassword string) error {
	if storedPassword == "" || loginPassword == "" {
		return errors.New("given password(s) is empty")
	}
	passwordBuf := bytes.Buffer{}
	passwordBuf.WriteString(loginPassword)
	passwordBuf.WriteString(repo.Salt)
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), passwordBuf.Bytes())
}

type signUpReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address"`
}

func (a *authHandler) register(ctx *gin.Context) {
	req := &signUpReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("Warn: invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	json, _ := json.Marshal(req)
	a.logger.Info("Info: user req" + string(json))

	if err := a.repo.InsertUser(schema.User{
		UserName: req.Username,
		Password: req.Password,
		Gender:   req.Gender,
		Age:      req.Age,
		Email:    req.Email,
		Phone:    req.Phone,
		Address:  req.Address,
	}); err != nil {
		a.logger.Error("Error: failed to insert user", zap.String("user", req.Username), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to insert user"})
		return
	}
}

func authenticate(repo repo.UserRepo, logger *zap.Logger) gin.HandlerFunc {
	authorizationHeaderField := "Authorization"
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get(authorizationHeaderField)
		prefix := "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			logger.Warn("Warn: invalid token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, prefix)
		if tokenStr == "" {
			logger.Warn("WARN: token not found")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userClaims, err := token.ParseJWT(tokenStr)
		if err == token.ErrJWTExpired {
			logger.Warn("WARN: token is expire: " + tokenStr)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if err != nil {
			logger.Warn("Warn: general error: " + err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if _, err := uuid.Parse(userClaims.UserID); err != nil {
			logger.Warn("Warn: invalid user id: " + userClaims.UserID + " " + err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(_ctxKey_UserID, userClaims.UserID)
	}
}
