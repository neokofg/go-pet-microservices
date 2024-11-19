package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/neokofg/go-pet-microservices/api-gateway/handlers"
	"github.com/neokofg/go-pet-microservices/api-gateway/middleware"
	"github.com/neokofg/go-pet-microservices/catalog-service/api/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.CORSMiddleware(),
		middleware.RequestLoggerMiddleware(logger),
		middleware.PrometheusMiddleware(),
		middleware.RateLimiterMiddleware(),
	)

	catalogConn, err := initGRPCClient(os.Getenv("CATALOG_SERVICE_ADDR"))
	if err != nil {
		logger.Fatal("Failed to connect to catalog service", zap.Error(err))
	}
	defer catalogConn.Close()

	catalogClient := proto.NewCatalogServiceClient(catalogConn)

	app := handlers.NewApp(
		logger,
		catalogClient,
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	app.RegisterHandlers(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		logger.Info("Starting server...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}

func initGRPCClient(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
