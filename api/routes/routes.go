// api/routes/routes.go

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/handlers"
	"github.com/seojoonrp/bbiyong-backend/api/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(
	router *gin.Engine,
	db *mongo.Database,
	authHandler *handlers.AuthHandler,
	meetingHandler *handlers.MeetingHandler,
	chatHandler *handlers.ChatHandler,
	friendHandler *handlers.FriendHandler,
	saveHandler *handlers.SaveHandler,
) {
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "bbiyong server is running!"})
		})

		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.GoogleLogin)
			auth.POST("/kakao", authHandler.KakaoLogin)
			auth.POST("/apple", authHandler.AppleLogin)
			auth.GET("/check-username", authHandler.CheckUsername)
		}

		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/auth/profile", authHandler.SetProfile)

			protected.POST("/meetings", meetingHandler.CreateMeeting)
			protected.GET("/meetings/nearby", meetingHandler.GetNearby)
			protected.POST("/meetings/:id/join", meetingHandler.Join)
			protected.POST("/meetings/:id/leave", meetingHandler.Leave)

			protected.GET("/ws/meetings/:id", chatHandler.ChatConnect)
			protected.GET("/meetings/:id/chats", chatHandler.GetChatHistory)
			protected.POST("/meetings/:id/save", saveHandler.SaveMeeting)

			protected.POST("/users/:id/friend", friendHandler.RequestFriend)
			protected.PATCH("/friendships/:id/accept", friendHandler.AcceptFriend)
			protected.GET("/friends", friendHandler.GetFriendList)
		}
	}
}
