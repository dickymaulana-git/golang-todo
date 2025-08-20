package helper

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

var (
	router    *mux.Router
	secretKey string = "secretkeyjwt" // Example secret key for JWT signing
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email string) (string, error) {
	var mySigningKey = []byte(secretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // Token valid for 30 minutes

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Error signing token: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
