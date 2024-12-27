package mail

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"subscription-mailing-service/internal/model"
	mail2 "subscription-mailing-service/storage/mail"
	"time"
)

type MailHandler interface {
	GetMailInfo() gin.HandlerFunc
	GetAllMails() gin.HandlerFunc
	CreateMail() gin.HandlerFunc
	SendMail() gin.HandlerFunc
	UpdateMail() gin.HandlerFunc
	DeleteMail() gin.HandlerFunc
	SearchMails() gin.HandlerFunc
}

type Handler struct {
	store  *mail2.MailStorage
	logger *slog.Logger
}

func NewHandler(store *mail2.MailStorage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) GetMailInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		mailID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid mail ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mail ID"})
			return
		}

		mail, err := h.store.Get(c.Request.Context(), mailID)
		if err != nil {
			h.logger.Error("Error fetching mail info", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching mail info"})
			return
		}

		if mail == nil {
			h.logger.Error("Mail not found", slog.Any("Error", err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Mail not found"})
			return
		}

		c.JSON(http.StatusOK, mail)
	}
}

func (h *Handler) GetAllMails() gin.HandlerFunc {
	return func(c *gin.Context) {
		mails, err := h.store.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error getting mails", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting mails"})
			return
		}

		c.JSON(http.StatusOK, mails)
	}
}

func (h *Handler) CreateMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var mail *model.Mail

		if err := c.ShouldBindJSON(&mail); err != nil {
			h.logger.Error("Invalid request", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		createdMail, err := h.store.Create(c.Request.Context(), mail)
		if err != nil {
			h.logger.Error("Error create mail", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error create mail"})
			return
		}

		c.JSON(http.StatusOK, createdMail)
	}
}

func (h *Handler) SendMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var mail *model.Mail

		if err := c.ShouldBindJSON(&mail); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if len(mail.To) == 0 || mail.Subject == "" || mail.Body == "" {
			h.logger.Error("Missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: 'to', 'subject', 'body'"})
			return
		}

		// TODO отпрака майла
		mail.SentAt = time.Now()

		c.JSON(http.StatusOK, gin.H{
			"message": "Mail sent successfully",
			"mail":    mail,
		})
	}
}

func (h *Handler) UpdateMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		mailID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid mail ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mail ID"})
			return
		}

		var mail *model.Mail
		if err := c.ShouldBindJSON(&mail); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if len(mail.To) == 0 || mail.Subject == "" || mail.Body == "" {
			h.logger.Error("Missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: 'to', 'subject', 'body'"})
			return
		}

		err = h.store.Update(c.Request.Context(), mail, mailID)
		if err != nil {
			h.logger.Error("Error updating mail", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating mail"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Mail updated successfully",
			"mail":    mail,
		})
	}
}

func (h *Handler) DeleteMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		mailID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid mail ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mail ID"})
			return
		}

		err = h.store.Delete(c.Request.Context(), mailID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Mail not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "Mail not found"})
				return
			}
			h.logger.Error("Error deleting mail", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting mail"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Mail deleted successfully"})
	}
}

//func (h *Handler) SearchMails() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		query := c.DefaultQuery("query", "")
//		mails, err := h.store.SearchMails(c.Request.Context(), query)
//		if err != nil {
//			h.logger.Error("Error searching mails", slog.Any("error", err))
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching mails"})
//			return
//		}
//
//		c.JSON(http.StatusOK, mails)
//	}
//}
