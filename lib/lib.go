package lib

import (
	"log"
	"os"
)

func FileToString(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %s\n", err)
	}
	contentString := string(content)
	return contentString
}
