// main.go

package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/seojoonrp/bbiyong-backend/api/routes"
	"github.com/seojoonrp/bbiyong-backend/config"
	"github.com/seojoonrp/bbiyong-backend/database"
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

	router := gin.Default()
	router.SetTrustedProxies(nil)

	routes.SetupRoutes(router, db)

	port := config.AppConfig.Port
	log.Printf("Starting server on port %s.", port)
	router.Run(":" + port)
}