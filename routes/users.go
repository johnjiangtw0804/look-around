package routes

import (
	api "look-around/external/api"
	"look-around/repository"
	"look-around/routes/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler interface {
	recommendEvents(ctx *gin.Context)
	listEvents(ctx *gin.Context)
	likeEvent(ctx *gin.Context)
	dislikeEvent(ctx *gin.Context)
}

func NewUserHandler(logger *zap.Logger, repo repository.UserRepo, mapUtil api.MapUtilities, eventSearcher api.EventsSearcher) UserHandler {
	return &userHandler{

		logger:        logger,
		userRepo:      repo,
		mapUtilities:  mapUtil,
		eventSearcher: eventSearcher,
	}
}

type userHandler struct {
	logger        *zap.Logger
	userRepo      repository.UserRepo
	mapUtilities  api.MapUtilities
	eventSearcher api.EventsSearcher
}

type listEventsResp struct {
	Events []entity.Event `json:"events"`
}

func (u *userHandler) listEvents(ctx *gin.Context) {
	userID, _ := ctx.Get(_ctxKey_UserID)
	user, err := u.userRepo.SelectUserByID(uuid.MustParse(userID.(string)))
	if err != nil {
		u.logger.Error("failed to get user", zap.Error(err))
		ctx.JSON(500, gin.H{"error": "failed to get user"})
		return
	}

	// get users lat long
	lat, long, err := u.mapUtilities.GetLatLong(user.Address)
	if err != nil {
		u.logger.Error("failed to get lat long", zap.Error(err))
		ctx.JSON(500, gin.H{"error": "failed to get lat long"})
		return
	}

	resp, err := u.eventSearcher.ListEvents(lat, long, 0, "")
	if err != nil {
		u.logger.Error("failed to get events", zap.Error(err))
		ctx.JSON(500, gin.H{"error": "failed to get events"})
		return
	}

	// compose responses
	eventsToOccur := make(map[*entity.Event]int)
	events := make([]entity.Event, 0)
	for _, event := range resp.Embedded.Events {
		var imageURL string
		if len(event.Images) > 0 {
			imageURL = event.Images[0].URL
		}
		long, _ := strconv.ParseFloat(event.EmbeddedVenues.Venues[0].Location.Longitude, 64)
		lat, _ := strconv.ParseFloat(event.EmbeddedVenues.Venues[0].Location.Latitude, 64)
		event := entity.Event{
			ID:   event.ID,
			Name: event.Name,
			Date: entity.Date{
				LocalDate: event.Dates.Start.LocalDate,
				LocalTime: event.Dates.Start.LocalTime,
				Timezone:  event.Dates.Timezone,
				Status:    event.Dates.Status.Code,
			},
			Address: event.EmbeddedVenues.Venues[0].Address.Line1 + ", " + event.
				EmbeddedVenues.Venues[0].City.Name + ", " +
				event.EmbeddedVenues.Venues[0].State.Name + ", " + event.EmbeddedVenues.Venues[0].Country.Name,
			Genres:    []string{event.Classifications[0].Genre.Name, event.Classifications[0].SubGenre.Name},
			ImageURL:  imageURL,
			URL:       event.URL,
			Longitude: long,
			Latitude:  lat,
		}

		// remove duplicates
		if _, ok := eventsToOccur[&event]; !ok {
			events = append(events, event)
		}
	}
	// calculate distance
	for i := range events {
		events[i].Distance, _ = u.mapUtilities.CalculateDistance(lat, long, events[i].Latitude, events[i].Longitude)
	}
	ctx.JSON(200, listEventsResp{Events: events})
}

type likeReq struct {
	Genre    string `json:"genre" binding:"required"`
	SubGenre string `json:"subgenre" binding:"required"`
}

type dislikeReq struct {
	Genre    string `json:"genre" binding:"required"`
	SubGenre string `json:"subgenre" binding:"required"`
}

func (u *userHandler) likeEvent(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	req := &likeReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("Warn: invalid request body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	genre := req.Genre
	subgenre := req.SubGenre
	u.userRepo.InsertUserLikeGenreAndSubGenre(uuid.MustParse(userID), genre, subgenre)

	ctx.JSON(201, gin.H{"message": "liked event"})
}

func (u *userHandler) dislikeEvent(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	req := &dislikeReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("Warn: invalid request body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	genre := req.Genre
	subgenre := req.SubGenre
	u.userRepo.InsertUserDisLikeGenreAndSubGenre(uuid.MustParse(userID), genre, subgenre)

	ctx.JSON(200, gin.H{"message": "disliked event"})
}

func (u *userHandler) recommendEvents(ctx *gin.Context) {

}
