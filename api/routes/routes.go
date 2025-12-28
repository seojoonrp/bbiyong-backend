// api/routes/routes.go

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/handlers"
	"github.com/seojoonrp/bbiyong-backend/api/repositories"
	"github.com/seojoonrp/bbiyong-backend/api/services"
	"github.com/seojoonrp/bbiyong-backend/middleware"
	"go.mongodb.org/mongo-driver/mongo"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRepo := repositories.NewUserRepository(db)
	meetingRepo := repositories.NewMeetingRepository(db)

	authService := services.NewAuthService(userRepo)
	meetingService := services.NewMeetingService(meetingRepo)

	authHandler := handlers.NewAuthHandler(authService)
	meetingHandler := handlers.NewMeetingHandler(meetingService)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "bbiyong server is running!"})
		})

		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
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
		}
	}
}
