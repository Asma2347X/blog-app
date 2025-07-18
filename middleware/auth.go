package middleware



import (

"net/http"

"strings"



"github.com/gin-gonic/gin"

"github.com/golang-jwt/jwt/v4"

"github.com/google/uuid"

)


// Secret must match the one used to sign JWT

var jwtKey = []byte("my_secret_key")



// AuthMiddleware protects routes by validating JWT

func AuthMiddleware() gin.HandlerFunc {

return func(c *gin.Context) {

// Read token from Authorization header

authHeader := c.GetHeader("Authorization")

if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {

c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
c.Abort()
return

}

// Strip the "Bearer " prefix
tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// Custom struct to extract the user ID

claims := &struct {
UserID uuid.UUID `json:"user_id"`

jwt.RegisteredClaims
}{}

// Parse and validate token
token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
return jwtKey, nil
})

if err != nil || !token.Valid {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
c.Abort()
return
}

// Attach user ID to context so controllers can use it
c.Set("userID", claims.UserID)
c.Next()
}
}