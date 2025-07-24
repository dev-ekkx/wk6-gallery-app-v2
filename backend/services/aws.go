package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	s3Client   *s3.Client
	bucketName string
)

func rds() *gorm.DB {
	// Load environment variables
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, username, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db.AutoMigrate(&Image{})

	return db
}

func InitAWS() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	bucketName = os.Getenv("S3_BUCKET_NAME")
	fmt.Println("Bucket Name: " + bucketName)

	// Init RDS
	rds()

}

func UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	files := form.File

	for key, headers := range files {
		if !strings.HasSuffix(key, ".file") {
			continue
		}

		for i, fileHeader := range headers {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
				return
			}
			defer file.Close()

			// Extract description
			descKey := fmt.Sprintf("images[%d].description", i)
			description := c.PostForm(descKey)
			fmt.Println("Description: ", description)

			// Create unique S3 key
			s3Key := fileHeader.Filename

			db := rds()
			// Check if the image already exists in the database
			var existingImage Image
			if err := db.Where("filename= ?", fileHeader.Filename).First(&existingImage).Error; err == nil {
				c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Image with key %s already exists", s3Key)})
				return
			}

			// Upload to S3
			_, err = s3Client.PutObject(c, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(s3Key),
				Body:        file,
				ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 upload failed", "details": err.Error()})
				return
			}

			// Insert Metadata into RDS
			image := Image{
				Filename:    fileHeader.Filename,
				Description: description,
				S3Key:       s3Key,
			}
			if err := db.Create(&image).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image metadata to database"})
				return
			}

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload & DB save complete",
	})
}

func GetImages(c *gin.Context) {
	continuationToken := c.Query("continuationToken")

	// Fetch images from RDS
	db := rds()
	var images []Image
	if err := db.Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images from database"})
		return
	}

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int32(10),
	}
	if continuationToken != "" {
		input.ContinuationToken = aws.String(continuationToken)
	}

	output, err := s3Client.ListObjectsV2(c, input)
	if err != nil {
		log.Println("List error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list images, %v", err.Error())})
		return
	}
	// fetch images from S3
	var imageStructs []ImageStruct
	for _, img := range images {
		for _, item := range output.Contents {
			// Generate signed URL
			presignClient := s3.NewPresignClient(s3Client)
			presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    &img.S3Key,
			}, s3.WithPresignExpires(15*time.Minute))
			if err != nil {
				log.Println("Failed to sign URL:", err)
				continue
			}
			urlStr := presignResult.URL

			imageStructs = append(imageStructs, ImageStruct{
				Key:         img.S3Key,
				Size:        aws.Int64Value(item.Size),
				ETag:        aws.StringValue(item.ETag),
				URL:         urlStr,
				Description: img.Description,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"images": imageStructs,
		"count":  len(imageStructs),
	})
}

func DeleteImage(c *gin.Context) {
	key := c.Param("key")
	fmt.Println("Key to delete:", key)
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image key"})
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	_, err := s3Client.DeleteObject(c, input)
	if err != nil {
		var noKey *types.NoSuchKey
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &noKey) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Image with key %s not found in %s", key, bucketName)})
			return
		} else if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDenied":
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to delete the image"})
				return
			case "InvalidBucketName":
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bucket name"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete image: %s", apiErr.ErrorMessage())})
				return
			}
		} else {
			err = s3.NewObjectNotExistsWaiter(s3Client).Wait(
				context.TODO(),
				&s3.HeadObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(key),
				},
				time.Minute,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Image deletion not confirmed"})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "Image deletion confirmed"})
				return
			}
		}
	}

	// Delete from RDS
	db := rds()
	if err := db.Where("s3_key = ?", key).Delete(&Image{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image metadata from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})

}
