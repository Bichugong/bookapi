package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}