package dto

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	Role         string
}
