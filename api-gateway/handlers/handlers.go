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
}

func NewApp(
	logger *zap.Logger,
	catalogSvc proto.CatalogServiceClient,
) *App {
	return &App{
		logger:     logger,
		catalogSvc: catalogSvc,
	}
}

func (a *App) RegisterHandlers(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Catalog endpoints
		v1.GET("/items", a.handleGetItems)
		v1.GET("/items/:id", a.handleGetItem)
		v1.POST("/items", a.handleCreateItem)
		v1.PUT("/items/:id", a.handleUpdateItem)
		v1.DELETE("/items/:id", a.handleDeleteItem)
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

func (a *App) handleGetItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, _ := c.Params.Get("id")

	resp, err := a.catalogSvc.GetItem(ctx, &proto.GetItemRequest{
		Id: id,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

type CreateItemRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
	ImageUrl    string   `json:"imageUrl" binding:"required"`
}

func (a *App) handleCreateItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, err)
		return
	}

	resp, err := a.catalogSvc.CreateItem(ctx, &proto.CreateItemRequest{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		ImageUrl:    req.ImageUrl,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

type UpdateItemRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ImageUrl    string   `json:"imageUrl"`
}

func (a *App) handleUpdateItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, _ := c.Params.Get("id")

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, err)
		return
	}

	resp, err := a.catalogSvc.UpdateItem(ctx, &proto.UpdateItemRequest{
		Id:          id,
		Title:       &req.Title,
		Description: &req.Description,
		Tags:        req.Tags,
		ImageUrl:    &req.ImageUrl,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *App) handleDeleteItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	id, _ := c.Params.Get("id")
	_, err := a.catalogSvc.DeleteItem(ctx, &proto.DeleteItemRequest{
		Id: id,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"message": "Item deleted"})
}

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
