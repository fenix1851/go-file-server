package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(username string, token_type string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	var secret string
	// Adding token into payload
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	if token_type == "access" {
		secret = GetSecrets().AccessSecret
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	} else if token_type == "refresh" {
		secret = GetSecrets().RefreshSecret
		claims["exp"] = time.Now().Add(time.Hour * 24 * 14).Unix()
	}

	// Use key to sign
	secretKey := []byte(secret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error:", err)
		return "error", err
	}

	return tokenString, nil
}

type Secrets struct {
	AccessSecret  string
	RefreshSecret string
}

func GetSecrets() Secrets {
	// get secrets file
	secretsFile, err := os.Open("data/secrets.json")
	if err != nil {
		fmt.Println(err)
	}
	// decode secrets from secrets file
	var secrets Secrets
	json.NewDecoder(secretsFile).Decode(&secrets)
	return secrets
}

type TokenPayload struct {
	Username string
	Exp      int64
}

func DecodeToken(tokenString string) (TokenPayload, error) {
	// Token parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSecrets().AccessSecret), nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return TokenPayload{}, err
	}
	float_exp := token.Claims.(jwt.MapClaims)["exp"].(float64)
	exp := int64(float_exp)
	return TokenPayload{
		Username: token.Claims.(jwt.MapClaims)["username"].(string),
		Exp:      exp,
	}, nil

}

func ValidateToken(tokenString string) error {
	// Token parse
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSecrets().AccessSecret), nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func GetTokenLifetime(token_type string) int64 {
	var lifetime int64
	if token_type == "access" {
		lifetime = 86400 + time.Now().Unix()
	} else if token_type == "refresh" {
		lifetime = 86400*14 + time.Now().Unix()
	}
	return lifetime
}
