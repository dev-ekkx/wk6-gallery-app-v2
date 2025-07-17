package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go/aws"
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

// func UploadImages(c *gin.Context) {
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
// 		return
// 	}

// 	files := form.File["images"]

// 	for _, file := range files {
// 		f, err := file.Open()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
// 			return
// 		}
// 		defer f.Close()

// 		key := file.Filename

// 		_, err = s3Client.PutObject(&s3.PutObjectInput{
// 			Bucket:      aws.String(bucketName),
// 			Key:         aws.String(key),
// 			Body:        f,
// 			ContentType: aws.String(file.Header.Get("Content-Type")),
// 		})

// 		if err != nil {
// 			log.Println("S3 upload error:", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"S3 upload error: ": err})
// 			return
// 		}
// 	}

// 	fmt.Println("Successfully uploaded files to S3:", output)
// 	c.JSON(http.StatusOK, gin.H{"message": "Uploaded successfully"})
// }

// func GetImages(c *gin.Context) {
// 	continuationToken := c.Query("continuationToken")

// 	input := &s3.ListObjectsV2Input{
// 		Bucket:  aws.String(bucketName),
// 		MaxKeys: aws.Int64(10),
// 	}
// 	if continuationToken != "" {
// 		input.ContinuationToken = aws.String(continuationToken)
// 	}

// 	output, err := s3Client.ListObjectsV2(input)
// 	if err != nil {
// 		log.Println("List error:", err.Error())
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list images, %v", err.Error())})
// 		return
// 	}

// 	var images []ImageStruct
// 	for _, item := range output.Contents {
// 		if *item.Key == "" || (*item.Key)[len(*item.Key)-1] == '/' {
// 			continue
// 		}

// 		// Generate signed URL
// 		req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
// 			Bucket: aws.String(bucketName),
// 			Key:    item.Key,
// 		})
// 		urlStr, err := req.Presign(15 * time.Minute)
// 		if err != nil {
// 			log.Println("Failed to sign URL:", err)
// 			continue
// 		}

// 		images = append(images, ImageStruct{
// 			Key:  *item.Key,
// 			Size: *item.Size,
// 			ETag: aws.StringValue(item.ETag),
// 			URL:  urlStr,
// 		})
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"images":      images,
// 		"nextToken":   aws.StringValue(output.NextContinuationToken),
// 		"isTruncated": aws.BoolValue(output.IsTruncated),
// 	})
// }

// func DeleteImage(c *gin.Context) {
// 	key := c.Param("key")
// 	if key == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image key"})
// 		return
// 	}

// 	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(key),
// 	})
// 	if err != nil {
// 		log.Println("Failed to delete object:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
// 		return
// 	}

// 	// Wait until the deletion is confirmed
// 	err = s3Client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(key),
// 	})
// 	if err != nil {
// 		log.Println("WaitUntilObjectNotExists failed:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image deletion not confirmed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Image deleted"})
// }

// func uploadImagesToRDS(c *gin.Context) {
// 	// This function is not implemented in the original code.
// 	// If you need to implement it, you can add the logic here.
// 	c.JSON(http.StatusNotImplemented, gin.H{"error": "uploadImagesToRDS not implemented"})
// }

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
			fmt.Println("Failed to open file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer f.Close()

		key := file.Filename

		output, err := s3Client.PutObject(c, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        f,
			ContentType: aws.String(file.Header.Get("Content-Type")),
		})

		if err != nil {
			fmt.Println("S3 upload error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"S3 upload error": err.Error()})
			return
		}

		fmt.Println("Successfully uploaded file to S3:", output)
	}
}
