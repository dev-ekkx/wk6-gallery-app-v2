package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
	"github.com/gin-gonic/gin"
)

var (
	s3Client   *s3.Client
	bucketName string
)

func InitAWS() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	bucketName = os.Getenv("S3_BUCKET_NAME")
	fmt.Println("Bucket Name: " + bucketName)
}

func UploadImages(c *gin.Context) {
	bucketName := os.Getenv("S3_BUCKET_NAME")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	files := form.File["images"]

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer f.Close()

		key := file.Filename

		_, err = s3Client.PutObject(c, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        f,
			ContentType: aws.String(file.Header.Get("Content-Type")),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"S3 upload error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully uploaded %s to S3", key)})
	}
}

// func UploadImages(c *gin.Context) {
// 	bucketName := os.Getenv("S3_BUCKET_NAME")

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
// 		return
// 	}

// 	files := form.File
// 	results := []string{}

// 	for key, headers := range files {
// 		if !strings.HasSuffix(key, ".file") {
// 			continue
// 		}

// 		for i, fileHeader := range headers {
// 			file, err := fileHeader.Open()
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
// 				return
// 			}
// 			defer file.Close()

// 			// Extract description
// 			descKey := fmt.Sprintf("images[%d].description", i)
// 			description := c.PostForm(descKey)

// 			// Create unique S3 key
// 			s3Key := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), fileHeader.Filename)

// 			// Upload to S3
// 			_, err = s3Client.PutObject(c, &s3.PutObjectInput{
// 				Bucket:      aws.String(bucketName),
// 				Key:         aws.String(s3Key),
// 				Body:        file,
// 				ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
// 			})
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 upload failed", "details": err.Error()})
// 				return
// 			}

// 			// Insert metadata into RDS
// 			_, err = db.Exec(`
// 				INSERT INTO images (filename, description, s3_key)
// 				VALUES ($1, $2, $3)
// 			`, fileHeader.Filename, description, s3Key)

// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert failed", "details": err.Error()})
// 				return
// 			}

// 			results = append(results, fmt.Sprintf("Saved %s", fileHeader.Filename))
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Upload & DB save complete",
// 		"files":   results,
// 	})
// }

func GetImages(c *gin.Context) {
	continuationToken := c.Query("continuationToken")

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

	var images []ImageStruct
	for _, item := range output.Contents {
		if *item.Key == "" || (*item.Key)[len(*item.Key)-1] == '/' {
			continue
		}

		// Generate signed URL
		presignClient := s3.NewPresignClient(s3Client)
		presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    item.Key,
		}, s3.WithPresignExpires(15*time.Minute))
		if err != nil {
			log.Println("Failed to sign URL:", err)
			continue
		}
		urlStr := presignResult.URL

		images = append(images, ImageStruct{
			Key:  *item.Key,
			Size: *item.Size,
			ETag: aws.StringValue(item.ETag),
			URL:  urlStr,
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"images":      images,
		"nextToken":   aws.StringValue(output.NextContinuationToken),
		"isTruncated": aws.BoolValue(output.IsTruncated),
	})
}

func DeleteImage(c *gin.Context) {
	key := c.Param("key")
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
				log.Println("WaitUntilObjectNotExists failed:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Image deletion not confirmed"})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})

}
