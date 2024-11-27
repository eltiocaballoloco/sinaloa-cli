package be

import (
    "os"
    "strings"
    "errors"
    "path/filepath"
)

// SaveFileLocally writes the data to the specified local path
func SaveFileLocally(filePath string, data []byte) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(data)
    return err
}

// ExtractFileExtension extracts the file extension from a file path
func ExtractFileExtension(filePath string) (string, error) {
    if filePath == "" {
        return "", errors.New("file path is empty")
    }
    ext := filepath.Ext(filePath)
    if ext == "" {
        return "", errors.New("no file extension found")
    }
    return ext, nil
}

// ExtractFileName extracts the name of the file from a file path.
func ExtractFileName(filePath string) (string, error) {
    if filePath == "" {
        return "", errors.New("file path is empty")
    }

    fileName := filepath.Base(filePath)
    if fileName == "." || fileName == "/" {
        return "", errors.New("invalid file path")
    }

    // Remove the extension from the file name
    ext := filepath.Ext(fileName)
    nameWithoutExt := strings.TrimSuffix(fileName, ext)

    return nameWithoutExt, nil
}
