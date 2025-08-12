package handlers

import (
	"net/http"
	"time"

	"github.com/musabgulfam/pumplink-backend/config"
	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the JSON payload for user registration
// This struct defines what data we expect when a user registers
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`    // User's email address (required, must be valid email)
	Password string `json:"password" binding:"required,min=6"` // User's password (required, minimum 6 characters)
}

// LoginRequest represents the JSON payload for user login
// This struct defines what data we expect when a user logs in
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"` // User's email address (required, must be valid email)
	Password string `json:"password" binding:"required"`    // User's password (required)
}

// Register handles user registration
// This endpoint allows new users to create accounts
// It validates the input, checks for existing users, creates the account, and returns success
func Register(c *gin.Context) {
	var req RegisterRequest

	// Parse and validate the JSON request body
	// The binding tag ensures email is valid and password is at least 6 characters
	if err := c.ShouldBindJSON(&req); err != nil {
		// If validation fails, return a 400 Bad Request with the error details
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	// Check if a user with this email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// If a user with this email exists, return a 409 Conflict error
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create a new user
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	if err := database.DB.Create(&user).Error; err != nil {
		// If database save fails, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Return success response with user data (excluding password)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// Login handles user authentication
// This endpoint allows existing users to log in and receive a JWT token
// It validates credentials and returns a token for subsequent authenticated requests
func Login(c *gin.Context) {
	var req LoginRequest

	// Parse and validate the JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	// Find the user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// If user not found, return a 401 Unauthorized error
		// We don't specify whether email or password is wrong for security
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Check if the provided password matches the stored hash
	if !user.CheckPassword(req.Password) {
		// If password doesn't match, return a 401 Unauthorized error
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	// JWT tokens contain user information and are signed with a secret
	// They allow users to access protected endpoints without logging in again
	cfg := config.Load()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,                              // Subject (user ID)
		"email": user.Email,                           // User's email
		"exp":   time.Now().Add(time.Hour * 2).Unix(), // Expires in 2 hours
		"iat":   time.Now().Unix(),                    // Issued at (current time)
	})

	// Sign the token with our secret key
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		// If token generation fails, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Return success response with token and user data
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
