package main

import (
	"adv_programming_3_4-main/internal/handler"
	"adv_programming_3_4-main/internal/repository"
	"adv_programming_3_4-main/service"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://hekxo:123456@localhost/barbershop?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.Use(handler.RateLimiter(time.Second))

	// Initialize the BarberRepository and BarberHandler
	barberRepo := repository.NewBarberRepository(db)
	barberHandler := handler.NewBarberHandler(barberRepo)

	userRepo := repository.NewSQLUserRepository(db)
	emailService := service.NewEmailService("smtp.example.com", 587, "hekxo", "password", "from@example.com")
	userHandler := handler.NewHandler(userRepo, emailService)

	// Define a route for the root path
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Barber routes
	router.GET("/barbers", barberHandler.GetBarbers)
	router.GET("/filtered-barbers", barberHandler.GetFilteredBarbers)

	// User routes
	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.Login)
	router.GET("/confirm-email", userHandler.ConfirmEmail)
	router.POST("/request-password-reset", userHandler.RequestPasswordReset)
	router.POST("/reset-password", userHandler.ResetPassword)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
