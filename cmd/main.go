package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"subscription-mailing-service/http-server/handlers/mail"
	"subscription-mailing-service/http-server/handlers/message"
	"subscription-mailing-service/http-server/handlers/subscription"
	"subscription-mailing-service/http-server/handlers/user"
	"subscription-mailing-service/internal/config"
	mail2 "subscription-mailing-service/storage/mail"
	message2 "subscription-mailing-service/storage/message"
	subscriber2 "subscription-mailing-service/storage/subscriber"
	user2 "subscription-mailing-service/storage/user"
)

func main() {
	cfg, err := config.LoadConfig("internal/config/config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	userStorage, err := user2.NewUserStorage(cfg)
	if err != nil {
		logger.Error("Failed to initialize user storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer userStorage.Close()

	userHandler := user.NewHandler(userStorage, logger)

	router := gin.Default()

	userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("/getall", userHandler.GetAllUsers())
		userRoutes.GET("/get/:id", userHandler.GetUserID())
		userRoutes.POST("/create", userHandler.CreateUser())
		userRoutes.PUT("/update/:id", userHandler.UpdateUser())
		userRoutes.DELETE("/delete/:id", userHandler.DeleteUser())
	}

	subscriberStorage, err := subscriber2.NewSubscriberStorage(cfg)
	if err != nil {
		logger.Error("Failed to initialize subscriber storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer subscriberStorage.Close()

	subscriberHandler := subscription.NewHandler(subscriberStorage, logger)

	subscriberRoutes := router.Group("/api/subscribers")
	{
		subscriberRoutes.GET("/getall", subscriberHandler.GetAllSubscribers())
		subscriberRoutes.GET("/get/:id", subscriberHandler.GetSubscriberID())
		subscriberRoutes.POST("/create", subscriberHandler.CreateSubscriber())
		subscriberRoutes.PUT("/update/:id", subscriberHandler.UpdateSubscriber())
		subscriberRoutes.DELETE("/delete/:id", subscriberHandler.DeleteSubscriber())
		subscriberRoutes.PUT("/updatelevel/:id", subscriberHandler.UpdateSubscriberLevel())
		subscriberRoutes.GET("/getall/:lvl", subscriberHandler.GetSubscribersByLevel())
	}

	messageStorage, err := message2.NewMessageStorage(cfg)
	if err != nil {
		logger.Error("Failed to initialize message storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer messageStorage.Close()

	messageHandler := message.NewHandler(messageStorage, logger)

	messageRoutes := router.Group("/api/messages")
	{
		messageRoutes.GET("/getall", messageHandler.GetAllMessages())
		messageRoutes.GET("/get/:id", messageHandler.GetMessageID())
		messageRoutes.POST("/create", messageHandler.CreateMessage())
		messageRoutes.PUT("/update/:id", messageHandler.UpdateMessage())
		messageRoutes.DELETE("/delete/:id", messageHandler.DeleteMessage())
	}

	mailStorage, err := mail2.NewMailStorage(cfg)
	if err != nil {
		logger.Error("Failed to initialize mail storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer mailStorage.Close()

	mailHandler := mail.NewHandler(mailStorage, logger)

	mailRoutes := router.Group("/api/mails")
	{
		mailRoutes.GET("/getall", mailHandler.GetAllMails())
		mailRoutes.GET("/get/:id", mailHandler.GetMailInfo())
		mailRoutes.POST("/create", mailHandler.CreateMail())
		mailRoutes.POST("/send", mailHandler.SendMail())
		mailRoutes.PUT("/update/:id", mailHandler.UpdateMail())
		mailRoutes.DELETE("/delete/:id", mailHandler.DeleteMail())
		//mailRoutes.GET("/search/:id", mailHandler.SearchMails())
	}

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := router.Run(serverAddr); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
