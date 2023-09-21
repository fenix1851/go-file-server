package startup

import (
	"crypto/sha256"
	"encoding/json"
	"fileserver/cli"
	"fileserver/repository"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/google/uuid"
)

var RootPath string
var defaultAdminPass = uuid.New().String()

func declareRootPaths() string {
	osName := runtime.GOOS
	switch osName {
	case "windows":
		return "C:\\"
	case "linux":
		return "/home/"
	default:
		return "/home/"
	}
}

func createAdmin() {
	user := repository.GetUser("admin")
	flag.Parse()
	password := *cli.AdminPass

	if user.Username != "" {
		fmt.Println("admin already exists")
		fmt.Println("passsword:" + password)
		//if admin exists and we declare flag
		if password != "" {
			fmt.Println("password isn`t empty")
			hashedPassword := sha256.Sum256([]byte(password))
			user.HashedPassword = hashedPassword
			repository.UpdateUser(user)
			return
		}
		return
	}

	// if admin dont exists and we didnt declare flag
	if password == "" {
		password = defaultAdminPass
	}
	user = repository.CreateUser("admin", password, 999)
	fmt.Println("created user admin with priviliges with password: \n" + password)
}

func Init() {
	fmt.Println("Initializing...")
	CreateDirectories()
	FillSecrets()
	createAdmin()
	RootPath = declareRootPaths()
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
