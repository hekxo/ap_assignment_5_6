package model

import "time"

type User struct {
	ID                int
	Email             string
	PasswordHash      string
	IsEmailConfirmed  bool
	ConfirmationToken string
	TokenCreatedAt    time.Time
}

type UserRegistration struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
