package handler

import (
	"adv_programming_3_4-main/internal/model"
	"adv_programming_3_4-main/internal/repository"
	"adv_programming_3_4-main/service"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"net/http"
)

type Handler struct {
	Repo         repository.UserRepository
	EmailService *service.EmailService
}

func NewHandler(repo repository.UserRepository, emailService *service.EmailService) *Handler {
	return &Handler{
		Repo:         repo,
		EmailService: emailService,
	}
}

func generateJWTToken(user model.User) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("saxfzAraC6DWW"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var newUser model.UserRegistration

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newUser.Email == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	confirmationToken, err := h.Repo.CreateUser(newUser.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	if err := sendEmail(newUser.Email, "Confirm Your Account", "Your confirmation token is: "+confirmationToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func sendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.example.com", 587, "user", "password")
	return d.DialAndSend(m)
}

func (h *Handler) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}
	err := h.Repo.ConfirmUserEmail(token)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not confirm email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email confirmed successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var creds model.LoginCredentials
	// Parse and validate login credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login credentials"})
		return
	}
	// Retrieve user by email
	user, err := h.Repo.GetUserByEmail(creds.Email)
	if err != nil {
		// User not found or other error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		return
	}
	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		// Password does not match
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		return
	}
	token, err := generateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func (h *Handler) RequestPasswordReset(c *gin.Context) {
	var emailReq struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&emailReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resetToken, err := h.Repo.CreatePasswordResetToken(emailReq.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to initiate password reset"})
		return
	}

	// Send the reset email
	resetLink := "http://yourfrontend/reset-password?token=" + resetToken
	// Call your send email function
	sendEmail(emailReq.Email, "Password Reset", "Click here to reset your password: "+resetLink)

	c.JSON(http.StatusOK, gin.H{"message": "If your email address is in our database, you will receive a password reset email shortly"})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var resetReq struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BindJSON(&resetReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.Repo.ResetPassword(resetReq.Token, resetReq.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reset password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
