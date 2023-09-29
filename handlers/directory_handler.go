package handlers

import (
	"fileserver/startup"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RootDirectoryPath is the path to the root directory to be served
var RootDirectoryPath = startup.RootPath

func DirectoryHandler(c *gin.Context) {
	requestPath := c.Request.URL.Path
	var CurrentDirectoryPath string = startup.RootPath
	fmt.Println(requestPath, "requestPath")
	if requestPath != "/" {
		// Преобразуем путь в стиле Windows в путь в стиле Unix
		requestPath = filepath.ToSlash(requestPath)
		CurrentDirectoryPath = filepath.Join(requestPath)
	}
	fmt.Println(CurrentDirectoryPath, "CurrentDirectoryPath")

	// Get the requested folder name from the URL
	folderPath := filepath.Join(CurrentDirectoryPath)
	folderName := filepath.Base(folderPath)

	// Get the directory contents
	directory, err := os.Open(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "Directory not found"})
			return
		}
		// if the error is Access is denied - redirect to /notallowed
		if os.IsPermission(err) {
			c.Redirect(302, "/notallowed")
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Check if the requested file is actually a directory
	fi, err := directory.Stat()
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if !fi.IsDir() {
		c.File(folderPath)
		return
	}
	defer directory.Close()

	//make pseudoDir
	parentDir := CurrentDirectoryPath + "/.."
	pseudoDir, err := os.Stat(parentDir)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("current path" + requestPath)
	fmt.Println("root path" + RootDirectoryPath)

	fileInfos, err := directory.Readdir(-1)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//get uploaded files
	uploadedFilesPath := filepath.Join("data/uploaded")
	uploadedFilesDir, err := os.Open(uploadedFilesPath)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		uploadedFilesDir.Close()
		return
	}
	uploadedFiles, err := uploadedFilesDir.Readdir(-1)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//make pseudoDir came first in slice
	var entries []os.FileInfo
	entries = append(entries, pseudoDir)
	//add file infos
	entries = append(entries, fileInfos...)
	startingUploadDirLen := len(entries)
	//add uploaded files
	if CurrentDirectoryPath == startup.RootPath {
		entries = append(entries, uploadedFiles...)
	}

	type Files struct {
		FileName string
		FilePath string
	}

	type Directories struct {
		DirectoryName string
		DirectoryPath string
	}
	// Create a slice to hold directory names and file links
	var directories []Directories
	var files []Files

	if CurrentDirectoryPath == startup.RootPath {
		entries = append(entries[:0], entries[0+1:]...)
	}

	// Process directory entries
	for i, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(folderPath, entryName)
		entryPath = entryPath[len(RootDirectoryPath):]
		if i >= startingUploadDirLen-1 {
			fmt.Print("\n1")
			baseDir, err := os.Getwd()
			if err != nil {
				fmt.Println("Ошибка при получении текущего рабочего каталога:", err)
				return
			}
			entryPath = filepath.Join(baseDir, "data", "uploaded", entryName)
		}

		entryPath = filepath.ToSlash(entryPath)
		if entry == pseudoDir {
			entryName = ".."
		}
		// remove letter: from path for Windows if it exists and if len > 3 using regexp
		regexp := regexp.MustCompile(`^[a-zA-Z]:/`)
		if regexp.MatchString(entryPath) {
			entryPath = entryPath[2:]
		}
		if entry.IsDir() {
			directories = append(directories, Directories{DirectoryName: entryName + "/", DirectoryPath: entryPath})
		} else {
			files = append(files, Files{FileName: entryName, FilePath: entryPath})
		}
	}
	// Render the template or return JSON data
	if len(directories) > 0 || len(files) > 0 {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.HTML(200, "directory.html", gin.H{
			"FolderName":  folderName,
			"Directories": directories,
			"Files":       files,
		})
	} else {
		c.JSON(200, gin.H{"message": "Empty directory"})
	}
}

func FileHandler(c *gin.Context) {
	// Get the requested file name from the URL
	filePath := c.Param("file_path")
	// Serve the file for download
	c.File(filePath)
}

func UploadHandler(c *gin.Context) {
	fmt.Print("getting Files...\n")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := fmt.Sprintf("id_%s_%s", uuid.New().String()[:8], file.Filename)
	fullPath := filepath.Join("data", "uploaded", fileName)

	err = c.SaveUploadedFile(file, fullPath)
	if err != nil {
		fmt.Printf("error has accured:%s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("konec")
	DirectoryHandler(c)
}
