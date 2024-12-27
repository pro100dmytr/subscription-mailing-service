package message

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"subscription-mailing-service/internal/model"
	message2 "subscription-mailing-service/storage/message"
)

type MessageHandler interface {
	GetMessageID() gin.HandlerFunc
	GetAllMessages() gin.HandlerFunc
	CreateMessage() gin.HandlerFunc
	UpdateMessage() gin.HandlerFunc
	DeleteMessage() gin.HandlerFunc
}

type Handler struct {
	store  *message2.MessageStorage
	logger *slog.Logger
}

func NewHandler(store *message2.MessageStorage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) GetMessageID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		messageID, err := strconv.Atoi(idStr)

		if err != nil {
			h.logger.Error("Invalid message ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
			return
		}

		message, err := h.store.Get(c.Request.Context(), messageID)
		if err != nil {
			h.logger.Error("Error getting message", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": message})
	}
}

func (h *Handler) GetAllMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		messages, err := h.store.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error getting messages", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}

func (h *Handler) CreateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var message *model.Message

		if err := c.ShouldBindJSON(&message); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		createdMessage, err := h.store.Create(c.Request.Context(), message)
		if err != nil {
			h.logger.Error("Error creating message", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": createdMessage})
	}
}

func (h *Handler) UpdateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		messageID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid message ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
			return
		}

		var message *model.Message
		if err := c.ShouldBindJSON(&message); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err = h.store.Update(c.Request.Context(), message, messageID)
		if err != nil {
			h.logger.Error("Error updating message", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": message})
	}
}

func (h *Handler) DeleteMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		messageID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid message ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
			return
		}

		err = h.store.Delete(c.Request.Context(), messageID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Message not found", slog.Any("Error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
				return
			}
			h.logger.Error("Error deleting message", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
	}
}
