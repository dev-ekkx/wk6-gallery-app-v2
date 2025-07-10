package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dev-ekkx/wk5-gallery-app/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	port := os.Getenv("PORT")
	fmt.Println("Port" + port)

	services.InitAWS()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/upload", services.UploadImages)
		api.GET("/images", services.GetImages)
		api.DELETE("/images/:key", services.DeleteImage)
	}

	if port == "" {
		port = "8080"
	}
	err := r.Run(":" + port)
	if err != nil {
		log.Fatal("Error loading env: ", err)
		return
	}
}
