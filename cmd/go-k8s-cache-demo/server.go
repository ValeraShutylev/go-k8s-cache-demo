package main

import (
	"cache-demo/internal/handlers"
	"cache-demo/internal/repository"
	"cache-demo/internal/routes"
	"cache-demo/internal/services"
	"cache-demo/internal/storage"
	"cache-demo/internal/routes/metrics"
	"cache-demo/utils"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port = os.Getenv("PORT")
)


type serverConfigs struct {
	port             string
}

func Run() {

	state := utils.NewGlobalState()
	serverConfigs := getConfigs()

	storage := storage.NewStorage()
	memRepository := repository.NewMemoryCache(storage)

	cacheService := services.NewCacheService(memRepository)
	cacheHandler := handlers.NewCacheHandler(cacheService)

	go memRepository.UploadDataToCacheFromStorage(state)

	go memRepository.CleanExpiredCacheItems()

	go metrics.ScrapCacheLenghtMetrics(memRepository)

	server := &http.Server{
		Addr:    serverConfigs.port,
		Handler: routes.NewCacheRouter(cacheHandler, state).GinRouter,
		ErrorLog: slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelError),
	}

	start(server)
}


func getConfigs() *serverConfigs {
	if port == "" {
		port = "8080"
	}
	return &serverConfigs{
		port: fmt.Sprintf(":%s", port),
	}
}

func start(server *http.Server) {
	slog.Info(
		"Started the server with graceful shutdown", 
		slog.String("Port", server.Addr),
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(
				"Cannot start server", 
				slog.String("Error", err.Error()),
			)
		}
	}()

	go func() {
		slog.Info(
			"Starting pprof server", 
			slog.Int("Port", 6060),
		)
		http.ListenAndServe(":6060", nil)
	}()


	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Graceful Shutdown Server in 3 seconds")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error(
			"Shutdown server with error", 
			slog.String("Error", err.Error()),
		)
	}

	<-ctx.Done()
	slog.Info("Server shutdown")
}