package main

import (
	"context"
	"entgo.io/ent/dialect"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/neokofg/go-pet-microservices/catalog-service/api/proto"
	"github.com/neokofg/go-pet-microservices/catalog-service/ent"
	"github.com/neokofg/go-pet-microservices/catalog-service/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	client, err := ent.Open(dialect.Postgres, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		logger.Fatal("Failed to create schema", zap.Error(err))
	}

	catalogService := service.NewCatalogService(client, logger)

	grpcServer := grpc.NewServer()
	proto.RegisterCatalogServiceServer(grpcServer, catalogService)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	grpcAddr := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("Starting gRPC server", zap.String("addr", grpcAddr))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
		Handler: router,
	}

	go func() {
		logger.Info("Starting HTTP server", zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("HTTP server forced to shutdown", zap.Error(err))
	}

	grpcServer.GracefulStop()

	logger.Info("Servers exited properly")
}
