package controllers

import (
"net/http"
"time"
"github.com/gin-gonic/gin"
"github.com/google/uuid"
"github.com/golang-jwt/jwt/v4"
"golang.org/x/crypto/bcrypt" //  For password hashing
"github.com/Asma2347X/blog-app/models" // User struct
"github.com/Asma2347X/blog-app/config" // DB connection
)
// Secret key used to sign JWT tokens
var jwtKey = []byte("my_secret_key")

// Claims defines what data goes into the JWT
type Claims struct {
UserID uuid.UUID `json:"user_id"` // We'll store the user ID
jwt.RegisteredClaims // Includes expiry
}

// Register a new user
func Register(c *gin.Context) {
var user models.User

// Bind request body to user struct
if err := c.ShouldBindJSON(&user); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
return
}

// Hash the password before saving
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
return
}
user.Password = string(hashedPassword) // Save hashed version
user.ID = uuid.New() // Assign a unique ID

// Save user to DB
if err := config.DB.Create(&user).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
return
}

c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login authenticates the user and returns a JWT
func Login(c *gin.Context) {
var input struct {
Email string `json:"email"`
Password string `json:"password"`
}
var user models.User

// Parse input from request
if err := c.ShouldBindJSON(&input); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login"})
return
}

// Find user by email
if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
return
}

// Compare hashed password with plain password
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
return
}

// Generate JWT
expTime := time.Now().Add(24 * time.Hour)
claims := &Claims{
UserID: user.ID,
RegisteredClaims: jwt.RegisteredClaims{
ExpiresAt: jwt.NewNumericDate(expTime),
},
}
// Sign the token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, _ := token.SignedString(jwtKey)
// Return token
c.JSON(http.StatusOK, gin.H{"token": tokenString})
}