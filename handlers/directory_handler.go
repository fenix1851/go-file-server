package handlers

import (
	"fileserver/startup"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// RootDirectoryPath is the path to the root directory to be served
var RootDirectoryPath = startup.RootPath

func DirectoryHandler(c *gin.Context) {
	requestPath := c.Request.URL.Path
	var CurrentDirectoryPath string = startup.RootPath
	fmt.Println(requestPath, "requestPath")
	if requestPath != "/" {
		// Convert a Windows-style path to a Unix-style path
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
		fmt.Println("ER", err)
		return
	}

	fileInfos, err := directory.Readdir(-1)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//make pseudoDir came first in slice
	var entries []os.FileInfo
	entries = append(entries, pseudoDir)
	//add file infos
	entries = append(entries, fileInfos...)
	//add uploaded files

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
	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(folderPath, entryName)
		entryPath = entryPath[len(RootDirectoryPath):]

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
	//sort files
	sort.Slice(files, func(i, j int) bool {
		return files[i].FileName < files[j].FileName
	})
	//sort directories
	sort.Slice(directories, func(i, j int) bool {
		return strings.ToLower(directories[i].DirectoryName) < strings.ToLower(directories[j].DirectoryName)
	})

	// Render the template or return JSON data
	if len(directories) > 0 || len(files) > 0 {
		//clear cache
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
	fmt.Print("\n\n\n________________________________\ngetting Files...")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	absUrl := c.PostForm("absUrl")
	fmt.Print("\nABSOLUTE URL:", absUrl)
	parsedURL, err := url.Parse(absUrl)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//getting path from url
	path := parsedURL.Path
	fullPath := filepath.Join(path, file.Filename)

	fmt.Println("\nfull path:", fullPath)
	fmt.Println("current Dir Path: ", path)

	_, err = os.Stat(fullPath)
	//if file doesnt already exists save it
	if os.IsNotExist(err) {
		err = c.SaveUploadedFile(file, fullPath)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if err == nil {
		fmt.Printf("file %s already exists in directory %s\n", file.Filename, fullPath)
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("\nfile saved\n________________________________")

	DirectoryHandler(c)
}
