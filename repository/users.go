package repository

import (
	"crypto/sha256"
	"encoding/json"
	"fileserver/utils"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

// User struct
type User struct {
	Uuid           string   `json:"uuid"`
	Username       string   `json:"username"`
	HashedPassword [32]byte `json:"password"`
	RefreshToken   string   `json:"refresh_token"`
	AccessToken    string   `json:"access_token"`
	Access         int      `json:"acces"`
}

func CreateUser(username string, password string, access int) User {
	// hash password
	hashedPassword := sha256.Sum256([]byte(password))
	// create user
	access_token, err := utils.CreateToken(username, "access")
	if err != nil {
		fmt.Println(err)
	}
	refresh_token, err := utils.CreateToken(username, "refresh")
	if err != nil {
		fmt.Println(err)
	}
	uuid := uuid.New().String()

	user := User{
		Uuid:           uuid,
		Username:       username,
		HashedPassword: hashedPassword,
		RefreshToken:   refresh_token,
		AccessToken:    access_token,
		Access:         access,
	}
	// create user file
	userFile, err := os.Create("data/users/" + username + ".json")
	if err != nil {
		fmt.Println(err)
	}
	// write user to user file
	json.NewEncoder(userFile).Encode(user)
	return user
}

func GetUser(username string) User {
	// get user file
	userFile, err := os.Open("data/users/" + username + ".json")
	if err != nil {
		fmt.Println(err)
	}
	// decode user from user file
	var user User
	json.NewDecoder(userFile).Decode(&user)
	return user
}

func GetUserByToken(token string, token_type string) (User, error) {
	//validate access token
	err := utils.ValidateToken(token)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}
	// decode access token
	tokenPayload, err := utils.DecodeToken(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(utils.GetTokenLifetime("access"), "token lifetime")
	fmt.Println(tokenPayload.Exp, "token exp")
	switch token_type {
	case "access":
		// check if token is expired
		if tokenPayload.Exp < time.Now().Unix() {
			fmt.Println("Access is expired")
			return User{}, err
		}
	case "refresh":
		// check if token is expired
		if tokenPayload.Exp < time.Now().Unix() {
			fmt.Println("Refresh is expired")
			return User{}, err
		}
	}
	// get user
	user := GetUser(tokenPayload.Username)
	return user, nil
}

func UpdateUser(user User) {
	fmt.Println("Updating user...")
	// get user file
	userFile, err := os.Open("data/users/" + user.Username + ".json")
	if err != nil {
		fmt.Println(err)
	}
	os.Remove("data/users/" + user.Username + ".json")
	// create user file
	userFile, err = os.Create("data/users/" + user.Username + ".json")
	if err != nil {
		fmt.Println(err)
	}

	// write user to user file
	json.NewEncoder(userFile).Encode(user)
}
