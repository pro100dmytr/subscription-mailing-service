package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"subscription-mailing-service/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	r := gin.Default()
	userGroup := r.Group("/users")
	{
		userGroup.GET("/:id", h.GetUserID())
		userGroup.GET("/", h.GetAllUsers())
		userGroup.POST("/", h.CreateUser())
		userGroup.PUT("/:id", h.UpdateUser())
		userGroup.DELETE("/:id", h.DeleteUser())
	}

	if err := r.Run(cfg.Server.Address); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
