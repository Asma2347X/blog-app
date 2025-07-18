package main 
import (
    "github.com/Asma2347x/blog-appconfig"
    "github.com/Asma2347x/blog-app/models"
    "fmt"
)
func main() {
    // Connect to the PostgreSQL database using your config
    config.ConnectDB()

    // Auto migrate the User and Post models (creates tables if not exist)
    config.DB.AutoMigrate(&models.User{}, &models.Post{})
    fmt.Println("Migration complete. Database is ready.")
}