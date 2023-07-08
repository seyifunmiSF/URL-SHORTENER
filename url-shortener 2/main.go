package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"url-shortener/domain"
	"url-shortener/handlers"
	"url-shortener/repositories"
	"url-shortener/services"
)

var (
	qrCodeService domain.QRCodeService
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
		panic(fmt.Sprintf("Failed to load .env file: %v", err))
	}
	// Initialize Gin router
	r := gin.Default()

	// Get the database connection string from the environment variable.
	connectionString := os.Getenv("DB_CONNECTION_STRING")

	// Connect to the database using the gorm library.
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	db.AutoMigrate(&domain.ShortenedURL{})
	redisClient := redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		Addr:     "rediss://red-cii6nh2ip7vpelqbg1r0:dBFzikkmZGXCHgKv86R7ImWytgCLdDz1@oregon-redis.render.com:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisRepository := repositories.NewRedisRepository(redisClient)

	var config = &domain.Config{
		BaseURL:         "http://localhost:8080",
		QRCodeDirectory: "./images",
		UrlImagePath:    "static",
	}

	qrCodeService = services.NewQRCodeService(config)

	shortenerRepository := repositories.NewShortenedURLRepository(db)

	shortenURLService := services.NewShortenURLService(*shortenerRepository, *redisRepository, config, qrCodeService)

	// API routes
	r.POST("/shorten", handlers.ShortenURLHandler(shortenURLService))
	r.GET("/history", handlers.LinkHistoryHandler(shortenerRepository))
	r.GET("/s/:alias", handlers.RedirectHandler(shortenURLService))
	r.GET("shorten/:alias/stats", handlers.StatisticsHandler(shortenURLService))
	r.GET("/static/:filepath", handlers.WildcardHandler)

	// Start the server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
