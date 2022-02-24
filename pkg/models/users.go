package models

type User struct {
	ID           uint   `gorm:"primarykey"`
	Username     string `gorm:"unique"`
	PasswordHash string
}

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
