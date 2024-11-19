package processing

import (
	"log"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/streadway/amqp"
	"github.com/lib/pq" // Import the pq package for PostgreSQL array handling
	"product-management/internal/db"
)

// Initialize a GORM model for updating compressed images
type Product struct {
	ID                     uint     `gorm:"primaryKey"`
	CompressedProductImages []string `gorm:"type:text[]"`
}

// Download and compress the image
func processImage(imageURL string) (string, error) {
	log.Printf("Processing image URL: %s", imageURL)

	// Download image
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode image
	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	// Compress image
	compressed := imaging.Resize(img, 800, 0, imaging.Lanczos)

	// Save compressed image to temporary file
	filePath := "/tmp/" + strings.ReplaceAll(imageURL, "/", "_")
	err = imaging.Save(compressed, filePath)
	if err != nil {
		return "", err
	}

	// Simulate uploading to S3 (replace with actual S3 logic if needed)
	compressedURL := "https://mock-s3-storage.com/" + filePath
	log.Println("Compressed image stored at:", compressedURL)

	return compressedURL, nil
}

// ConsumeQueue processes messages from RabbitMQ
func ConsumeQueue(queueName string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue (make sure it exists before consuming)
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Waiting for messages from queue: %s", queueName)
	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)

		// Split the message body into image URLs
		imageURLs := strings.Split(string(msg.Body), ",")
		var compressedImages []string

		// Process each image URL
		for _, imageURL := range imageURLs {
			compressedURL, err := processImage(imageURL)
			if err != nil {
				log.Printf("Failed to process image %s: %v", imageURL, err)
				continue
			}
			compressedImages = append(compressedImages, compressedURL)
		}

		// Update database with compressed image URLs
		if len(compressedImages) > 0 {
			var product Product
			product.CompressedProductImages = compressedImages

			// Update the "compressed_product_images" column using pq.Array for array handling
			err = db.DB.Model(&Product{}).
				Where("product_images @> ARRAY[?]::text[]", pq.Array(imageURLs)).
				Update("compressed_product_images", pq.Array(compressedImages)).Error

			if err != nil {
				log.Printf("Failed to update database: %v", err)
			} else {
				log.Printf("Database updated with compressed images: %v", compressedImages)
			}
		}
	}
}
