package config
import (
"fmt"
"log"
"os"
"github.com/Asma2347X/blog-app/config"     // loads .env
"gorm.io/driver/postgres"        // Postgres driver for GORM
"gorm.io/gorm"                   // GORM ORM
"github.com/joho/godotenv"
)
var DB *gorm.DB // This will be used in other files

func ConnectDB() {
// Load .env file
err := godotenv.Load()
if err != nil {
log.Fatal("Error loading .env file")
}
// Get variables from .env
host := os.Getenv("DB_HOST")
user := os.Getenv("DB_USER")
password := os.Getenv("DB_PASSWORD")
dbname := os.Getenv("DB_NAME")
port := os.Getenv("DB_PORT")
// Create connection string
dsn := fmt.Sprintf(
"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
host, user, password, dbname, port,
)

// Connect to DB
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
log.Fatal("Failed to connect to database: ", err)
}
fmt.Println("gConnected to PostgreSQL!")
DB = db
}