package main

import (
	"github.com/dev-ekkx/wk5-gallery-app/services"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	services.InitAWS()

	r := gin.Default()
	r.POST("/upload", services.UploadImages)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := r.Run(":" + port)
	if err != nil {
		log.Fatal("Error loading env: ", err)
		return
	}
}
