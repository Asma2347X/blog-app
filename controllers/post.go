package controllers
import (
"net/http"
"github.com/gin-gonic/gin"
"github.com/Asma2347X/blog-app/models" // import your Post/User models
"github.com/Asma2347X/blog-app/config" // DB connection
"github.com/golang-jwt/jwt/v4"         // for JWT claims
"github.com/google/uuid"               // for UUIDs
)
// GetAllPosts returns all posts in the database
func GetAllPosts(c *gin.Context) {
var posts []models.Post
result := config.DB.Preload("Author").Find(&posts) // Preload author info (User)
if result.Error != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
return
}
c.JSON(http.StatusOK, posts) // Return list of posts as JSON
}
// GetPostByID returns a single post by ID
func GetPostByID(c *gin.Context) {
id := c.Param("id") // Get post ID from URL
var post models.Post
result := config.DB.Preload("Author").First(&post, "id = ?", id)
if result.Error != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
return
}
c.JSON(http.StatusOK, post) // Return the post
}
// CreatePost allows an authenticated user to create a new post
func CreatePost(c *gin.Context) {
// Get user ID from JWT
user, exists := c.Get("user")
if !exists {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
return
}
claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
userID := claims["user_id"].(string)
var post models.Post
if err := c.ShouldBindJSON(&post); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
post.ID = uuid.New()                // Generate UUID for the new post
post.AuthorID, _ = uuid.Parse(userID) // Set AuthorID from JWT
if err := config.DB.Create(&post).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusCreated, post) // Return the created post
}
// UpdatePost allows a user to update their own post
func UpdatePost(c *gin.Context) {
// Auth check
user, exists := c.Get("user")
if !exists {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
return
}
claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
userID := claims["user_id"].(string)
postID := c.Param("id")
var post models.Post
if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
return
}
// Check if post belongs to user
if post.AuthorID.String() != userID {
c.JSON(http.StatusForbidden, gin.H{"error": "Not your post"})
return
}
var updatedData models.Post
if err := c.ShouldBindJSON(&updatedData); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
post.Title = updatedData.Title
post.Content = updatedData.Content
config.DB.Save(&post)
c.JSON(http.StatusOK, post) // Return updated post
}
// DeletePost allows a user to delete their own post
func DeletePost(c *gin.Context) {
user, exists := c.Get("user")
if !exists {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
return
}
claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
userID := claims["user_id"].(string)
postID := c.Param("id")
var post models.Post
if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {

c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})

return
}

if post.AuthorID.String() != userID {
c.JSON(http.StatusForbidden, gin.H{"error": "Not your post"})
return
}
if err := config.DB.Delete(&post).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
return
}
}