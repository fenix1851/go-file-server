package repository

import (
	"encoding/hex"
	"fileserver/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User struct
type User struct {
	gorm.Model
	Uuid           string `gorm:"uniqueIndex"`
	Username       string
	HashedPassword string
	RefreshToken   string
	AccessToken    string
	Access         int
}

func CreateUser(db *gorm.DB, username string, password string, access int) (User, error) {
	// hash password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err)
		return User{}, err
	}
	var hashedPassword [32]byte

	// Copy the bytes from []byte to [32]byte
	copy(hashedPassword[:], hashedPasswordBytes)
	// create user
	access_token, err := utils.CreateToken(username, "access")
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	refresh_token, err := utils.CreateToken(username, "refresh")
	if err != nil {
		fmt.Println(err)
	}
	uuid := uuid.New().String()

	user := User{
		Uuid:           uuid,
		Username:       username,
		HashedPassword: hex.EncodeToString(hashedPassword[:]),
		RefreshToken:   refresh_token,
		AccessToken:    access_token,
		Access:         access,
	}
	// Save the user to the database
	result := db.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

func GetUser(db *gorm.DB, username string) (User, error) {
	var user User

	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil && result.Error.Error() == "record not found" {
		return User{}, nil
	}
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User

	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
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
	user, err := GetUser(DB, tokenPayload.Username)
	if err != nil {
		fmt.Println("error while getting user:", err)
		return User{}, err
	}
	return user, nil
}

func UpdateUser(db *gorm.DB, user User) error {
	fmt.Println("Updating user...")
	existingUser, err := GetUser(db, user.Username)
	if err != nil {
		fmt.Println("Error while updating user:", err)
		return err
	}

	existingUser.Access = user.Access
	existingUser.AccessToken = user.AccessToken
	existingUser.CreatedAt = user.CreatedAt
	existingUser.DeletedAt = user.DeletedAt
	existingUser.HashedPassword = user.HashedPassword
	existingUser.ID = user.ID
	existingUser.Model = user.Model
	existingUser.RefreshToken = user.RefreshToken
	existingUser.UpdatedAt = user.UpdatedAt
	existingUser.Username = user.Username
	existingUser.Uuid = user.Uuid
	existingUser.DeletedAt.Time = user.DeletedAt.Time

	result := db.Save(&existingUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
