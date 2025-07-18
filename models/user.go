package models
import "github.com/google/uuid"

// User model maps to the "users" table in PostgreSQL
type User struct {
    ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name     string    `json:"name"`
    Email    string    `gorm:"unique"`
    Password string    `json:"password"`
	Posts    []Post     `gorm:"foreignKey:AuthorID"`
}