package main

import (
	https "auth-system/http"
	"auth-system/middleware"
	"auth-system/models"
	"auth-system/repositories"
	"auth-system/services"
	"auth-system/usecases"
	"errors"
	"fmt"
	"log"

	//"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setUpDatabase() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if host == "" || port == "" || user == "" || password == "" || name == "" || sslmode == "" {
		log.Fatal("One or more required DB environment variables are missing")
		fmt.Printf("%s", host)
		return nil, errors.New("Empty variable")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, name, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to get db handle: %v", err)
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping failed: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Access_Token{}, &models.ResetToken{})
	if err != nil {
		log.Fatalf("Automigrate failed: %v", err)
		return nil, err
	}

	return db, nil
}



func main() {
	router := gin.Default()
	db, err := setUpDatabase()
	if err != nil {
		fmt.Println("Failed to set up")
		return
	}

	secretKey := os.Getenv("SECRET_KEY")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")	

	repo := repositories.NewRepository(db)
	emailService := services.NewEmailService(smtpServer, smtpPort, smtpUser, smtpPassword, repo)
	useCase := usecases.NewUseCase(repo, secretKey, emailService)
	handler := https.NewHttp(useCase)
	middleware := middleware.NewAuthMiddleware(repo, secretKey)


	router.POST("user/register", handler.Register)
	router.POST("user/login", handler.Login)
	router.POST("user/logout", middleware.Authorization(), handler.Logout)
	router.POST("user/forgot", handler.Forgot)
	router.POST("user/reset", handler.Reset)

	router.Run()
}