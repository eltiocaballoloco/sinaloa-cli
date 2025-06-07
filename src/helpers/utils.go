package helpers

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ReplaceInFile(filePath, oldValue, newValue string) error {
	// Read the content of the file
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	content := string(contentBytes)

	// Exchange the old value with the new value
	modifiedContent := strings.ReplaceAll(content, oldValue, newValue)

	// Write the modified content back to the file
	err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func UpdateImageTagWithRegex(filePath, newTag string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// This regex assumes the structure under containers[0].image includes `tag: something`
	re := regexp.MustCompile(`(?m)^\s*tag:\s*.*$`)
	updatedContent := re.ReplaceAllString(string(content), fmt.Sprintf("  tag: %s", newTag))

	if err := os.WriteFile(filePath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
