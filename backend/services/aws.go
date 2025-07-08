package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var (
	s3Client   *s3.S3
	bucketName string
)

func InitAWS() {
	errr := godotenv.Load()
	if errr != nil {
		log.Fatal("Error loading .env file")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		log.Fatal("Failed to create AWS session:", err)
	}

	s3Client = s3.New(sess)
	bucketName = os.Getenv("AWS_BUCKET")
}

func UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form"})
		return
	}

	files := form.File["images"] // "images" is the field name
	var urls []string

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		key := file.Filename

		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        src,
			ContentType: aws.String(file.Header.Get("Content-Type")),
			ACL:         aws.String("public-read"),
		})
		if err == nil {
			url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
			urls = append(urls, url)
		}
	}

	c.JSON(http.StatusOK, gin.H{"uploaded": urls})
}
