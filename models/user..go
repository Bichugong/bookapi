package models

import "time"

type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"`
	Email        string `jso:"email"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}