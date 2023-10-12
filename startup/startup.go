package startup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fileserver/cli"
	"fileserver/repository"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/google/uuid"
)

var RootPath string
var PORT string = ":4001"
var defaultAdminPass = uuid.New().String()

func declareRootPaths() string {
	osName := runtime.GOOS
	switch osName {
	case "windows":
		return "C:\\"
	case "linux":
		// trying to get /sjdhfgjhfdjkhg/ path
		// if it doesnt exist, return /home/
		homePath := "/home/"

		// Проверяем существование папки
		if _, err := os.Stat(homePath); os.IsNotExist(err) {
			return "/root/"
		} else {
			return "/home/"
		}

	default:
		return "/home/"
	}
}

func createAdmin() error {
	DB, err := repository.GetDBInstance()
	if err != nil {
		return err
	}
	user, err := repository.GetUser(DB, "admin")
	if err != nil {
		return err
	}
	flag.Parse()
	password := *cli.AdminPass

	if user.Username != "" {
		fmt.Println("admin already exists")
		fmt.Println("passsword:" + password)
		//if admin exists and we declare flag
		if password != "" {
			fmt.Println("password isn`t empty")
			hashedPassword := sha256.Sum256([]byte(password))
			user.HashedPassword = hex.EncodeToString(hashedPassword[:])
			repository.UpdateUser(DB, user)
			if err != nil {
				fmt.Println("Error while updating user:", err)
				return err
			}

			return nil
		}
		return nil
	}

	// if admin dont exists and we didnt declare flag
	if password == "" {
		password = defaultAdminPass
	}
	user, err = repository.CreateUser(DB, "admin", password, 999)
	if err != nil {
		panic(err)
	}
	fmt.Println("created user admin with priviliges with password: \n" + password)
	return nil
}

func setPort() {
	flag.Parse()
	port := *cli.Port
	//trying to get port from flag
	//if port is not converted to int, set default port
	intPort, err := strconv.Atoi(port)
	if err == nil {
		PORT = ":" + strconv.Itoa(intPort)
		fmt.Println("port is set to " + port)
	} else {
		fmt.Println("port is set to default")
	}
}

func Init() {
	fmt.Println("Initializing...")
	CreateDirectories()
	FillSecrets()
	createAdmin()
	setPort()
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
