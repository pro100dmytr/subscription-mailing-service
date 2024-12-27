package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"subscription-mailing-service/internal/model"
	user2 "subscription-mailing-service/storage/user"
)

type UserHandler interface {
	GetUserID() gin.HandlerFunc
	GetAllUsers() gin.HandlerFunc
	CreateUser() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
}

type Handler struct {
	store  *user2.UserStorage
	logger *slog.Logger
}

func NewHandler(store *user2.UserStorage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) GetUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, err := strconv.Atoi(idStr)

		if err != nil {
			h.logger.Error("Invalid user ID", slog.Any("Error", "Invalid user ID"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		user, err := h.store.Get(c.Request.Context(), userID)
		if err != nil {
			h.logger.Error("Error getting user", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user"})
			return
		}

		if user == nil {
			h.logger.Error("User not found", slog.Any("Error", err))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (h *Handler) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.store.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error getting users", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *model.User

		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if user.Login == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login is required"})
			return
		}

		createdUser, err := h.store.Create(c.Request.Context(), user)
		if err != nil {
			h.logger.Error("Error create user", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error create user"})
			return
		}

		c.JSON(http.StatusOK, createdUser)
	}
}

func (h *Handler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid user ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user *model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err = h.store.Update(c.Request.Context(), user, userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("Error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Error update user", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (h *Handler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid user ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.store.Delete(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("Error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Error delete user", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
