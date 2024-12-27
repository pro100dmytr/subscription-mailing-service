package subscription

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"subscription-mailing-service/internal/model"
	subs "subscription-mailing-service/storage/subscriber"
)

type SubscriberHandler interface {
	GetSubscriberID() gin.HandlerFunc
	GetAllSubscribers() gin.HandlerFunc
	GetSubscribersByLevel() gin.HandlerFunc
	CreateSubscriber() gin.HandlerFunc
	UpdateSubscriber() gin.HandlerFunc
	UpdateSubscriberLevel() gin.HandlerFunc
	DeleteSubscriber() gin.HandlerFunc
}

type Handler struct {
	store  *subs.SubscriberStorage
	logger *slog.Logger
}

func NewHandler(store *subs.SubscriberStorage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) GetSubscriberID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		subscriberID, err := strconv.Atoi(idStr)

		if err != nil {
			h.logger.Error("Invalid subscriber ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscriber ID"})
			return
		}

		subscriber, err := h.store.Get(c.Request.Context(), subscriberID)
		if err != nil {
			h.logger.Error("Error getting subscriber", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting subscriber"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": subscriber})
	}
}

func (h *Handler) GetAllSubscribers() gin.HandlerFunc {
	return func(c *gin.Context) {
		subscribers, err := h.store.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error getting subscribers", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting subscribers"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": subscribers})
	}
}
func (h *Handler) CreateSubscriber() gin.HandlerFunc {
	return func(c *gin.Context) {
		var subscriber *model.Subscriber

		if err := c.ShouldBindJSON(&subscriber); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err := h.store.Create(c.Request.Context(), subscriber)
		if err != nil {
			h.logger.Error("Error creating subscriber", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating subscriber"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": subscriber})
	}
}

func (h *Handler) UpdateSubscriber() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		subscriberID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid subscriber ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscriber ID"})
			return
		}

		var subscriber *model.Subscriber
		if err := c.ShouldBindJSON(&subscriber); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err = h.store.Update(c.Request.Context(), subscriber, subscriberID)
		if err != nil {
			h.logger.Error("Error updating subscriber", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating subscriber"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": subscriber})
	}
}

func (h *Handler) DeleteSubscriber() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		subscriberID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid subscriber ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscriber ID"})
			return
		}

		err = h.store.Delete(c.Request.Context(), subscriberID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Subscriber not found", slog.Any("Error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "Subscriber not found"})
				return
			}
			h.logger.Error("Error deleting subscriber", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting subscriber"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscriber deleted successfully"})
	}
}

func (h *Handler) UpdateSubscriberLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		subscriberID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid subscriber ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscriber ID"})
			return
		}

		if subscriberID <= 0 {
			h.logger.Error("Invalid subscriber ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscriber ID"})
			return
		}

		var subscriber *model.Subscriber
		if err := c.ShouldBindJSON(&subscriber); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err = h.store.LevelUp(c.Request.Context(), subscriber, subscriberID)
		if err != nil {
			h.logger.Error("Error updating level subscriber", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating level subscriber"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": subscriber})
	}
}

func (h *Handler) GetSubscribersByLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		level := c.Param("lvl")
		if level == "" {
			h.logger.Error("Invalid level subscriber", slog.Any("Error", "Invalid level subscriber"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level subscriber"})
			return
		}

		subscribers, err := h.store.GetByLevel(c.Request.Context(), level)
		if err != nil {
			h.logger.Error("Error getting subscribers by level", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting subscribers by level"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": subscribers})
	}
}
