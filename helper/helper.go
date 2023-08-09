package helper

import (
	"log"
	"os"
	"path/filepath"
	// "log"
	// "path/filepath"
)

type BodyJson struct {
    Subject string `json:"subject"`
    Message string `json:"message"`
}

func GenerateAbsolutePath(relativePath string) string {
    currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
    }
    
    bodyPath := filepath.Join(currentPath, relativePath)
    return bodyPath
}

func ReadFile(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err!= nil {
        return nil, err
    }
    defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
        return nil, err
	}
	fileSize := fileInfo.Size()

	fileContent := make([]byte, fileSize)
	_, err = file.Read(fileContent)
	if err != nil {
        return nil, err
	}
    
    return fileContent, nil
}

