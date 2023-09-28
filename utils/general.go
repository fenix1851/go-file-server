package utils

import (
	"crypto/sha256"
	"mime/multipart"
	"os"
	"time"
)

func HashPassword(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}

func UploadedFileToFileInfo(fileHeader *multipart.FileHeader) os.FileInfo {
	// Создайте структуру, которая реализует интерфейс os.FileInfo
	return &fileInfo{
		name:    fileHeader.Filename,
		size:    fileHeader.Size,
		mode:    0,          // Здесь можно указать соответствующий режим доступа файла
		modTime: time.Now(), // Здесь можно указать время модификации файла
	}
}

// structure that implemets os.fileInfo interface
type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi *fileInfo) Name() string {
	return fi.name
}

func (fi *fileInfo) Size() int64 {
	return fi.size
}

func (fi *fileInfo) Mode() os.FileMode {
	return fi.mode
}

func (fi *fileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *fileInfo) IsDir() bool {
	return false
}

func (fi *fileInfo) Sys() interface{} {
	return nil
}
