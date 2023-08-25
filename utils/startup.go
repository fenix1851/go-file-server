package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func Init() {
	fmt.Println("Initializing...")
	CreateDirectories()
	FillSecrets()
}

func CreateDirectories() {
	// create data directory

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		fmt.Println("Creating data directory...")
		os.Mkdir("data", 0755)
	}

}

func FillSecrets() {
	type Secrets struct {
		AccessSecret  string
		RefreshSecret string
	}
	fmt.Println("Checking for secrets...")
	// check if data/secrets.json exists
	if _, err := os.Stat("data/secrets.json"); os.IsNotExist(err) {
		fmt.Println("No secrets found!")
		// create data/secrets.json
		fmt.Println("Creating secrets...")
		secretsFile, err := os.Create("data/secrets.json")
		if err != nil {
			fmt.Println(err)
		}
		access_uuid := uuid.New().String()
		refresh_uuid := uuid.New().String()
		// make readeable secrets
		access_secret := sha256.Sum256([]byte(access_uuid))
		refresh_secret := sha256.Sum256([]byte(refresh_uuid))
		// create secrets
		secrets := Secrets{
			AccessSecret:  string(access_secret[:]),
			RefreshSecret: string(refresh_secret[:]),
		}
		// write secrets to secrets file
		json.NewEncoder(secretsFile).Encode(secrets)
	} else {
		fmt.Println("Secrets found!")
	}
}
