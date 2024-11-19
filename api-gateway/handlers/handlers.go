package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/neokofg/go-pet-microservices/catalog-service/api/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	router     *gin.Engine
	logger     *zap.Logger
	catalogSvc proto.CatalogServiceClient
	//searchSvc    SearchServiceClient
	//recommendSvc RecommendServiceClient
}

func NewApp(
	logger *zap.Logger,
	catalogSvc proto.CatalogServiceClient,
	// searchSvc proto.SearchServiceClient,
	// recommendSvc proto.RecommendationServiceClient,
) *App {
	return &App{
		logger:     logger,
		catalogSvc: catalogSvc,
		//searchSvc:    searchSvc,
		//recommendSvc: recommendSvc,
	}
}

func (a *App) RegisterHandlers(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Catalog endpoints
		v1.GET("/items", a.handleGetItems)
		//v1.GET("/items/:id", a.handleGetItem)
		//v1.POST("/items", a.handleCreateItem)

		// Search endpoints
		//v1.GET("/search", a.handleSearch)

		// Recommendation endpoints
		//v1.GET("/recommendations", a.handleGetRecommendations)
	}
}

func (a *App) handleGetItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		handleError(c, err)
		return
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	if err != nil {
		handleError(c, err)
		return
	}

	resp, err := a.catalogSvc.GetItems(ctx, &proto.GetItemsRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

//func (a *App) handleSearch(c *gin.Context) {
//	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
//	defer cancel()
//
//	query := c.Query("q")
//	if query == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
//		return
//	}
//
//	resp, err := a.searchSvc.Search(ctx, &SearchRequest{
//		Query: query,
//	})
//
//	if err != nil {
//		handleError(c, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, resp)
//}
//
//func (a *App) handleGetRecommendations(c *gin.Context) {
//	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
//	defer cancel()
//
//	userID := c.GetString("user_id")
//	if userID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
//		return
//	}
//
//	resp, err := a.recommendSvc.GetRecommendations(ctx, &RecommendationRequest{
//		UserId: userID,
//	})
//
//	if err != nil {
//		handleError(c, err)
//	}
//
//	c.JSON(http.StatusOK, resp)
//}

func handleError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	switch st.Code() {
	case codes.NotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
	case codes.Unauthenticated:
		c.JSON(http.StatusUnauthorized, gin.H{"error": st.Message()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
