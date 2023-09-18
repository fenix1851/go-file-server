package handlers

import (
	"fileserver/startup"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// RootDirectoryPath is the path to the root directory to be served
var RootDirectoryPath = startup.RootPath

func DirectoryHandler(c *gin.Context) {
	requestPath := c.Request.URL.Path
	var CurrentDirectoryPath string = startup.RootPath
	fmt.Println(requestPath, "requestPath")
	if requestPath != "/" {
		CurrentDirectoryPath = filepath.Join(startup.RootPath, requestPath)
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

	// Itterate over the directory entries
	entries, err := directory.Readdir(-1)
	if err != nil {

		c.JSON(500, gin.H{"error": err.Error()})
		return
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

	// Process directory entries
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(folderPath, entryName)
		entryPath = entryPath[len(RootDirectoryPath):]
		if entry.IsDir() {
			directories = append(directories, Directories{DirectoryName: entryName + "/", DirectoryPath: entryPath})
		} else {
			files = append(files, Files{FileName: entryName, FilePath: entryPath})
		}
	}

	// Render the template or return JSON data
	if len(directories) > 0 || len(files) > 0 {
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
