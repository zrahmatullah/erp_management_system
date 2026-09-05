package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"cafe-erp-system/backend/internal/config"
	appHttp "cafe-erp-system/backend/internal/delivery/http"
	"cafe-erp-system/backend/internal/delivery/http/handler"
	"cafe-erp-system/backend/internal/repository/postgres"
	"cafe-erp-system/backend/internal/usecase/auth"
	"cafe-erp-system/backend/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	logger.InitLogger(cfg.Server.Env)
	logger.Log.Info().Msg("🚀 Starting Cafe ERP Backend Server...")

	// PostgreSQL Pool
	dbpool, err := pgxpool.New(context.Background(), cfg.DB.DSN)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer dbpool.Close()

	// Verify DB Connection
	if err := dbpool.Ping(context.Background()); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database ping failed")
	}
	logger.Log.Info().Msg("✅ PostgreSQL Connected Successfully")

	// Optional Redis Client (graceful if not running)
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.URL})
	defer redisClient.Close()
	ctxRedis, cancelRedis := context.WithTimeout(context.Background(), 1*time.Second)
	if err := redisClient.Ping(ctxRedis).Err(); err != nil {
		logger.Log.Warn().Err(err).Msg("⚠️ Redis not reachable, continuing in standalone mode")
	} else {
		logger.Log.Info().Msg("✅ Redis Connected Successfully")
	}
	cancelRedis()

	// Repositories & Services
	userRepo := postgres.NewUserRepository(dbpool)
	authUsecase := auth.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)
	masterHandler := handler.NewMasterHandler(dbpool)
	opHandler := handler.NewOperationalHandler(dbpool)

	// Router
	router := appHttp.SetupRouter(authHandler, masterHandler, opHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		logger.Log.Info().Msgf("🌐 Cafe ERP API Server listening on http://localhost:%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Server failed to listen")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info().Msg("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server shutdown error")
	}
	logger.Log.Info().Msg("Server stopped.")
}
