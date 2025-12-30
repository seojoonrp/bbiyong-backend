// main.go

package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/events"
	"github.com/seojoonrp/bbiyong-backend/api/handlers"
	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/api/routes"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/api/ws"
	"github.com/seojoonrp/bbiyong-backend/config"
	"github.com/seojoonrp/bbiyong-backend/database"
	"github.com/seojoonrp/bbiyong-backend/models"
)

// @title 삐용(BBIYONG) API
// @version 1.0
// @description 어른들의 동심 놀이 매칭 서비스 앱 삐용(BBIYONG)의 백엔드 API 문서입니다.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	config.LoadConfig()

	client, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(config.AppConfig.DBName)

	chatHub := ws.NewHub()
	go chatHub.Run()

	meetingEventChan := make(chan models.MeetingEvent, 100)

	userRepo := repositories.NewUserRepository(db)
	meetingRepo := repositories.NewMeetingRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	friendRepo := repositories.NewFriendRepository(db)

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	meetingService := services.NewMeetingService(meetingRepo, meetingEventChan)
	chatService := services.NewChatService(chatRepo, userRepo, meetingRepo)
	friendService := services.NewFriendService(friendRepo)

	authHandler := handlers.NewAuthHandler(authService)
	meetingHandler := handlers.NewMeetingHandler(meetingService)
	chatHandler := handlers.NewChatHandler(chatHub, chatService, userService)
	friendHandler := handlers.NewFriendHandler(friendService)

	go events.StartMeetingWorker(meetingEventChan, chatService, chatHub)

	router := gin.Default()
	router.Use(cors.Default())
	router.SetTrustedProxies(nil)

	routes.SetupRoutes(
		router,
		db,
		authHandler,
		meetingHandler,
		chatHandler,
		friendHandler,
	)

	port := config.AppConfig.Port
	log.Printf("Starting server on port %s.", port)
	router.Run(":" + port)
}
